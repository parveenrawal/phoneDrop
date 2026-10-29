// PhoneDrop — fast iPhone <-> Windows file & text sharing over your local Wi-Fi.
// Run PhoneDrop.exe, scan the QR code with your iPhone camera, done.
package main

import (
	"bufio"
	"crypto/rand"
	"embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	qrcode "github.com/skip2/go-qrcode"
)

//go:embed web/*
var webFS embed.FS

const version = "1.4"

var (
	port    int
	token   string
	baseDir string // Downloads\PhoneDrop
	recvDir string // files received from iPhone
	outDir  string // files shared to iPhone
)

// ---------- live transfer tracking ----------

type transfer struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	Kind     string  `json:"kind"` // "in" (iPhone->PC), "out" (PC->iPhone), "share" (added to share list)
	Size     int64   `json:"size"`
	Done     int64   `json:"done"`
	Speed    float64 `json:"speed"` // bytes/sec
	Finished bool    `json:"finished"`
	Failed   bool    `json:"failed"`

	done    atomic.Int64
	start   time.Time
	end     time.Time
	samples []speedSample // recent progress, for a live (not average) speed
}

type speedSample struct {
	at   time.Time
	done int64
}

var (
	tMu       sync.Mutex
	transfers = map[int64]*transfer{}
	nextID    int64
)

func newTransfer(name, kind string, size int64) *transfer {
	tMu.Lock()
	defer tMu.Unlock()
	nextID++
	t := &transfer{ID: nextID, Name: name, Kind: kind, Size: size, start: time.Now()}
	transfers[t.ID] = t
	return t
}

func (t *transfer) finish(failed bool) {
	tMu.Lock()
	t.Finished, t.Failed, t.end = true, failed, time.Now()
	tMu.Unlock()
}

func snapshotTransfers() []transfer {
	tMu.Lock()
	defer tMu.Unlock()
	out := []transfer{}
	for id, t := range transfers {
		if t.Finished && time.Since(t.end) > 10*time.Second {
			delete(transfers, id)
			continue
		}
		end := time.Now()
		if t.Finished {
			end = t.end
		}
		done := t.done.Load()
		if t.Size >= 0 && done > t.Size {
			done = t.Size
		}
		sp := 0.0
		if t.Finished {
			if el := end.Sub(t.start).Seconds(); el > 0.05 {
				sp = float64(done) / el // average for finished files
			}
		} else {
			now := time.Now()
			if n := len(t.samples); n == 0 || now.Sub(t.samples[n-1].at) > 400*time.Millisecond {
				t.samples = append(t.samples, speedSample{now, done})
			}
			for len(t.samples) > 2 && now.Sub(t.samples[1].at) > 3*time.Second {
				t.samples = t.samples[1:]
			}
			if first := t.samples[0]; now.Sub(first.at) > 300*time.Millisecond {
				sp = float64(done-first.done) / now.Sub(first.at).Seconds()
			} else if el := now.Sub(t.start).Seconds(); el > 0.05 {
				sp = float64(done) / el
			}
			if sp < 0 {
				sp = 0
			}
		}
		out = append(out, transfer{ID: t.ID, Name: t.Name, Kind: t.Kind, Size: t.Size, Done: done, Speed: sp, Finished: t.Finished, Failed: t.Failed})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID > out[j].ID })
	return out
}

type countWriter struct {
	w io.Writer
	t *transfer
}

func (c *countWriter) Write(p []byte) (int, error) {
	n, err := c.w.Write(p)
	c.t.done.Add(int64(n))
	return n, err
}

type countResponse struct {
	http.ResponseWriter
	t *transfer
}

func (c *countResponse) Write(p []byte) (int, error) {
	n, err := c.ResponseWriter.Write(p)
	c.t.done.Add(int64(n))
	return n, err
}

// ---------- text clipboard ----------

type textItem struct {
	ID   int64  `json:"id"`
	From string `json:"from"`
	Text string `json:"text"`
	At   int64  `json:"at"`
}

var (
	txMu   sync.Mutex
	texts  []textItem
	textID int64
)

// ---------- activity log (iPhone reports what happens; shown in the console + PC page) ----------

var (
	logMu   sync.Mutex
	logList []string
)

