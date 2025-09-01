# Installation & Usage Guide

## Overview
This Go program generates native macOS screensavers from any MP4 video file. It creates complete `.saver` bundles that can be installed and used as system screensavers.

## Quick Start

### Prerequisites
- **macOS**: macOS 10.12 Sierra or later
- **Architecture**: Intel or Apple Silicon
- **Build Tools**: Xcode Command Line Tools (`xcode-select --install`)

### Usage
```bash
go run main.go mac -in video.mp4 -out MyScreensaver.saver -name "My Screensaver"
```

**Example:**
```bash
go run main.go mac -in sunset.mp4 -out SunsetSaver.saver -name "Beautiful Sunset"
```

## Installation

### Method 1: Double-Click Installation (Recommended)
1. **Double-click** the generated `.saver` bundle
2. System Preferences will open automatically
3. Click **"Install"** when prompted
4. The screensaver will be added to your system

### Method 2: Manual Installation
1. **Copy** the `.saver` bundle to one of these locations:
   - `~/Library/Screen Savers/` (for current user only)
   - `/Library/Screen Savers/` (for all users - requires admin rights)

2. **Open System Preferences**:
   - Apple menu → System Preferences → Desktop & Screen Saver
   - OR: Apple menu → System Settings → Wallpaper (macOS 13+)

3. **Select the screensaver**:
   - Click on "Screen Saver" tab
   - Find your screensaver in the list
   - Click to select it

4. **Configure settings** (optional):
   - Set activation time
   - Choose display options
   - Test with "Preview" button

### ⚠️ Important: Restart Required
**After installing the screensaver, you must restart your Mac for it to work properly.** This is a macOS system requirement for loading new screensaver bundles. Without a restart, you may see colored screens instead of your video.

## Troubleshooting

- **Screensaver shows colored screens instead of video**: **Restart your Mac** - this fixes the issue in most cases
- **Security warning**: Go to System Preferences → Security & Privacy → General, and click "Allow" if blocked
- **Doesn't appear**: Try logging out and back in, or restart your Mac
- **Black screen**: Check that the video file is properly embedded in the bundle

## How It Works

### Bundle Structure
The program creates a standard macOS screensaver bundle:
```
YourScreensaver.saver/
├── Contents/
│   ├── Info.plist
│   ├── MacOS/
│   │   └── VideoSaver (compiled Swift shared library)
│   └── Resources/
│       └── payload.mp4 (your video file)
```

### Technical Implementation

#### Swift Screensaver Code
- Uses `AVFoundation` framework for native video playback
- Inherits from `ScreenSaverView` (macOS screensaver base class)
- Features:
  - Automatic video looping
  - Aspect-fill scaling to fit screen
  - Muted playback
  - Black background fallback
  - Compatible with both preview and full-screen modes

#### Bundle Configuration
- **Info.plist**: Contains screensaver metadata and configuration
- **CFBundleExecutable**: Points to the compiled Swift shared library
- **NSPrincipalClass**: Specifies the main screensaver class (`VideoSaverView`)

#### Compilation Process
1. Creates proper bundle directory structure
2. Copies input MP4 to `Resources/payload.mp4`
3. Generates Swift source code with video player logic
4. Compiles Swift code to shared library using:
   ```bash
   swiftc -framework ScreenSaver -framework AVFoundation -framework Cocoa \
          -emit-library -module-name VideoSaver -o VideoSaver VideoSaver.swift
   ```
5. Creates Info.plist with proper metadata
6. Sets executable permissions on the shared library

## Customization

### Replacing the Video
To use a different video in an existing screensaver:

1. **Right-click** on the `.saver` bundle → Show Package Contents
2. **Navigate** to `Contents/Resources/`
3. **Replace** `payload.mp4` with your video (keep the same filename)
4. **Reinstall** the screensaver

### Creating New Screensavers
Simply run the generator with different input files:

```bash
go run main.go mac -in sunset.mp4 -out SunsetSaver.saver -name "Sunset"
go run main.go mac -in ocean.mp4 -out OceanSaver.saver -name "Ocean Waves"
```

## Uninstallation

1. **Delete** the `.saver` file from:
   - `~/Library/Screen Savers/` or `/Library/Screen Savers/`
2. **Restart** System Preferences if open
3. **Select** a different screensaver in System Preferences

## Features
- ✅ Automatic video looping
- ✅ Full-screen video playback
- ✅ Preview support in System Preferences
- ✅ Muted audio (screensaver appropriate)
- ✅ Responsive scaling for different screen sizes
- ✅ Swift compiler and Xcode build system support
- ✅ Proper shared library compilation

## Limitations
- **macOS Only**: Creates `.saver` bundles specific to macOS
- **MP4 Format**: Optimized for MP4 video input
- **Swift Dependency**: Requires Swift compiler on build machine

## Support

For issues or questions:
1. Ensure your video file is a valid MP4 format
2. Verify you have Xcode Command Line Tools installed
3. Check that you have proper permissions to install screensavers  
4. Try restarting your system after installation
5. Ensure your system meets the technical requirements above

---

*Generated using Go and Swift for native macOS screensaver creation*