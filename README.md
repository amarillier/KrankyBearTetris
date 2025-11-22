# KrankyBear Tetris

A cross-platform Tetris game built with Go and the Fyne GUI library. Enjoy classic Tetris gameplay with modern graphics and sound effects.

## Features

- **Cross-Platform Support**: Works on desktop and mobile platforms
  - **Desktop Platforms**:
    - **Linux**: Works on all desktop environments (GNOME, KDE, XFCE, Cinnamon, MATE, etc.) with X11 or Wayland
    - **macOS**: macOS 10.13 (High Sierra) or later
    - **Windows**: Windows 10 or later
  - **Mobile Platforms**:
    - **iOS**: iOS 13.0 or later (requires sideloading)
    - **Android**: Android 5.0 (API level 21) or later

# Features & Known Issues

### 1. **Fyne GUI** (Primary - Requires OpenGL)
- Modern interface
- Custom icons support
- Auto-close with timeout
- Full feature set

### 2. **WebView GUI** (Optional Fallback - HTML/CSS/JavaScript)
- Used when OpenGL is not available (if compiled with `-tags webview`)
- Web based gradient UI with animations
- Auto-close with countdown timer
- Works in VMs and Remote Desktop
- **Requires**: Build with `-tags webview` flag
- **Runtime**: WebView2 on Windows, WebKit on macOS/Linux


## Building

### Desktop Platforms

Build scripts are provided for each platform:
- `./compile-mac.sh` - Build for macOS (Intel and Apple Silicon)
- `./compile-linux.sh` - Build for Linux (amd64 and arm64)
- `./compile-windows.ps1` - Build for Windows (requires Windows or WSL)
- `./compile-all.sh` - Build and package for all desktop platforms

Or use the Makefile:
```bash
make build-all          # Build for all desktop platforms
make build-darwin       # Build for macOS
make build-linux        # Build for Linux (requires fyne-cross or native Linux)
make build-windows      # Build for Windows (requires mingw-w64)
```

### Mobile Platforms

#### Prerequisites

**For iOS:**
- macOS with Xcode installed
- iOS SDK (included with Xcode)
- Apple Developer account (for device installation)

**For Android:**
- Android SDK and NDK
- Set `ANDROID_NDK_HOME` environment variable to your NDK path
- Android device with USB debugging enabled (for sideloading)

#### Building Mobile Apps

1. Install the Fyne command-line tool:
   ```bash
   go install fyne.io/fyne/v2/cmd/fyne@latest
   ```

2. Build for iOS:
   ```bash
   fyne package -os ios -appID com.github.amarillier.KrankyBearTetris -name "KrankyBear Tetris"
   ```
   This creates `KrankyBearTetris.app` bundle.

3. Build for Android:
   ```bash
   fyne package -os android -appID com.github.amarillier.KrankyBearTetris -name "KrankyBear Tetris"
   ```
   This creates `KrankyBearTetris.apk` file.

Or use the Makefile:
```bash
make build-ios      # Build for iOS (macOS only)
make build-android  # Build for Android
make build-all-platforms  # Build for all platforms including mobile
```

## Sideloading Mobile Apps

### iOS Sideloading

**Option 1: Using Xcode (Recommended for Development)**
1. Connect your iOS device to your Mac via USB
2. Open Xcode
3. Go to **Window** → **Devices and Simulators**
4. Select your connected device
5. Drag the `KrankyBearTetris.app` bundle onto the "Installed Apps" list
6. Trust the developer certificate on your device: **Settings** → **General** → **Device Management** → Trust your developer account

**Option 2: Using AltStore (No Developer Account Required)**
1. Install [AltStore](https://altstore.io/) on your iOS device
2. Transfer the `.ipa` file (convert `.app` to `.ipa` if needed) to your device
3. Open AltStore and install the app
4. Refresh the app weekly using AltStore (free account limitation)

**Option 3: Using Sideloadly (Windows/Mac)**
1. Download [Sideloadly](https://sideloadly.io/)
2. Connect your iOS device
3. Drag the `.ipa` file into Sideloadly
4. Enter your Apple ID (non-developer account works)
5. Click "Start" to install

**Note:** Apps installed via sideloading expire after 7 days (free Apple ID) or 1 year (paid developer account). You'll need to reinstall periodically.

### Android Sideloading

**Option 1: Using ADB (Recommended for Development)**
1. Enable **Developer Options** on your Android device:
   - Go to **Settings** → **About Phone**
   - Tap **Build Number** 7 times
2. Enable **USB Debugging**:
   - Go to **Settings** → **Developer Options**
   - Enable **USB Debugging**
3. Connect your device to your computer via USB
4. Install the APK:
   ```bash
   adb install KrankyBearTetris.apk
   ```

**Option 2: Direct Installation**
1. Transfer the `KrankyBearTetris.apk` file to your Android device
2. On your device, go to **Settings** → **Security** → Enable **Unknown Sources** (or **Install Unknown Apps** on newer Android versions)
3. Open the APK file using a file manager
4. Tap **Install** when prompted

**Option 3: Using Wireless ADB**
1. Connect your device via USB initially
2. Enable wireless debugging:
   ```bash
   adb tcpip 5555
   adb connect <device-ip>:5555
   ```
3. Disconnect USB and install wirelessly:
   ```bash
   adb install KrankyBearTetris.apk
   ```

**Security Note:** Sideloading apps from unknown sources can pose security risks. Only install apps from trusted sources.

## Dependencies

- [Fyne](https://fyne.io/) v2.7.1 - Cross-platform GUI toolkit
- Go 1.21 or later
- For mobile builds: Fyne command-line tool (`go install fyne.io/fyne/v2/cmd/fyne@latest`)

## License

This project is provided as-is, free for personal, educational and commercial use, under GNU GPL-3.0

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## Author

Allan Marillier

## Acknowledgments

- Built with [Fyne](https://fyne.io/) - An easy-to-use GUI toolkit for Go
- Classic Tetris gameplay with modern cross-platform support
