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

**Installation Steps for iOS:**

1. **Install Xcode** (required):
   ```bash
   # Install Xcode from Mac App Store, or via Homebrew:
   brew install --cask xcode
   ```

2. **Install Xcode Command Line Tools** (if not already installed):
   ```bash
   xcode-select --install
   ```

3. **Accept Xcode License**:
   ```bash
   sudo xcodebuild -license accept
   ```

4. **Install Go** (if not already installed):
   ```bash
   brew install go
   ```

5. **Install Fyne Command-Line Tool**:
   ```bash
   go install fyne.io/fyne/v2/cmd/fyne@latest
   ```

6. **Verify Installation**:
   ```bash
   xcodebuild -version
   go version
   fyne version
   ```

**For Android:**
- Android SDK and NDK
- Set `ANDROID_NDK_HOME` environment variable to your NDK path
- Android device with USB debugging enabled (for sideloading)

**Installation Steps for Android:**

1. **Install Homebrew** (if not already installed):
   ```bash
   /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
   ```

2. **Install Java Development Kit (JDK)**:
   ```bash
   brew install openjdk@17
   # Add to your shell profile (~/.zshrc or ~/.bash_profile):
   echo 'export PATH="/opt/homebrew/opt/openjdk@17/bin:$PATH"' >> ~/.zshrc
   echo 'export JAVA_HOME="/opt/homebrew/opt/openjdk@17"' >> ~/.zshrc
   source ~/.zshrc
   ```

3. **Install Android Command Line Tools**:
   ```bash
   # Create Android SDK directory
   mkdir -p ~/Library/Android/sdk
   cd ~/Library/Android/sdk
   
   # Download command line tools (replace with latest version from https://developer.android.com/studio#command-tools)
   curl -o cmdline-tools.zip https://dl.google.com/android/repository/commandlinetools-mac-11076708_latest.zip
   unzip cmdline-tools.zip
   mkdir -p cmdline-tools/latest
   mv cmdline-tools/* cmdline-tools/latest/ 2>/dev/null || true
   rm cmdline-tools.zip
   ```

4. **Set Android Environment Variables**:
   ```bash
   # Add to your shell profile (~/.zshrc or ~/.bash_profile):
   cat >> ~/.zshrc << 'EOF'
   export ANDROID_HOME="$HOME/Library/Android/sdk"
   export PATH="$ANDROID_HOME/cmdline-tools/latest/bin:$PATH"
   export PATH="$ANDROID_HOME/platform-tools:$PATH"
   export PATH="$ANDROID_HOME/tools:$PATH"
   export PATH="$ANDROID_HOME/tools/bin:$PATH"
   EOF
   source ~/.zshrc
   ```

5. **Install Android SDK Components**:
   ```bash
   # Accept licenses
   yes | sdkmanager --licenses
   
   # Install required SDK components
   sdkmanager "platform-tools" "platforms;android-33" "build-tools;33.0.0"
   
   # Install NDK (required for Fyne Android builds)
   sdkmanager "ndk;25.2.9519653"
   ```

6. **Set NDK Environment Variable**:
   ```bash
   # Add to your shell profile (~/.zshrc or ~/.bash_profile):
   echo 'export ANDROID_NDK_HOME="$HOME/Library/Android/sdk/ndk/25.2.9519653"' >> ~/.zshrc
   source ~/.zshrc
   ```

7. **Install Go** (if not already installed):
   ```bash
   brew install go
   ```

8. **Install Fyne Command-Line Tool**:
   ```bash
   go install fyne.io/fyne/v2/cmd/fyne@latest
   ```

9. **Install ADB (Android Debug Bridge)** for sideloading:
   ```bash
   # ADB is included in platform-tools, but you can also install via Homebrew:
   brew install android-platform-tools
   ```

10. **Verify Installation**:
    ```bash
    java -version
    echo $ANDROID_HOME
    echo $ANDROID_NDK_HOME
    adb version
    go version
    fyne version
    ```

**Alternative: Install Android Studio (GUI Method)**

If you prefer a GUI, you can install Android Studio which includes all necessary tools:

```bash
brew install --cask android-studio
```

After installation:
1. Open Android Studio
2. Go to **Tools** → **SDK Manager**
3. Install **Android SDK Platform 33** (or latest)
4. Install **Android SDK Build-Tools**
5. Install **NDK (Side by side)** - version 25.2.9519653 or later
6. Set environment variables as shown above (Android Studio sets `ANDROID_HOME` automatically)

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

### Required Dependencies

