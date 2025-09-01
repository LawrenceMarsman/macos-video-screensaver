// macOS Video Screensaver Generator
// A Go CLI that converts MP4 videos into native macOS .saver bundles.
//
// Creates a complete Swift-based screensaver that:
// - Uses AVFoundation for native video playback
// - Inherits from ScreenSaverView for proper system integration
// - Embeds the MP4 video file in the bundle resources
// - Supports both preview and full-screen modes
// - Automatically loops video with muted audio
//
// Usage:
//   go run main.go mac -in video.mp4 -out MyScreensaver.saver [-name "My Screensaver"]
//
// Requirements:
// - macOS with Xcode Command Line Tools
// - Swift compiler (swiftc) or Xcode
// - Valid MP4 input file
//
// Output: Complete .saver bundle ready for installation

package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func main() {
    if len(os.Args) < 2 {
        usage()
        os.Exit(2)
    }
    cmd := os.Args[1]
    flagSet := flag.NewFlagSet(cmd, flag.ExitOnError)
    in := flagSet.String("in", "", "input MP4 file")
    out := flagSet.String("out", "", "output .saver bundle")
    name := flagSet.String("name", "MyScreensaver", "screensaver display name")
    if err := flagSet.Parse(os.Args[2:]); err != nil {
        fatal(err)
    }
    if *in == "" || *out == "" {
        usage()
        os.Exit(2)
    }

    switch cmd {
    case "mac":
        if err := buildMacSaver(*in, *out, *name); err != nil {
            fatal(err)
        }
        fmt.Println("✅ Built macOS screensaver:", *out)
    default:
        usage()
        os.Exit(2)
    }
}

func usage() {
    fmt.Println(`macOS Video Screensaver Generator - convert MP4 to macOS .saver bundle

Usage:
  go run main.go mac -in video.mp4 -out MyScreensaver.saver [-name "My Screensaver"]

Example:
  go run main.go mac -in sunset.mp4 -out SunsetSaver.saver -name "Beautiful Sunset"`)
}

func fatal(err error) {
    fmt.Fprintln(os.Stderr, "ERROR:", err)
    os.Exit(1)
}

func sanitizeName(s string) string {
    s = strings.TrimSpace(s)
    if s == "" { return "Screensaver" }
    return s
}

// ---------------- macOS (.saver) ----------------

func buildMacSaver(in, out, name string) error {
    if runtime.GOOS != "darwin" {
        fmt.Println("[warn] Building a macOS .saver requires macOS.")
    }
    
    // Try alternative approaches before requiring Xcode
    if _, err := exec.LookPath("swiftc"); err == nil {
        return buildMacSaverSwift(in, out, name)
    }
    
    // Fallback to Xcode if available
    if _, err := exec.LookPath("xcodebuild"); err == nil {
        return buildMacSaverXcode(in, out, name)
    }
    
    return errors.New("Neither swiftc nor xcodebuild found. Install Xcode command line tools: xcode-select --install")
}

func buildMacSaverSwift(in, out, name string) error {
    fmt.Println("[info] Using Swift compiler directly")
    
    tempDir, err := os.MkdirTemp("", "scrgen-mac-*")
    if err != nil { return err }
    // Don't remove tempDir immediately to allow debugging
    fmt.Printf("[debug] Temp directory: %s\n", tempDir)

    projName := sanitizeName(name)
    
    // Create bundle structure manually
    bundlePath := filepath.Join(tempDir, projName+".saver")
    contentsPath := filepath.Join(bundlePath, "Contents")
    macosPath := filepath.Join(contentsPath, "MacOS")
    resourcesPath := filepath.Join(contentsPath, "Resources")
    
    for _, dir := range []string{bundlePath, contentsPath, macosPath, resourcesPath} {
        if err := os.MkdirAll(dir, 0755); err != nil { return err }
    }
    
    // Copy video
    if err := copyFile(in, filepath.Join(resourcesPath, "payload.mp4")); err != nil {
        return err
    }
    
    // Create Info.plist
    plistPath := filepath.Join(contentsPath, "Info.plist")
    if err := os.WriteFile(plistPath, []byte(infoPlist(projName)), 0644); err != nil {
        return err
    }
    
    // Create Swift source
    swiftPath := filepath.Join(tempDir, "VideoSaver.swift")
    if err := os.WriteFile(swiftPath, []byte(swiftSaverClass(projName)), 0644); err != nil {
        return err
    }
    
    // Compile with swiftc as a shared library for screensavers
    execPath := filepath.Join(macosPath, "VideoSaver")
    fmt.Printf("[debug] Compiling Swift to: %s\n", execPath)
    if err := run(tempDir, "swiftc", 
        "-framework", "ScreenSaver",
        "-framework", "AVFoundation", 
        "-framework", "AVKit",
        "-framework", "Cocoa",
        "-emit-library",
        "-module-name", "VideoSaver",
        "-o", execPath,
        "VideoSaver.swift"); err != nil {
        return fmt.Errorf("swift compilation failed: %w", err)
    }
    
    // Verify executable was created
    if _, err := os.Stat(execPath); err != nil {
        return fmt.Errorf("executable not created: %s", execPath)
    }
    fmt.Printf("[debug] Executable created successfully: %s\n", execPath)
    
    // Copy bundle to output
    if err := copyDir(bundlePath, out); err != nil { return err }
    
    // Ensure executable permissions on the final output
    finalExecPath := filepath.Join(out, "Contents/MacOS/VideoSaver")
    if err := os.Chmod(finalExecPath, 0755); err != nil {
        fmt.Printf("[warn] Could not set executable permissions: %v\n", err)
    }
    
    return nil
}

