# KrankyBear Tetris

A cross-platform notification application built with Go and the Fyne GUI library. This application ...

## Features

- **Cross-Platform Support**: Works on macOS, Windows, and Linux
  - **Linux**: Works on all desktop environments (GNOME, KDE, XFCE, Cinnamon, MATE, etc.) with X11 or Wayland
  - **macOS**: macOS 10.13 (High Sierra) or later
  - **Windows**: Windows 10 or later

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


## Dependencies

- [Fyne](https://fyne.io/) v2.6.3 - Cross-platform GUI toolkit
- Go standard library

## License

This project is provided as-is, free for personal, educational and commercial use, under GNU GPL-3.0

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## Author

Allan Marillier

## Acknowledgments

- Built with [Fyne](https://fyne.io/) - An easy-to-use GUI toolkit for Go
- Inspired by the need for simple, cross-platform notification systems