func addLog(src, msg string) {
	msg = strings.ReplaceAll(strings.TrimSpace(msg), "\n", " ")
	if len(msg) > 600 {
		msg = msg[:600] + "…"
	}
	line := time.Now().Format("15:04:05") + "  [" + src + "] " + msg
	fmt.Println("  " + line)
	logMu.Lock()
	logList = append([]string{line}, logList...)
	if len(logList) > 150 {
		logList = logList[:150]
	}
	logMu.Unlock()
}

// ---------- helpers ----------

func clientIP(r *http.Request) string {
	h, _, _ := net.SplitHostPort(r.RemoteAddr)
	return h
}

func isLocal(r *http.Request) bool {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return false
	}
	ip := net.ParseIP(host)
	if ip == nil || !ip.IsLoopback() {
		return false
	}
	// Block DNS-rebinding: the PC page must be opened as localhost/127.0.0.1.
	h := r.Host
	if hh, _, err := net.SplitHostPort(h); err == nil {
		h = hh
	}
	h = strings.Trim(h, "[]")
	return h == "localhost" || h == "127.0.0.1" || h == "::1"
}

func authed(r *http.Request) bool {
	if isLocal(r) {
		return true
	}
	if r.URL.Query().Get("t") == token {
		return true
	}
	if c, err := r.Cookie("pd"); err == nil && c.Value == token {
		return true
	}
	return false
}

var badChars = strings.NewReplacer("<", "_", ">", "_", ":", "_", "\"", "_", "/", "_", "\\", "_", "|", "_", "?", "_", "*", "_")

func cleanName(n string) string {
	n = strings.ReplaceAll(n, "\\", "/")
	n = filepath.Base("/" + n)
	n = badChars.Replace(n)
	n = strings.Map(func(r rune) rune {
		if r < 32 {
			return -1
		}
		return r
	}, n)
	n = strings.TrimRight(strings.TrimSpace(n), ". ")
	if n == "" || n == "." || n == "/" {
		n = "file-" + time.Now().Format("20060102-150405")
	}
	stem := strings.ToUpper(strings.TrimSuffix(n, filepath.Ext(n)))
	switch stem {
	case "CON", "PRN", "AUX", "NUL", "COM1", "COM2", "COM3", "COM4", "LPT1", "LPT2", "LPT3":
		n = "_" + n
	}
	if len(n) > 200 {
		ext := filepath.Ext(n)
		if len(ext) > 20 {
			ext = ""
		}
		n = n[:200-len(ext)] + ext
	}
	return n
}

func uniquePath(dir, name string) string {
	p := filepath.Join(dir, name)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return p
	}
	ext := filepath.Ext(name)
	stem := strings.TrimSuffix(name, ext)
	for i := 1; ; i++ {
		p = filepath.Join(dir, fmt.Sprintf("%s (%d)%s", stem, i, ext))
		if _, err := os.Stat(p); os.IsNotExist(err) {
			return p
		}
	}
}

func listDir(dir string) []map[string]any {
	ents, _ := os.ReadDir(dir)
	out := []map[string]any{}
	for _, e := range ents {
		if e.IsDir() || strings.HasSuffix(e.Name(), ".part") || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		out = append(out, map[string]any{"name": e.Name(), "size": info.Size(), "mod": info.ModTime().UnixMilli()})
	}
	sort.Slice(out, func(i, j int) bool { return out[i]["mod"].(int64) > out[j]["mod"].(int64) })
	return out
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	json.NewEncoder(w).Encode(v)
}

