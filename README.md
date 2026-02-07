<div style="width:100%" align="center"><img alt="Supernote Toolkit" src="Icon.png" /></div>
<h1 align="center">Supernote Toolkit</h1>
<p align="center">A collection of tools for tinkering with supernote files.</p>
<br />

## Showcase (v0.1.0)

https://github.com/user-attachments/assets/944129de-6cde-452c-9e60-741618ff3739


## Features

For now, it only converts a folder of .note files into corresponding png/pdf files. See [roadmap](#Roadmap) for more.

## Install

Download the latest release from below:

| Platform              | Download                                                                                                                     |
| --------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| macOS (Apple Silicon) | <a href="https://github.com/p-tupe/supernote-toolkit/raw/main/dist/darwin-arm64.dmg" download>darwin-arm64.dmg</a>           |
| Linux (amd64)         | <a href="https://github.com/p-tupe/supernote-toolkit/raw/main/dist/linux-amd64.tar.xz" download>linux-amd64.tar.xz</a>       |
| Linux (arm64)         | <a href="https://github.com/p-tupe/supernote-toolkit/raw/main/dist/linux-arm64.tar.xz" download>linux-arm64.tar.xz</a>       |
| Linux (arm)           | <a href="https://github.com/p-tupe/supernote-toolkit/raw/main/dist/linux-arm.tar.xz" download>linux-arm.tar.xz</a>           |
| Windows (amd64)       | <a href="https://github.com/p-tupe/supernote-toolkit/raw/main/dist/windows-amd64.exe.zip" download>windows-amd64.exe.zip</a> |
| Windows (arm64)       | <a href="https://github.com/p-tupe/supernote-toolkit/raw/main/dist/windows-arm64.exe.zip" download>windows-arm64.exe.zip</a> |

Verify your download against <a href="https://github.com/p-tupe/supernote-toolkit/raw/main/dist/checksum.txt" download>checksum.txt</a> (SHA-256) included in each release.

### macOS

The app is unsigned. On first launch, macOS will block it.

```bash
xattr -d com.apple.quarantine "Supernote Toolkit.app"
```

Or: System Settings > Privacy & Security > Open Anyway.

### Linux

```bash
tar xf "Supernote Toolkit.tar.xz"
make user-install # or ./supernote-toolkit
```

### Windows

Extract the zip and run `Supernote Toolkit.exe`. Windows may show a SmartScreen warning since the app is unsigned -- click "More info" then "Run anyway".

### Android

Enable "Install from unknown sources" in your device settings, then install the APK.

## Build from source

Requires [Go](https://go.dev/dl/) 1.24+ and [Fyne](https://docs.fyne.io/started/) dependencies.

```bash
git clone https://github.com/p-tupe/supernote-toolkit.git
cd supernote-toolkit
go run .
```

## Roadmap

### In-progress

- [ ] Extract Text (from Realtime notes)
- [ ] Option to recurse through sub-folders
- [ ] Improve UI somehow
  - [ ] Make errors more visible
  - [ ] Show a log?

### Don't hold your breath

- [ ] Input from file server (from device)
- [ ] Input from private dav server
- [ ] Input from supernote cloud
- [ ] Automatically convert notes on change in a pre-configured folder
- [ ] Convert to SVG
- [ ] Convert Text using OCR
