# macOS Video Screensaver Generator

A simple Go tool that converts any MP4 video into a native macOS screensaver.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![macOS](https://img.shields.io/badge/macOS-10.12%2B-blue)](https://developer.apple.com/macos/)
[![Swift](https://img.shields.io/badge/Swift-5.0%2B-orange)](https://swift.org)
[![Go](https://img.shields.io/badge/Go-1.19%2B-00ADD8)](https://golang.org)

## What it does

Transform your favorite videos into beautiful, looping macOS screensavers with a single command:

```bash
go run main.go mac -in your-video.mp4 -out MyScreensaver.saver -name "My Custom Screensaver"
```

The tool generates a complete `.saver` bundle that integrates seamlessly with macOS System Preferences, just like built-in screensavers.

## Features

- 🎥 **Native video playback** - Uses macOS AVFoundation for smooth performance
- 🔄 **Automatic looping** - Videos repeat seamlessly 
- 🖥️ **Full-screen support** - Scales to fit any screen size
- 👁️ **Preview support** - Works in System Preferences preview pane
- 🔇 **Muted audio** - Appropriate for screensaver use
- ⚡ **Fast compilation** - Uses Swift compiler or Xcode build system
- 📦 **Self-contained** - Video embedded in the screensaver bundle

## Quick Start

### Prerequisites
- macOS 10.12 Sierra or later
- Xcode Command Line Tools: `xcode-select --install`
- Go 1.19+ (if building from source)

### Generate a Screensaver
```bash
# Clone the repository
git clone https://github.com/LawrenceMarsman/macos-video-screensaver.git
cd macos-video-screensaver

# Generate screensaver from your video
go run main.go mac -in path/to/your/video.mp4 -out MyScreensaver.saver -name "My Video"

# Install by double-clicking MyScreensaver.saver
# Then select it in System Preferences > Desktop & Screen Saver
```

**Important:** Restart your Mac after installation for the screensaver to load properly.

## Examples

```bash
# Nature video screensaver
go run main.go mac -in forest-timelapse.mp4 -out ForestSaver.saver -name "Forest Timelapse"

# Abstract art screensaver  
go run main.go mac -in abstract-loops.mp4 -out AbstractSaver.saver -name "Abstract Art"

# Personal memories screensaver
go run main.go mac -in family-vacation.mp4 -out VacationSaver.saver -name "Family Vacation"
```

## How it Works

1. **Creates bundle structure** - Generates proper macOS `.saver` bundle layout
2. **Embeds video** - Copies your MP4 into the bundle resources
3. **Compiles Swift code** - Builds native screensaver using AVFoundation
4. **Links frameworks** - Integrates with ScreenSaver and system frameworks
5. **Sets permissions** - Ensures proper executable permissions

The result is a native macOS screensaver that appears in System Preferences alongside built-in options.

## Documentation

- **[📖 Installation & Usage Guide](INSTALLATION.md)** - Detailed setup instructions and troubleshooting
- **[🔧 Technical Details](INSTALLATION.md#how-it-works)** - Implementation specifics and customization options

## Supported Formats

- **Input**: MP4 video files (recommended for best compatibility)
- **Output**: Native macOS `.saver` bundles
- **Platforms**: macOS 10.12+ (Intel and Apple Silicon)

## Related Projects

- **Windows Users**: Check out [EasyVideoScreensaver](https://github.com/tonyfederer/EasyVideoScreensaver) for creating video screensavers on Windows

## Contributing

Contributions are welcome! Please feel free to submit issues, feature requests, or pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with Go for cross-platform development
- Uses Swift and AVFoundation for native macOS video playback
- Integrates with macOS ScreenSaver framework for system compatibility

---

**Transform your videos into screensavers in seconds!** ✨