func randHex(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// lanIPs returns usable IPv4 addresses, best guess first.
func lanIPs() []string {
	var ips []string
	seen := map[string]bool{}
	add := func(ip string) {
		if ip != "" && !seen[ip] {
			seen[ip] = true
			ips = append(ips, ip)
		}
	}
	// The address used for the default route (no packets are actually sent).
	if c, err := net.Dial("udp", "8.8.8.8:80"); err == nil {
		if a, ok := c.LocalAddr().(*net.UDPAddr); ok && a.IP.IsPrivate() {
			add(a.IP.String())
		}
		c.Close()
	}
	ifaces, _ := net.Interfaces()
	skip := []string{"vethernet", "virtualbox", "vmware", "wsl", "hyper-v", "loopback", "bluetooth", "docker", "veth", "br-", "tailscale", "zerotier", "vpn", "tap", "tun"}
	for _, ifc := range ifaces {
		if ifc.Flags&net.FlagUp == 0 || ifc.Flags&net.FlagLoopback != 0 {
			continue
		}
		ln := strings.ToLower(ifc.Name)
		bad := false
		for _, s := range skip {
			if strings.Contains(ln, s) {
				bad = true
			}
		}
		addrs, _ := ifc.Addrs()
		for _, a := range addrs {
			ipn, ok := a.(*net.IPNet)
			if !ok || ipn.IP.To4() == nil || !ipn.IP.IsPrivate() {
				continue
			}
			// Windows Mobile Hotspot is always 192.168.137.1 — keep it even if the adapter looks virtual.
			if bad && ipn.IP.String() != "192.168.137.1" {
				continue
			}
			add(ipn.IP.String())
		}
	}
	return ips
}

func phoneURL(ip string) string {
	return fmt.Sprintf("http://%s:%d/?t=%s", ip, port, token)
}

func openBrowser(u string) {
	switch runtime.GOOS {
	case "windows":
		exec.Command("rundll32", "url.dll,FileProtocolHandler", u).Start()
	case "darwin":
		exec.Command("open", u).Start()
	default:
		exec.Command("xdg-open", u).Start()
	}
}

func openFolder(dir string) {
	switch runtime.GOOS {
	case "windows":
		exec.Command("explorer", dir).Start()
	case "darwin":
		exec.Command("open", dir).Start()
	default:
		exec.Command("xdg-open", dir).Start()
	}
}

func loadToken() string {
	cfg, err := os.UserConfigDir()
	if err != nil {
		cfg = baseDir
	}
	dir := filepath.Join(cfg, "PhoneDrop")
	os.MkdirAll(dir, 0o700)
	p := filepath.Join(dir, "token.txt")
	if b, err := os.ReadFile(p); err == nil {
		if t := strings.TrimSpace(string(b)); len(t) >= 16 {
			return t
		}
	}
	t := randHex(10)
	os.WriteFile(p, []byte(t), 0o600)
	return t
}

// ---------- handlers ----------

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		http.Error(w, "use PUT or POST", 405)
		return
	}
	dest, kind := recvDir, "in"
	if r.URL.Query().Get("to") == "share" {
		if !isLocal(r) {
			http.Error(w, "forbidden", 403)
			return
		}
		dest, kind = outDir, "share"
	}
	ct := r.Header.Get("Content-Type")
	saved := []string{}

	saveStream := func(name string, size int64, body io.Reader) error {
		name = cleanName(name)
		t := newTransfer(name, kind, size)
		tmp := filepath.Join(dest, "."+randHex(6)+".part")
		f, err := os.Create(tmp)
		if err != nil {
			t.finish(true)
			return err
		}
		bw := bufio.NewWriterSize(f, 4<<20)
		_, err = io.CopyBuffer(&countWriter{bw, t}, body, make([]byte, 1<<20))
		if err == nil {
			err = bw.Flush()
		}
		if cerr := f.Close(); err == nil {
			err = cerr
		}
		if err != nil {
			os.Remove(tmp)
			t.finish(true)
			return err
		}
		final := uniquePath(dest, name)
		if err := os.Rename(tmp, final); err != nil {
			os.Remove(tmp)
			t.finish(true)
			return err
		}
		t.finish(false)
		saved = append(saved, filepath.Base(final))
		return nil
	}

	if mt, _, _ := mime.ParseMediaType(ct); mt == "multipart/form-data" {
		// Used by iPhone Shortcuts ("Form" body) or plain HTML forms.
		mr, err := r.MultipartReader()
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		for {
			part, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
			if part.FileName() == "" {
				// A text field from a Shortcut → treat as text message.
				b, _ := io.ReadAll(io.LimitReader(part, 1<<20))
				if s := strings.TrimSpace(string(b)); s != "" {
					addText("iPhone", s)
				}
				continue
			}
			if err := saveStream(part.FileName(), -1, part); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
		}
	} else {
		name := r.URL.Query().Get("name")
		if name == "" {
			name = r.Header.Get("X-File-Name")
		}
		if name == "" {
			ext := ""
			if exts, _ := mime.ExtensionsByType(ct); len(exts) > 0 {
				ext = exts[0]
			}
			name = "iPhone-" + time.Now().Format("20060102-150405") + ext
		}
		if un, err := url.QueryUnescape(name); err == nil {
			name = un
		}
		if err := saveStream(name, r.ContentLength, r.Body); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}
	writeJSON(w, map[string]any{"ok": true, "saved": saved})
}

