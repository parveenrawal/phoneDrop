=====================================================================
  PhoneDrop 1.4
  Fast iPhone <-> Windows file sharing over your own Wi-Fi or hotspot
=====================================================================

PhoneDrop is an AirDrop-style tool for Windows. Run one small program on
your PC, scan a QR code with your iPhone, and send photos, videos, files,
text and links in both directions.

  - Nothing to install on the iPhone (it uses Safari)
  - Files go straight across your local network - never to the internet
  - Big videos supported, with automatic resume if the Wi-Fi drops
  - Live speed, progress and time left on both screens


---------------------------------------------------------------------
 CONTENTS
---------------------------------------------------------------------
  1. Requirements
  2. Download and run
  3. Connect your iPhone
  4. Sending files: iPhone -> PC
  5. Sending files: PC -> iPhone
  6. Text and links
  7. Getting the best speed
  8. Troubleshooting
  9. Privacy and security
 10. Where things are stored
 11. Build from source
 12. What's new in 1.4


---------------------------------------------------------------------
 1. REQUIREMENTS
---------------------------------------------------------------------
  PC      Windows 10 or 11, 64-bit
  Phone   iPhone or iPad with Safari (any recent iOS)
          Android phones also work in Chrome.
  Network PC and phone on the SAME Wi-Fi, or the phone connected to
          the PC's Mobile Hotspot. No internet connection is needed.


---------------------------------------------------------------------
 2. DOWNLOAD AND RUN
---------------------------------------------------------------------
  1. Download PhoneDrop-1.4.zip from the Releases page of this repo.
  2. Right-click the zip -> "Extract All...".
  3. Open the extracted folder and double-click PhoneDrop.exe.

  First run only:
  - "Windows protected your PC" (SmartScreen):
      click "More info" -> "Run anyway".
      This appears because the app isn't code-signed. You can check the
      code yourself in the "source" folder, or build it (section 11).
  - Windows Firewall asks for access:
      tick "Private networks" and click "Allow access".
      Without this, the iPhone can't connect.

  A black window opens (keep it open - closing it stops PhoneDrop) and
  your browser opens the PhoneDrop control page with a QR code.


---------------------------------------------------------------------
 3. CONNECT YOUR IPHONE
---------------------------------------------------------------------
  1. Join the same Wi-Fi as the PC.
  2. Open the Camera app and point it at the QR code on the PC.
  3. Tap the link that appears. PhoneDrop opens in Safari.

  The iPhone page shows "Connected to PC" with a green dot, and the
  version (v1.4) at the bottom.

  Tip: in Safari tap Share -> "Add to Home Screen". PhoneDrop then opens
  like an app, and you don't need to scan again next time.

  The QR code stays the same each time you run PhoneDrop on that PC.


---------------------------------------------------------------------
 4. SENDING FILES: IPHONE -> PC
---------------------------------------------------------------------
  1. Tap "Send to PC".
  2. Choose "Photo Library" (or "Choose File" for the Files app) and
     select one or many photos, videos or files.
  3. Tap "Add". Sending starts by itself.

  Files are saved on the PC in:
      Downloads\PhoneDrop\From iPhone
  (Click "Open received folder" on the PC page to go there.)

  Big videos:
  - After tapping "Add", the iPhone may first convert videos to a more
    compatible format. This can take a while and nothing moves yet -
    that is the iPhone working, not PhoneDrop.
  - To skip the conversion and send the original video at full
    quality: in the Photos picker tap "Options" (at the top) and choose
    "Current" before selecting.

  Keep the PhoneDrop page open on the iPhone until it says "Done".
  If the screen locks or the Wi-Fi drops, the transfer pauses and then
  continues from where it stopped. Anything that still fails shows a
  "Retry failed" button.


---------------------------------------------------------------------
 5. SENDING FILES: PC -> IPHONE
---------------------------------------------------------------------
  1. On the PC page, drag files onto "Drop files here" (or click it).
  2. They appear on the iPhone under "From your PC" within a second
     or two.
  3. On the iPhone:
       "Save" -> downloads the file (Files app -> Downloads)
       "Open" -> opens photos, videos and PDFs; to put a photo or video
                 in your Photos app, long-press it -> "Save to Photos"

  Click the X next to a file on the PC page to stop sharing it.


---------------------------------------------------------------------
 6. TEXT AND LINKS
---------------------------------------------------------------------
  Use the "Text & links" box on either device to send text, links or
  OTPs. On the PC press Ctrl+Enter to send. Tap "Copy" on the other
  device; links also get an "Open" button.


---------------------------------------------------------------------
 7. GETTING THE BEST SPEED
---------------------------------------------------------------------
  - Connect the iPhone to the 5 GHz Wi-Fi network, not 2.4 GHz.
  - Stay close to the router.
  - If the PC is on Wi-Fi too, every file crosses the air twice
    (iPhone -> router -> PC). Plugging the PC into the router with an
    Ethernet cable roughly doubles the speed.
  - No router, or a slow one? Use the PC's own hotspot:
      Windows Settings -> Network & internet -> Mobile hotspot -> On
      Join that hotspot on the iPhone, then on the PC page choose
      192.168.137.1 in the "Network" list and scan the new QR code.
  - Send many files in one go - PhoneDrop sends two at a time and starts
    small files first.

  The big number on the iPhone is the live speed (last 3 seconds), with
  the time left. The PC page shows the same under "Live transfers".


