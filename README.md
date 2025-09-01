# macOS Video Screensaver Generator

A Go program that converts any MP4 video into a native macOS screensaver bundle.

## Usage

```bash
go run main.go mac -in video.mp4 -out MyScreensaver.saver -name "My Screensaver"
```

This creates a complete `.saver` bundle that can be installed on macOS.

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

## Technical Requirements

- **OS**: macOS 10.12 Sierra or later
- **Architecture**: Intel or Apple Silicon
- **Build Tools**: Xcode Command Line Tools (`xcode-select --install`)
- **Dependencies**: Swift compiler, AVFoundation, ScreenSaver frameworks (built-in)

## How It Works

The program creates a complete macOS screensaver bundle with:
- Swift-based video player using AVFoundation
- Proper screensaver inheritance from `ScreenSaverView`
- Embedded MP4 video file in the bundle resources
- Automatic video looping and aspect-fill scaling
- Muted audio playback (screensaver appropriate)

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

## Project Structure

- `main.go` - Primary screensaver generator
- `go.mod` - Go module definition
- `README.md` - This file

## Support

For issues or questions:
1. Ensure your video file is a valid MP4 format
2. Verify you have Xcode Command Line Tools installed
3. Check that you have proper permissions to install screensavers  
4. Try restarting your system after installation
5. Ensure your system meets the technical requirements above

---

*Generated using Go and Swift for native macOS screensaver creation*