// ---------- resumable uploads ----------
// The iPhone sends a whole file in one request. If the connection breaks, everything
// received so far is kept, and the iPhone asks /api/status how much arrived and sends
// only the rest. Ranges may overlap (a stale request can still be finishing); the file
// is complete when the contiguous data from byte 0 reaches the full size.

type upState struct {
	mu     sync.Mutex
	f      *os.File
	tmp    string
	name   string
	total  int64
	ranges map[int64]int64 // offset -> longest length received there
	t      *transfer
	last   time.Time
	done   bool
	final  string
}

var chunkHook func(off, total int64, body io.Reader) io.Reader // test builds only

var (
	upMu    sync.Mutex
	uploads = map[string]*upState{}
)

func validID(id string) bool {
	if len(id) < 8 || len(id) > 64 {
		return false
	}
	for _, c := range id {
		if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '-') {
			return false
		}
	}
	return true
}

// have returns how many bytes from the start of the file have arrived (call with u.mu held).
func (u *upState) have() int64 {
	offs := make([]int64, 0, len(u.ranges))
	for o := range u.ranges {
		offs = append(offs, o)
	}
	sort.Slice(offs, func(i, j int) bool { return offs[i] < offs[j] })
	var cur int64
	for _, o := range offs {
		if o > cur {
			break
		}
		if e := o + u.ranges[o]; e > cur {
			cur = e
		}
	}
	return cur
}

func (u *upState) record(off, n int64) {
	if n > u.ranges[off] {
		u.ranges[off] = n
	}
	u.last = time.Now()
}

// finishIfComplete renames the file into place once all bytes are there (call with u.mu held).
func (u *upState) finishIfComplete(r *http.Request) error {
	if u.done || u.have() < u.total {
		return nil
	}
	u.done = true
	cerr := u.f.Close()
	final := uniquePath(recvDir, u.name)
	if cerr != nil || os.Rename(u.tmp, final) != nil {
		os.Remove(u.tmp)
		u.t.finish(true)
		return fmt.Errorf("could not save file")
	}
	u.final = filepath.Base(final)
	u.t.done.Store(u.total)
	u.t.finish(false)
	addLog("PC", fmt.Sprintf("Saved %s (%.1f MB) from %s", u.final, float64(u.total)/1048576, clientIP(r)))
	return nil
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	upMu.Lock()
	u := uploads[r.URL.Query().Get("id")]
	upMu.Unlock()
	if u == nil {
		writeJSON(w, map[string]any{"known": false, "have": 0, "done": false})
		return
	}
	u.mu.Lock()
	defer u.mu.Unlock()
	writeJSON(w, map[string]any{"known": true, "have": u.have(), "done": u.done, "saved": u.final})
}