func buildMacSaverXcode(in, out, name string) error {
    fmt.Println("[info] Using Xcode build system")

    tempDir, err := os.MkdirTemp("", "scrgen-mac-*")
    if err != nil { return err }
    // Do not defer RemoveAll; leave for inspection on failures

    projName := sanitizeName(name)
    swiftFiles := map[string]string{
        filepath.Join(tempDir, "VideoSaver.xcodeproj/project.pbxproj"): xcodeprojPbxproj(projName),
        filepath.Join(tempDir, "VideoSaver/VideoSaver.swift"):          swiftSaverClass(projName),
        filepath.Join(tempDir, "VideoSaver/Info.plist"):               infoPlist(projName),
    }
    for p, content := range swiftFiles {
        if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil { return err }
        if err := os.WriteFile(p, []byte(content), 0644); err != nil { return err }
    }
    // Copy MP4 into the target bundle resources folder (referenced by project file)
    payloadDest := filepath.Join(tempDir, "VideoSaver/payload.mp4")
    if err := copyFile(in, payloadDest); err != nil { return err }

    // Build .saver
    if err := run(tempDir, "xcodebuild", "-project", "VideoSaver.xcodeproj", "-scheme", "VideoSaver", "-configuration", "Release", "build"); err != nil {
        return fmt.Errorf("xcodebuild failed: %w", err)
    }

    // Find built .saver in build/Release
    built := filepath.Join(tempDir, "build/Release/"+projName+".saver")
    if _, err := os.Stat(built); err != nil {
        return fmt.Errorf("expected output not found: %s", built)
    }
    // Copy to desired out path
    if err := copyDir(built, out); err != nil { return err }

    // Cleanup temp dir after a short delay (to allow inspection if needed)
    go func(dir string) {
        time.Sleep(3 * time.Second)
        os.RemoveAll(dir)
    }(tempDir)

    return nil
}

func copyFile(src, dst string) error {
    in, err := os.Open(src)
    if err != nil { return err }
    defer in.Close()
    if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil { return err }
    out, err := os.Create(dst)
    if err != nil { return err }
    defer out.Close()
    if _, err := io.Copy(out, in); err != nil { return err }
    return out.Close()
}

func copyDir(src, dst string) error {
    info, err := os.Stat(src)
    if err != nil { return err }
    if info.IsDir() {
        return copyDirRecursive(src, dst)
    }
    return copyFile(src, dst)
}

func copyDirRecursive(src, dst string) error {
    return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
        if err != nil { return err }
        rel, _ := filepath.Rel(src, path)
        target := filepath.Join(dst, rel)
        if d.IsDir() {
            return os.MkdirAll(target, 0755)
        }
        return copyFile(path, target)
    })
}

func run(dir string, name string, args ...string) error {
    cmd := exec.Command(name, args...)
    cmd.Dir = dir
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    return cmd.Run()
}

// ---------- Xcode project templates ----------

