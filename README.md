<div align="center">

# 📲 PhoneDrop 

### AirDrop-style sharing between your iPhone and Windows PC

**Send photos, videos, files and text in both directions — over your own Wi-Fi.**
**No cables. No cloud. Nothing to install on your iPhone.**

<br>

![Version](https://img.shields.io/badge/version-1.4-0f6e64?style=for-the-badge)
![Windows](https://img.shields.io/badge/Windows-10%20%7C%2011-0078D4?style=for-the-badge&logo=windows&logoColor=white)
![iPhone](https://img.shields.io/badge/iPhone-Safari-000000?style=for-the-badge&logo=apple&logoColor=white)
![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go&logoColor=white)

<br>

### [⬇️ **Download PhoneDrop for Windows**](../../releases/latest)

<sub>One small `.exe` · no installer · works offline on your home network</sub>

<br>

<img src="docs/screenshot-pc.png" alt="PhoneDrop control page on Windows" width="72%">
&nbsp;
<img src="docs/screenshot-iphone.png" alt="PhoneDrop on iPhone" width="23%">

</div>

<br>

## ✨ Why PhoneDrop?

AirDrop only works between Apple devices. PhoneDrop gives you the same quick, wireless sharing between an **iPhone and a Windows PC** — just run one app and scan a QR code.

| | |
|---|---|
| ⚡ **Fast** | Files travel straight across your Wi-Fi at full local speed — not through the internet |
| 🎬 **Big-video ready** | Send large `.MOV` videos; if the Wi-Fi drops, it **continues where it stopped** |
| 📱 **No iPhone app** | Works in Safari — add it to your Home Screen to use it like an app |
| 🔁 **Both directions** | iPhone → PC *and* PC → iPhone |
| 💬 **Text & links** | Send links, notes and OTPs between devices in one tap |
| 📊 **Live progress** | Real-time speed, progress bars and time left on both screens |
| 🔒 **Private** | Nothing leaves your network; only devices that scanned your QR code can connect |
| 📶 **Works without internet** | Use your home Wi-Fi or your PC's own Mobile Hotspot |

<br>

## 🚀 Quick start

> **3 steps, about 1 minute.**

### 1️⃣ Download & run

1. Download **`PhoneDrop-1.4.zip`** from the [**Releases page**](../../releases/latest).
2. Right-click the zip → **Extract All…**
3. Double-click **`PhoneDrop.exe`**.

> [!NOTE]
> **First run only:**
> - 🛡️ **"Windows protected your PC"** → click **More info** → **Run anyway**. *(The app isn't code-signed — the full source is in this repo.)*
> - 🧱 **Windows Firewall** asks for access → tick **Private networks** → **Allow access**. *Without this, your iPhone can't connect.*

### 2️⃣ Scan the QR code

A page opens on your PC with a QR code. On your iPhone:

1. 📶 Join the **same Wi-Fi** as the PC
2. 📷 Open the **Camera** and point it at the QR code
3. 👆 Tap the link that appears

You'll see **"Connected to PC"** with a green dot. 🟢

> [!TIP]
> In Safari, tap **Share** → **Add to Home Screen**. PhoneDrop then opens like an app and you won't need to scan again.

### 3️⃣ Share!

| 📤 **iPhone → PC** | 📥 **PC → iPhone** |
|---|---|
| Tap **Send to PC** and pick photos, videos or files. | Drag files onto **Drop files here** on the PC page. |
| They're saved to **`Downloads\PhoneDrop\From iPhone`** | On the iPhone, tap **Save** (→ Files app) or **Open** |

<br>

## 🎬 Sending big videos

> [!IMPORTANT]
> **For the fastest, full-quality videos:** in the iPhone Photos picker, tap **Options** (at the top) → choose **Current** before selecting your videos.
>
> Otherwise the iPhone first **converts** each video after you tap *Add* — that can take a while and nothing moves on screen yet. That's the iPhone working, not PhoneDrop.

- 📱 **Keep the PhoneDrop page open** until it says **Done** — PhoneDrop keeps the screen awake for you.
- 🔄 If the screen locks or the Wi-Fi drops, the transfer **pauses and continues on its own**.
- 🔁 Anything that still fails gets a **Retry failed** button.

<br>

## ⚡ Get the best speed

| Tip | Why it helps |
|---|---|
| 📶 Use **5 GHz Wi-Fi** | Several times faster than 2.4 GHz |
| 📍 Stay **close to the router** | Stronger signal = higher speed |
| 🔌 Plug the PC in with **Ethernet** | Data crosses the air once instead of twice — roughly **2× faster** |
| 🔥 Use the PC's **Mobile Hotspot** | No router needed — a direct link between iPhone and PC |

<details>
<summary><b>🔥 How to use Mobile Hotspot</b></summary>

<br>

1. On the PC: **Settings → Network & internet → Mobile hotspot → On**
2. On the iPhone: join the hotspot's Wi-Fi network
3. On the PC page, choose **`192.168.137.1`** in the **Network** list
4. Scan the new QR code

</details>

<br>

## 🛠️ Troubleshooting

> [!TIP]
> **Open the Activity log first.** It's at the bottom of the iPhone page (and on the PC page). It shows every file picked, sent and saved — and the **exact error** if something fails.

<details>
<summary><b>❌ iPhone can't open the page / "Safari cannot connect"</b></summary>

<br>

- ✅ Make sure **both devices are on the same Wi-Fi**.
- ✅ Set your PC's Wi-Fi to **Private**: *Settings → Network & internet → Wi-Fi → (your network) → Network profile type → **Private***
- ✅ Allow PhoneDrop through the firewall: *Start → "Allow an app through Windows Firewall" → Change settings → tick **PhoneDrop** under **Private***
- ✅ Turn off **VPN apps** on the PC and iPhone while sharing.
- ⚠️ **Office, college, hotel and public Wi-Fi** often block devices from talking to each other — use the PC's **Mobile Hotspot** instead.

</details>

<details>
<summary><b>🔢 The QR code shows the wrong address</b></summary>

<br>

Pick another address in the **Network** list on the PC page. The right one usually starts with **`192.168.`**

</details>

<details>
<summary><b>📭 Files I picked on the iPhone don't arrive</b></summary>

<br>

Check the **Activity log**:

| Log message | What to do |
|---|---|
| *"iPhone could not read this file"* | Safari lost access to the file — pick it again (use **Options → Current** for videos) |
| *"not connected — scan the QR code again"* | Close the tab and scan the QR code again |
| *"continuing from where the PC got to"* | Nothing — the Wi-Fi blipped and PhoneDrop is resuming by itself ✅ |

</details>

<details>
<summary><b>🔴 iPhone says "PC not reachable"</b></summary>

<br>

The PhoneDrop window on the PC was closed, or the PC went to sleep. Start **`PhoneDrop.exe`** again — the iPhone reconnects by itself.

</details>

> Still stuck? **[Open an issue](../../issues/new)** and include a screenshot of the Activity log. 🙏

<br>

## 🔒 Privacy & security

- 🏠 **Everything stays on your local network.** Files and text are never uploaded to the internet.
- 🔑 **The QR link contains a secret key.** Only devices that scanned it can connect — anyone else is blocked and shown in the Activity log.
- 🖥️ **Only the PC** can see received files, manage shared files or show the QR code.
- 🌐 Connections use plain HTTP inside your network — use PhoneDrop on networks you trust (home Wi-Fi or your own hotspot), **not public Wi-Fi**.

<details>
<summary><b>🔄 How to reset the secret key</b></summary>

<br>

Close PhoneDrop, delete `%APPDATA%\PhoneDrop\token.txt`, and start PhoneDrop again. Every device will need to scan the new QR code.

</details>

<br>

## 📁 Where things are stored

| What | Location |
|---|---|
| 📥 Received from iPhone | `Downloads\PhoneDrop\From iPhone` |
| 📤 Shared to iPhone | `Downloads\PhoneDrop\To iPhone` |
| 🔑 Secret key | `%APPDATA%\PhoneDrop\token.txt` |

**Use a different folder:**

```bat
PhoneDrop.exe "D:\Shared"
```

<sub>PhoneDrop uses port **8080** (or 8081, 8088, 8090, 8765, 9090 if that's busy). Text messages are kept only while PhoneDrop is running.</sub>

<br>

## ⚙️ How it works

```mermaid
flowchart LR
    A["📱 iPhone<br/>(Safari)"] <-- "📶 Your Wi-Fi / Hotspot<br/>files · text · links" --> B["💻 PhoneDrop.exe<br/>(small web server)"]
    B --> C["📂 Downloads\PhoneDrop"]
```

**`PhoneDrop.exe`** runs a tiny web server on your PC. Your iPhone opens its page in Safari via the QR code, and files are sent straight across your local network. Each file is sent **whole**; if the connection breaks, the PC keeps what already arrived and the iPhone **resumes from that point**.

<br>

## 👩‍💻 Build from source

A single Go program with the web pages built in.

<details>
<summary><b>📂 Project structure</b></summary>

<br>

```text
├── main.go                    # the server
├── web/
│   ├── phone.html             # iPhone page
│   ├── pc.html                # PC control page
│   ├── icon.svg
│   └── icon.png
├── phonedrop.ico              # Windows app icon
├── rsrc_windows_amd64.syso    # icon resource (used automatically on Windows builds)
├── go.mod / go.sum            # dependencies (QR code library)
└── docs/                      # README screenshots
```

</details>

**Requirements:** [Go 1.24+](https://go.dev/dl/)

```bash
# 🪟 Build on Windows
go build -trimpath -ldflags "-s -w" -o PhoneDrop.exe .

# 🍎🐧 Build the Windows app from macOS / Linux
GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o PhoneDrop.exe .

# 🍎🐧 Run on macOS / Linux (files go to ~/Downloads/PhoneDrop)
go build -o phonedrop . && ./phonedrop
```

<details>
<summary><b>🎨 Optional: change the app icon</b></summary>

<br>

```bash
go install github.com/akavel/rsrc@latest
rsrc -ico phonedrop.ico -arch amd64 -o rsrc_windows_amd64.syso
```

Set `PHONEDROP_NO_BROWSER=1` to stop the browser opening automatically.

</details>

<br>

## 📝 What's new

| Version | Changes |
|---|---|
| **🆕 1.4** | 🎬 **Big `.MOV` videos no longer get stuck.** Files are sent whole; if the connection breaks, the PC keeps what arrived and the iPhone continues from there. |
| **1.3** | ⏱️ Stuck transfers detected and restarted automatically · live speed and time left |
| **1.2** | 📭 Fixed picked photos/videos sometimes not arriving · added the **Activity log** |
| **1.1** | 🔁 Reliable big uploads with automatic retry · screen stays awake while sending |
| **1.0** | 🎉 First release |

<br>

---

<div align="center">

**Made with ❤️ by Parveen rawal for everyone who owns an iPhone and a Windows PC.**

⭐ **If PhoneDrop saves you time, give it a star!** ⭐

</div>