func handleChunk(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		http.Error(w, "use PUT", 405)
		return
	}
	q := r.URL.Query()
	id := q.Get("id")
	off, err1 := strconv.ParseInt(q.Get("off"), 10, 64)
	total, err2 := strconv.ParseInt(q.Get("total"), 10, 64)
	if !validID(id) || err1 != nil || err2 != nil || off < 0 || total < 0 || off > total {
		http.Error(w, "bad request", 400)
		return
	}

	upMu.Lock()
	u := uploads[id]
	if u == nil && off != 0 {
		// We don't know this upload (PhoneDrop was restarted) — tell the iPhone to start the file again.
		upMu.Unlock()
		http.Error(w, "unknown upload", 409)
		return
	}
	if u == nil {
		name := cleanName(q.Get("name"))
		tmp := filepath.Join(recvDir, "."+id+".part")
		f, err := os.OpenFile(tmp, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o644)
		if err != nil {
			upMu.Unlock()
			http.Error(w, err.Error(), 500)
			return
		}
		u = &upState{f: f, tmp: tmp, name: name, total: total, ranges: map[int64]int64{}, t: newTransfer(name, "in", total), last: time.Now()}
		uploads[id] = u
	}
	upMu.Unlock()

	u.mu.Lock()
	if u.total != total {
		u.mu.Unlock()
		http.Error(w, "size mismatch", 400)
		return
	}
	if u.done {
		final := u.final
		u.mu.Unlock()
		writeJSON(w, map[string]any{"ok": true, "done": true, "saved": final, "have": total})
		return
	}
	if total == 0 { // empty file
		err := u.finishIfComplete(r)
		u.mu.Unlock()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		writeJSON(w, map[string]any{"ok": true, "done": true, "saved": u.final, "have": 0})
		return
	}
	u.mu.Unlock()

	var rbody io.Reader = r.Body
	if chunkHook != nil {
		rbody = chunkHook(off, total, rbody)
	}
	// Write as the bytes arrive, recording progress every few MB so a broken
	// connection keeps what already arrived.
	ow := io.NewOffsetWriter(u.f, off)
	body := io.LimitReader(rbody, total-off)
	buf := make([]byte, 1<<20)
	var n int64
	var err error
	lastRec := int64(0)
	for {
		k, rerr := body.Read(buf)
		if k > 0 {
			if _, werr := ow.Write(buf[:k]); werr != nil {
				err = werr
				break
			}
			n += int64(k)
			u.t.done.Add(int64(k))
			if n-lastRec >= 4<<20 {
				u.mu.Lock()
				u.record(off, n)
				u.mu.Unlock()
				lastRec = n
			}
		}
		if rerr == io.EOF {
			break
		}
		if rerr != nil {
			err = rerr
			break
		}
	}
	if err == nil && r.ContentLength >= 0 && n != r.ContentLength {
		err = fmt.Errorf("connection ended early")
	}

	u.mu.Lock()
	defer u.mu.Unlock()
	if n > 0 && !u.done {
		u.record(off, n)
	}
	if err != nil {
		have := u.have()
		addLog("PC", fmt.Sprintf("%s: connection broke at %.1f of %.1f MB (%v) - kept what arrived", u.name, float64(have)/1048576, float64(u.total)/1048576, err))
		http.Error(w, err.Error(), 500)
		return
	}
	if ferr := u.finishIfComplete(r); ferr != nil {
		http.Error(w, ferr.Error(), 500)
		return
	}
	writeJSON(w, map[string]any{"ok": true, "done": u.done, "saved": u.final, "have": u.have()})
}

// Forget finished uploads after a while; drop abandoned ones after 2 hours.
func cleanupUploads() {
	for range time.Tick(time.Minute) {
		upMu.Lock()
		for id, u := range uploads {
			u.mu.Lock()
			idle := time.Since(u.last)
			if u.done && idle > 15*time.Minute {
				delete(uploads, id)
			} else if !u.done && idle > 2*time.Hour {
				u.f.Close()
				os.Remove(u.tmp)
				u.t.finish(true)
				delete(uploads, id)
			}
			u.mu.Unlock()
		}
		upMu.Unlock()
	}
}

func handleDownload(w http.ResponseWriter, r *http.Request) {
	name := cleanName(strings.TrimPrefix(r.URL.Path, "/dl/"))
	dir := outDir
	if r.URL.Query().Get("from") == "recv" && isLocal(r) {
		dir = recvDir
	}
	p := filepath.Join(dir, name)
	f, err := os.Open(p)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil || st.IsDir() {
		http.NotFound(w, r)
		return
	}
	disp := "attachment"
	if r.URL.Query().Get("view") == "1" {
		disp = "inline"
	}
	w.Header().Set("Content-Disposition", fmt.Sprintf("%s; filename=\"%s\"; filename*=UTF-8''%s",
		disp, strings.Map(func(r rune) rune {
			if r > 126 || r == '"' {
				return '_'
			}
			return r
		}, name), url.PathEscape(name)))
	w.Header().Set("Cache-Control", "no-store")
	// Small range requests (video previews/seeking) aren't shown as transfers.
	if rg := r.Header.Get("Range"); r.Method == http.MethodHead || (rg != "" && !strings.HasPrefix(rg, "bytes=0-")) {
		http.ServeContent(w, r, name, st.ModTime(), f)
		return
	}
	t := newTransfer(name, "out", st.Size())
	http.ServeContent(&countResponse{w, t}, r, name, st.ModTime(), f)
	t.finish(false)
}

func addText(from, s string) {
	txMu.Lock()
	defer txMu.Unlock()
	textID++
	texts = append([]textItem{{ID: textID, From: from, Text: s, At: time.Now().UnixMilli()}}, texts...)
	if len(texts) > 30 {
		texts = texts[:30]
	}
}