func xcodeprojPbxproj(name string) string {
    // Minimal pbxproj for a Screen Saver target named `name`.
    // To keep this file compact, we embed a pre-made pbxproj zipped and substitute name, but here we
    // generate a very small one inline.
    // For brevity and reliability, we use a single target with sources: VideoSaver.swift, Info.plist, payload.mp4

    // This pbxproj is simplified and works with modern Xcode. If Xcode changes formats, you may need to
    // refresh it. The UUIDs are fixed for simplicity.
    return strings.ReplaceAll(`// !$*UTF8*$!
{
  archiveVersion = 1;
  classes = {};
  objectVersion = 56;
  objects = {

/* Begin PBXFileReference section */
    000000000000000000000001 /* VideoSaver.swift */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.swift; path = VideoSaver.swift; sourceTree = "<group>"; };
    000000000000000000000002 /* Info.plist */ = {isa = PBXFileReference; lastKnownFileType = text.plist.xml; path = Info.plist; sourceTree = "<group>"; };
    000000000000000000000003 /* payload.mp4 */ = {isa = PBXFileReference; lastKnownFileType = file; path = payload.mp4; sourceTree = "<group>"; };
    000000000000000000000010 /* `+name+`.saver */ = {isa = PBXFileReference; explicitFileType = wrapper.cfbundle; includeInIndex = 0; path = "`+name+`.saver"; sourceTree = BUILT_PRODUCTS_DIR; };
/* End PBXFileReference section */

/* Begin PBXGroup section */
    000000000000000000000100 = {isa = PBXGroup; children = (
            000000000000000000000200 /* VideoSaver */,
            000000000000000000000300 /* Products */,
        ); sourceTree = "<group>"; };
    000000000000000000000200 /* VideoSaver */ = {isa = PBXGroup; children = (
            000000000000000000000001 /* VideoSaver.swift */,
            000000000000000000000002 /* Info.plist */,
            000000000000000000000003 /* payload.mp4 */,
        ); path = VideoSaver; sourceTree = "<group>"; };
    000000000000000000000300 /* Products */ = {isa = PBXGroup; children = (
            000000000000000000000010 /* `+name+`.saver */,
        ); name = Products; sourceTree = "<group>"; };
/* End PBXGroup section */

/* Begin PBXNativeTarget section */
    000000000000000000000400 /* VideoSaver */ = {isa = PBXNativeTarget; buildConfigurationList = 000000000000000000000800 /* Build configuration list for PBXNativeTarget "VideoSaver" */; buildPhases = (
            000000000000000000000500 /* Sources */,
            000000000000000000000600 /* Resources */,
        ); buildRules = ( ); dependencies = ( ); name = VideoSaver; productName = VideoSaver; productReference = 000000000000000000000010 /* `+name+`.saver */; productType = "com.apple.product-type.bundle"; };
/* End PBXNativeTarget section */

/* Begin PBXProject section */
    000000000000000000000700 /* Project object */ = {isa = PBXProject; buildConfigurationList = 000000000000000000000900 /* Build configuration list for PBXProject "VideoSaver" */; compatibilityVersion = "Xcode 14.0"; developmentRegion = en; hasScannedForEncodings = 0; knownRegions = (en); mainGroup = 000000000000000000000100; productRefGroup = 000000000000000000000300 /* Products */; projectDirPath = ""; projectRoot = ""; targets = (000000000000000000000400 /* VideoSaver */); };
/* End PBXProject section */

/* Begin PBXResourcesBuildPhase section */
    000000000000000000000600 /* Resources */ = {isa = PBXResourcesBuildPhase; files = (
            000000000000000000000604 /* payload.mp4 in Resources */,
            000000000000000000000603 /* Info.plist in Resources */,
        ); };
/* End PBXResourcesBuildPhase section */

/* Begin PBXSourcesBuildPhase section */
    000000000000000000000500 /* Sources */ = {isa = PBXSourcesBuildPhase; files = (
            000000000000000000000501 /* VideoSaver.swift in Sources */,
        ); };
/* End PBXSourcesBuildPhase section */

/* Begin PBXBuildFile section */
    000000000000000000000501 /* VideoSaver.swift in Sources */ = {isa = PBXBuildFile; fileRef = 000000000000000000000001 /* VideoSaver.swift */; };
    000000000000000000000603 /* Info.plist in Resources */ = {isa = PBXBuildFile; fileRef = 000000000000000000000002 /* Info.plist */; };
    000000000000000000000604 /* payload.mp4 in Resources */ = {isa = PBXBuildFile; fileRef = 000000000000000000000003 /* payload.mp4 */; };
/* End PBXBuildFile section */

/* Begin XCBuildConfiguration section */
    000000000000000000000901 /* Debug */ = {isa = XCBuildConfiguration; buildSettings = {
        PRODUCT_NAME = "`+name+`";
        INFOPLIST_FILE = VideoSaver/Info.plist;
        WRAPPER_EXTENSION = saver;
        CODE_SIGNING_ALLOWED = NO;
        CODE_SIGNING_REQUIRED = NO;
        MACOSX_DEPLOYMENT_TARGET = 11.0;
        SWIFT_VERSION = 5.0;
    }; name = Debug; };
    000000000000000000000902 /* Release */ = {isa = XCBuildConfiguration; buildSettings = {
        PRODUCT_NAME = "`+name+`";
        INFOPLIST_FILE = VideoSaver/Info.plist;
        WRAPPER_EXTENSION = saver;
        CODE_SIGNING_ALLOWED = NO;
        CODE_SIGNING_REQUIRED = NO;
        MACOSX_DEPLOYMENT_TARGET = 11.0;
        SWIFT_VERSION = 5.0;
    }; name = Release; };
/* End XCBuildConfiguration section */

/* Begin XCConfigurationList section */
    000000000000000000000800 /* Build configuration list for PBXNativeTarget "VideoSaver" */ = {isa = XCConfigurationList; buildConfigurations = (
            000000000000000000000901 /* Debug */,
            000000000000000000000902 /* Release */,
        ); defaultConfigurationIsVisible = 0; defaultConfigurationName = Release; };
    000000000000000000000900 /* Build configuration list for PBXProject "VideoSaver" */ = {isa = XCConfigurationList; buildConfigurations = (
            000000000000000000000901 /* Debug */,
            000000000000000000000902 /* Release */,
        ); defaultConfigurationIsVisible = 0; defaultConfigurationName = Release; };
/* End XCConfigurationList section */

  };
  rootObject = 000000000000000000000700 /* Project object */;
}
`, "`+name+`", name)
}