- **Go** 1.21 or later - Programming language
- **[Fyne](https://fyne.io/)** v2.7.1 - Cross-platform GUI toolkit
- **Fyne command-line tool** - For mobile builds (`go install fyne.io/fyne/v2/cmd/fyne@latest`)

### Installation

#### macOS (using Homebrew)

1. **Install Homebrew** (if not already installed):
   ```bash
   /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
   ```

2. **Install Go**:
   ```bash
   brew install go
   ```

3. **Verify Go Installation**:
   ```bash
   go version
   # Should show: go version go1.21.x or later
   ```

4. **Install Project Dependencies**:
   ```bash
   # Clone or navigate to the project directory
   cd /path/to/KrankyBearTetris
   
   # Download and install all Go dependencies
   go mod download
   go mod tidy
   ```

5. **Install Fyne Command-Line Tool** (for mobile builds):
   ```bash
   go install fyne.io/fyne/v2/cmd/fyne@latest
   ```

6. **Verify Fyne Installation**:
   ```bash
   fyne version
   ```

#### Linux (Ubuntu/Debian)

1. **Install Go**:
   ```bash
   # Using apt (Ubuntu/Debian)
   sudo apt-get update
   sudo apt-get install -y golang-go
   
   # Or install latest version from official Go website:
   # https://go.dev/dl/
   ```

2. **Set Go Environment Variables** (if needed):
   ```bash
   export GOPATH=$HOME/go
   export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
   ```

3. **Install Project Dependencies**:
   ```bash
   cd /path/to/KrankyBearTetris
   go mod download
   go mod tidy
   ```

4. **Install Fyne Command-Line Tool**:
   ```bash
   go install fyne.io/fyne/v2/cmd/fyne@latest
   ```

5. **Install System Dependencies for Fyne** (Linux):
   ```bash
   # Ubuntu/Debian
   sudo apt-get install -y \
     libgl1-mesa-dev \
     xorg-dev \
     libxcursor-dev \
     libxrandr-dev \
     libxinerama-dev \
     libxi-dev \
     libglfw3-dev \
     libx11-dev \
     libx11-xcb-dev \
     libxkbcommon-x11-dev \
     libxkbcommon-dev \
     libxcb1-dev \
     libxcb-keysyms1-dev \
     libxcb-icccm4-dev \
     libxcb-image0-dev \
     libxcb-shm0-dev \
     libxcb-util1-dev \
     libxcb-render0-dev \
     libxcb-render-util0-dev \
     libxcb-xfixes0-dev \
     libxcb-xinerama0-dev \
     libxcb-dri3-dev \
     libasound2-dev \
     libpulse-dev
   ```

#### Windows

1. **Install Go**:
   - Download from: https://go.dev/dl/
   - Run the installer and follow the prompts
   - Verify installation:
     ```powershell
     go version
     ```

2. **Install Project Dependencies**:
   ```powershell
   cd C:\path\to\KrankyBearTetris
   go mod download
   go mod tidy
   ```

3. **Install Fyne Command-Line Tool**:
   ```powershell
   go install fyne.io/fyne/v2/cmd/fyne@latest
   ```

4. **Install System Dependencies for Fyne** (Windows):
   - Fyne on Windows uses OpenGL which is typically pre-installed
   - For development, you may need:
     - **MinGW-w64** (for cross-compilation from other platforms):
       ```bash
       # On macOS with Homebrew:
       brew install mingw-w64
       ```

### Go Module Dependencies

The project uses the following main dependencies (automatically managed by Go modules):

- `fyne.io/fyne/v2` - GUI toolkit
- `github.com/amarillier/go-update-checker` - Version update checking
- `github.com/gopxl/beep/v2` - Audio/sound effects

All dependencies are automatically downloaded when you run:
```bash
go mod download
```

Or when building:
```bash
go build
```

### Troubleshooting

**Issue: `go: command not found`**
- Solution: Install Go and ensure it's in your PATH
- Verify: `which go` or `go version`

**Issue: `fyne: command not found`**
- Solution: Install Fyne CLI tool: `go install fyne.io/fyne/v2/cmd/fyne@latest`
- Ensure `$GOPATH/bin` or `$HOME/go/bin` is in your PATH

**Issue: Build fails with missing C dependencies (Linux)**
- Solution: Install system dependencies listed above
- Or use the provided `compile-linux.sh` script which auto-installs dependencies

**Issue: `go mod download` fails**
- Solution: Check internet connection
- Try: `go env -w GOPROXY=https://proxy.golang.org,direct`
- Or: `go env -w GOSUMDB=sum.golang.org`

## License

This project is provided as-is, free for personal, educational and commercial use, under GNU GPL-3.0

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## Author

Allan Marillier

## Acknowledgments

- Built with [Fyne](https://fyne.io/) - An easy-to-use GUI toolkit for Go
- Classic Tetris gameplay with modern cross-platform support
