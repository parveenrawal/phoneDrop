PhoneDrop 1.4 - fast iPhone <-> Windows sharing over your Wi-Fi
================================================================

START
1. Double-click PhoneDrop.exe.
   - "Windows protected your PC"? Click "More info" -> "Run anyway"
     (the app isn't code-signed, that's normal for small tools).
   - Windows Firewall asks? Tick "Private networks" and click Allow.
2. A page opens in your browser with a QR code.
3. iPhone: join the SAME Wi-Fi, open the Camera, scan the code, tap the link.
4. Optional: in Safari tap Share -> "Add to Home Screen" - now it opens like an app.

iPHONE -> PC
- Tap "Send to PC", pick photos/videos/files. They save to:
  Downloads\PhoneDrop\From iPhone
- Big videos: after tapping "Add", the iPhone converts videos first
  (can take a minute). To skip that: in the Photos picker tap
  "Options" (top) -> "Current".
- Keep the PhoneDrop screen open until it says Done. If the Wi-Fi drops
  or the screen locks, it pauses and continues by itself. Anything that
  still fails gets a "Retry failed" button.

PC -> iPHONE
- Drop files on the PhoneDrop page on your PC (or click to choose).
- On the iPhone tap "Save" (goes to Files -> Downloads) or "Open"
  (for photos/videos: long-press -> Save to Photos).

TEXT & LINKS
- Send text/links/OTPs both ways from the "Text & links" box.

SPEED TIPS
- The big number on the iPhone is the LIVE speed (last 3 seconds), with
  time left.
- Each file is sent whole. If the Wi-Fi drops or a transfer gets stuck,
  the PC keeps what already arrived and the iPhone continues from there
  by itself (the log shows "continuing from where the PC got to").
- Use 5 GHz Wi-Fi, stay close to the router.
- If the PC is also on Wi-Fi, every byte crosses the air twice
  (iPhone -> router -> PC). An Ethernet cable on the PC roughly doubles speed.
- No router? Windows Settings -> Network -> Mobile hotspot -> On,
  join it from the iPhone, then pick 192.168.137.1 in the Network list.

TROUBLESHOOTING
- Files not arriving? Open "Activity log" at the bottom of the iPhone page
  (or on the PC page). It shows each file picked, sent and saved, and the
  exact error if something fails. Send me a screenshot of it.
- iPhone says it can't open the page: set your Wi-Fi to "Private" in
  Windows (Settings -> Network -> Wi-Fi -> your network -> Private),
  and allow PhoneDrop in Windows Firewall.
- Wrong network shown? Pick another one in the Network dropdown.
- Keep the black PhoneDrop window open while sharing; close it to stop.

PRIVACY
- Everything stays on your local network - nothing goes to the internet.
- The QR link contains a secret key, so only devices that scanned it can connect.
- Custom save folder: run  PhoneDrop.exe "D:\MyFolder"