func swiftSaverClass(_ string) string {
    return `import ScreenSaver
import AVFoundation
import Cocoa

@objc(VideoSaverView)
public class VideoSaverView: ScreenSaverView {
    var player: AVPlayer?
    var playerLayer: AVPlayerLayer?
    
    public override init?(frame: NSRect, isPreview: Bool) {
        super.init(frame: frame, isPreview: isPreview)
        setupPlayer()
    }
    
    required init?(coder: NSCoder) {
        super.init(coder: coder)
        setupPlayer()
    }
    
    func setupPlayer() {
        self.wantsLayer = true
        self.layer = CALayer()
        self.layer?.backgroundColor = NSColor.black.cgColor
        
        // Try to find the video file
        guard let url = Bundle(for: type(of: self)).url(forResource: "payload", withExtension: "mp4") else { 
            return 
        }
        
        let item = AVPlayerItem(url: url)
        self.player = AVPlayer(playerItem: item)
        self.player?.isMuted = true
        
        self.playerLayer = AVPlayerLayer(player: self.player)
        self.playerLayer?.videoGravity = .resizeAspectFill
        self.playerLayer?.frame = self.bounds
        
        if let playerLayer = self.playerLayer {
            self.layer?.addSublayer(playerLayer)
        }
        
        NotificationCenter.default.addObserver(
            self, 
            selector: #selector(loopVideo(_:)), 
            name: .AVPlayerItemDidPlayToEndTime, 
            object: item
        )
        
        // Start playback after a brief delay
        DispatchQueue.main.asyncAfter(deadline: .now() + 0.1) {
            self.player?.play()
        }
    }
    
    @objc func loopVideo(_ note: Notification) {
        self.player?.seek(to: .zero)
        self.player?.play()
    }
    
    public override func animateOneFrame() {
        super.animateOneFrame()
        // Update layer frame if needed
        if let playerLayer = self.playerLayer {
            playerLayer.frame = self.bounds
        }
    }
    
    public override var hasConfigureSheet: Bool { false }
    public override var configureSheet: NSWindow? { nil }
}
`
}

func infoPlist(name string) string {
    return `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleDevelopmentRegion</key>
    <string>en</string>
    <key>CFBundleExecutable</key>
    <string>VideoSaver</string>
    <key>CFBundleIdentifier</key>
    <string>com.example.` + strings.ToLower(strings.ReplaceAll(name, " ", "")) + `</string>
    <key>CFBundleInfoDictionaryVersion</key>
    <string>6.0</string>
    <key>CFBundleName</key>
    <string>` + name + `</string>
    <key>CFBundlePackageType</key>
    <string>BNDL</string>
    <key>CFBundleShortVersionString</key>
    <string>1.0</string>
    <key>CFBundleVersion</key>
    <string>1</string>
    <key>NSPrincipalClass</key>
    <string>VideoSaverView</string>
</dict>
</plist>
`
}