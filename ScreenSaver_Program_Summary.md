# macOS Video Screensaver Generator

## Overview
This Go program (`main.go`) generates native macOS screensavers from any MP4 video file. It creates complete `.saver` bundles that can be installed and used as system screensavers.

## How It Works

### 1. Program Structure
- **Language**: Go
- **Input**: MP4 video file, output path, screensaver name
- **Output**: Complete macOS screensaver bundle (`.saver` directory)

### 2. Bundle Creation Process

The program creates a standard macOS screensaver bundle structure:
```
YourScreensaver.saver/
├── Contents/
│   ├── Info.plist
│   ├── MacOS/
│   │   └── VideoSaver (compiled Swift shared library)
│   └── Resources/
│       └── payload.mp4 (your video file)
```

### 3. Technical Implementation

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

### 4. Compilation Process
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

### 5. Usage
```bash
go run main.go mac -in <input.mp4> -out <output.saver> [-name "Display Name"]
```

**Example:**
```bash
go run main.go mac -in video.mp4 -out MyScreensaver.saver -name "My Video Screensaver"
```

### 6. Installation
Once generated, the `.saver` bundle can be:
1. Double-clicked to install in System Preferences
2. Manually copied to `~/Library/Screen Savers/` or `/Library/Screen Savers/`
3. Selected from System Preferences > Desktop & Screen Saver

**Important:** Restart your Mac after installation for the screensaver to load properly.

## Features
- ✅ Automatic video looping
- ✅ Full-screen video playback
- ✅ Preview support in System Preferences
- ✅ Muted audio (screensaver appropriate)
- ✅ Responsive scaling for different screen sizes
- ✅ Swift compiler and Xcode build system support
- ✅ Proper shared library compilation

## Technical Requirements
- **macOS**: macOS 10.12 Sierra or later
- **Architecture**: Intel or Apple Silicon
- **Build Tools**: Xcode Command Line Tools (`xcode-select --install`)
- **Dependencies**: Swift compiler, AVFoundation, ScreenSaver frameworks (built-in)

## Limitations
- **macOS Only**: Creates `.saver` bundles specific to macOS
- **MP4 Format**: Optimized for MP4 video input
- **Swift Dependency**: Requires Swift compiler on build machine