---------------------------------------------------------------------
 8. TROUBLESHOOTING
---------------------------------------------------------------------
  The iPhone can't open the page / "Safari cannot connect"
    - Make sure both devices are on the same Wi-Fi.
    - Set the PC's Wi-Fi to Private:
        Settings -> Network & internet -> Wi-Fi -> (your network)
        -> Network profile type -> Private
    - Allow PhoneDrop through the firewall:
        Start -> "Allow an app through Windows Firewall" ->
        Change settings -> tick PhoneDrop under "Private".
    - VPN apps on the PC or iPhone can block local connections. Turn
      them off while sharing.
    - Office, college, hotel and public Wi-Fi often block devices from
      reaching each other. Use the PC's Mobile Hotspot instead.

  The QR code / Network list shows the wrong address
    - Pick another address in the "Network" list on the PC page. The
      right one usually starts with 192.168.

  "No network found"
    - The PC isn't connected to any network. Connect it to Wi-Fi or turn
      on Mobile Hotspot, then reload the page.

  Files picked on the iPhone don't arrive
    - Open "Activity log" at the bottom of the iPhone page (it's also on
      the PC page and in the black window). It shows every file picked,
      sent and saved, and the exact error if something fails.
    - "iPhone could not read this file": Safari lost access to the file.
      Pick it again, using Options -> Current for videos.
    - "not connected - scan the QR code again": the page lost its key.
      Close the tab and scan the QR code again.

  The iPhone says "PC not reachable"
    - The black PhoneDrop window was closed, or the PC went to sleep.
      Start PhoneDrop again; the iPhone reconnects by itself.

  Still stuck? Open an issue in this repo and include a screenshot of
  the Activity log.


---------------------------------------------------------------------
 9. PRIVACY AND SECURITY
---------------------------------------------------------------------
  - Everything stays on your local network. PhoneDrop never uploads
    your files or text anywhere on the internet.
  - The QR link contains a secret key. Only devices that scanned it can
    connect; anyone else on the network is blocked (and shown in the
    Activity log).
  - Only the PC itself can see received files, change shared files or
    see the QR code. The iPhone can only send files and see what you
    chose to share.
  - Connections use plain HTTP inside your network. Use PhoneDrop on
    networks you trust (home Wi-Fi or your own hotspot), not public
    Wi-Fi.
  - To change the secret key (for example, to stop a device that has
    the old link): close PhoneDrop, delete
        %APPDATA%\PhoneDrop\token.txt
    and start PhoneDrop again. Every device then has to scan again.


---------------------------------------------------------------------
 10. WHERE THINGS ARE STORED
---------------------------------------------------------------------
  Received from iPhone   Downloads\PhoneDrop\From iPhone
  Shared to iPhone       Downloads\PhoneDrop\To iPhone
  Secret key             %APPDATA%\PhoneDrop\token.txt

  Use a different folder by starting PhoneDrop with a path, for example
  from Command Prompt:
      PhoneDrop.exe "D:\Shared"
  or make a shortcut to PhoneDrop.exe and add the path after the target.

  PhoneDrop uses port 8080. If that port is busy it tries 8081, 8088,
  8090, 8765 and 9090, then any free port. The address in use is shown
  in the black window.

  Partly-received files are hidden while transferring (".part") and
  cleaned up automatically. Text messages are kept only while PhoneDrop
  is running.


---------------------------------------------------------------------
 11. BUILD FROM SOURCE
---------------------------------------------------------------------
  The full source is in the "source" folder of the zip (and in this
  repo). It is a single Go program with the web pages built in.

  Files
    main.go                  the server
    web\phone.html           iPhone page
    web\pc.html              PC control page
    web\icon.svg, icon.png   icons
    phonedrop.ico            Windows app icon
    rsrc_windows_amd64.syso  icon resource (picked up automatically
                             by Go when building for Windows)
    go.mod, go.sum           dependencies (QR code library)

  Requirements
    Go 1.24 or newer - https://go.dev/dl/

  Build on Windows (Command Prompt, inside the source folder)
    go build -trimpath -ldflags "-s -w" -o PhoneDrop.exe .

  Build the Windows app from macOS or Linux
    GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o PhoneDrop.exe .

  Run on macOS or Linux (works the same; files go to ~/Downloads/PhoneDrop)
    go build -o phonedrop . && ./phonedrop

  Run without opening the browser automatically
    set PHONEDROP_NO_BROWSER=1        (Windows)
    PHONEDROP_NO_BROWSER=1 ./phonedrop (macOS / Linux)

  Changing the icon (optional)
    go install github.com/akavel/rsrc@latest
    rsrc -ico phonedrop.ico -arch amd64 -o rsrc_windows_amd64.syso


---------------------------------------------------------------------
 12. WHAT'S NEW IN 1.4
---------------------------------------------------------------------
  1.4  Big .mov videos no longer get stuck. Each file is sent whole; if
       the connection breaks, the PC keeps what arrived and the iPhone
       continues from there.
  1.3  Stuck transfers are detected and restarted automatically; live
       speed and time left.
  1.2  Fixed picked photos/videos sometimes not arriving; added the
       Activity log.
  1.1  Reliable big uploads with automatic retry; screen stays awake
       while sending.
  1.0  First release.