func handleText(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", 405)
		return
	}
	var body struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(io.LimitReader(r.Body, 1<<20)).Decode(&body); err != nil || strings.TrimSpace(body.Text) == "" {
		http.Error(w, "empty", 400)
		return
	}
	from := "iPhone"
	if isLocal(r) {
		from = "PC"
	}
	addText(from, body.Text)
	writeJSON(w, map[string]any{"ok": true})
}

func handleState(w http.ResponseWriter, r *http.Request) {
	txMu.Lock()
	tx := append([]textItem{}, texts...)
	txMu.Unlock()
	st := map[string]any{
		"shared":    listDir(outDir),
		"texts":     tx,
		"transfers": snapshotTransfers(),
		"version":   version,
	}
	if isLocal(r) {
		ips := lanIPs()
		urls := []map[string]string{}
		for _, ip := range ips {
			urls = append(urls, map[string]string{"ip": ip, "url": phoneURL(ip)})
		}
		logMu.Lock()
		st["log"] = append([]string{}, logList...)
		logMu.Unlock()
		st["received"] = listDir(recvDir)
		st["urls"] = urls
		st["recvDir"] = recvDir
		st["port"] = port
	}
	writeJSON(w, st)
}

func handleQR(w http.ResponseWriter, r *http.Request) {
	ip := r.URL.Query().Get("ip")
	if net.ParseIP(ip) == nil {
		if ips := lanIPs(); len(ips) > 0 {
			ip = ips[0]
		} else {
			ip = "127.0.0.1"
		}
	}
	png, err := qrcode.Encode(phoneURL(ip), qrcode.Medium, 640)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "no-store")
	w.Write(png)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", 405)
		return
	}
	q := r.URL.Query()
	switch q.Get("what") {
	case "shared":
		os.Remove(filepath.Join(outDir, cleanName(q.Get("name"))))
	case "text":
		id, _ := strconv.ParseInt(q.Get("id"), 10, 64)
		txMu.Lock()
		for i, t := range texts {
			if t.ID == id {
				texts = append(texts[:i], texts[i+1:]...)
				break
			}
		}
		txMu.Unlock()
	}
	writeJSON(w, map[string]any{"ok": true})
}

func page(name string, w http.ResponseWriter) {
	b, _ := webFS.ReadFile("web/" + name)
	s := strings.ReplaceAll(string(b), "{{TOKEN}}", token)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	io.WriteString(w, s)
}

func router(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	if p == "/manifest.json" || p == "/icon.png" || p == "/icon.svg" {
		serveAsset(w, r)
		return
	}
	if !authed(r) {
		if p == "/" || strings.HasPrefix(p, "/api/") {
			addLog("PC", "Blocked a device without the QR key: "+clientIP(r)+" "+p+" - scan the QR code again")
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(401)
		io.WriteString(w, `<!doctype html><meta name=viewport content="width=device-width,initial-scale=1"><body style="font:17px -apple-system,system-ui;padding:32px;text-align:center;color:#222"><h2>PhoneDrop</h2><p>Scan the QR code shown on your PC to connect.</p></body>`)
		return
	}
	// Remember the iPhone so links/home-screen bookmark keep working.
	if r.URL.Query().Get("t") == token && !isLocal(r) {
		http.SetCookie(w, &http.Cookie{Name: "pd", Value: token, Path: "/", MaxAge: 400 * 24 * 3600, SameSite: http.SameSiteLaxMode, HttpOnly: true})
	}
	localOnly := func(h http.HandlerFunc) {
		if !isLocal(r) {
			http.Error(w, "forbidden", 403)
			return
		}
		h(w, r)
	}
	switch {
	case p == "/":
		if isLocal(r) {
			page("pc.html", w)
		} else {
			page("phone.html", w)
		}
	case p == "/phone":
		page("phone.html", w)
	case p == "/api/log":
		b, _ := io.ReadAll(io.LimitReader(r.Body, 4096))
		src := "iPhone"
		if isLocal(r) {
			src = "PC"
		}
		addLog(src, string(b))
		writeJSON(w, map[string]any{"ok": true})
	case p == "/api/status":
		handleStatus(w, r)
	case p == "/api/chunk":
		handleChunk(w, r)
	case p == "/api/upload":
		handleUpload(w, r)
	case strings.HasPrefix(p, "/dl/"):
		handleDownload(w, r)
	case p == "/api/text":
		handleText(w, r)
	case p == "/api/state":
		handleState(w, r)
	case p == "/qr.png":
		localOnly(handleQR)
	case p == "/api/delete":
		localOnly(handleDelete)
	case p == "/api/open":
		localOnly(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Query().Get("which") == "shared" {
				openFolder(outDir)
			} else {
				openFolder(recvDir)
			}
			writeJSON(w, map[string]any{"ok": true})
		})
	case p == "/api/quit":
		localOnly(func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, map[string]any{"ok": true})
			go func() { time.Sleep(300 * time.Millisecond); os.Exit(0) }()
		})
	default:
		http.NotFound(w, r)
	}
}

func serveAsset(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/icon.svg":
		b, _ := webFS.ReadFile("web/icon.svg")
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Write(b)
	case "/icon.png":
		b, _ := webFS.ReadFile("web/icon.png")
		w.Header().Set("Content-Type", "image/png")
		w.Write(b)
	case "/manifest.json":
		writeJSON(w, map[string]any{
			"name": "PhoneDrop", "short_name": "PhoneDrop", "display": "standalone",
			"start_url": "/?t=" + token, "background_color": "#f6f4ef", "theme_color": "#f6f4ef",
			"icons": []map[string]string{{"src": "/icon.png", "sizes": "180x180", "type": "image/png"}},
		})
	}
}

func main() {
	home, _ := os.UserHomeDir()
	baseDir = filepath.Join(home, "Downloads", "PhoneDrop")
	if len(os.Args) > 1 && os.Args[1] != "" {
		baseDir = os.Args[1]
	}
	recvDir = filepath.Join(baseDir, "From iPhone")
	outDir = filepath.Join(baseDir, "To iPhone")
	for _, d := range []string{recvDir, outDir} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			fmt.Println("Cannot create folder:", d, err)
			pause()
			return
		}
	}
	// Clean up leftovers from interrupted transfers.
	for _, d := range []string{recvDir, outDir} {
		ents, _ := os.ReadDir(d)
		for _, e := range ents {
			if strings.HasSuffix(e.Name(), ".part") {
				os.Remove(filepath.Join(d, e.Name()))
			}
		}
	}
	token = loadToken()
	go cleanupUploads()

	var ln net.Listener
	var err error
	for _, p := range []int{8080, 8081, 8088, 8090, 8765, 9090, 0} {
		ln, err = net.Listen("tcp4", fmt.Sprintf("0.0.0.0:%d", p))
		if err == nil {
			break
		}
	}
	if err != nil {
		fmt.Println("Could not start server:", err)
		pause()
		return
	}
	port = ln.Addr().(*net.TCPAddr).Port

	srv := &http.Server{
		Handler:           http.HandlerFunc(router),
		ReadHeaderTimeout: 15 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	ips := lanIPs()
	fmt.Println()
	fmt.Println("  ==============================================")
	fmt.Println("    PhoneDrop " + version + "  -  iPhone <-> PC sharing")
	fmt.Println("  ==============================================")
	fmt.Println()
	fmt.Printf("  Control page:  http://localhost:%d\n", port)
	if len(ips) > 0 {
		fmt.Printf("  iPhone link:   %s\n", phoneURL(ips[0]))
	} else {
		fmt.Println("  No Wi-Fi/LAN network found - connect the PC to Wi-Fi or turn on Mobile Hotspot.")
	}
	fmt.Printf("  Files saved to: %s\n", baseDir)
	fmt.Println()
	fmt.Println("  If Windows Firewall asks, tick 'Private networks' and click Allow.")
	fmt.Println("  Keep this window open while sharing. Close it to stop PhoneDrop.")
	fmt.Println()

	if os.Getenv("PHONEDROP_NO_BROWSER") == "" {
		go func() {
			time.Sleep(400 * time.Millisecond)
			openBrowser(fmt.Sprintf("http://localhost:%d/", port))
		}()
	}
	if err := srv.Serve(ln); err != nil {
		fmt.Println("Server stopped:", err)
		pause()
	}
}

func pause() {
	fmt.Println("Press Enter to exit...")
	bufio.NewReader(os.Stdin).ReadString('\n')
}
