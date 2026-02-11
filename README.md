<div style="width:100%" align="center"><img alt="Supernote Toolkit" src="Icon.png" /></div>
<h1 align="center">Supernote Toolkit</h1>
<p align="center">A collection of tools for tinkering with supernote files.</p>
<br />

## Showcase (v0.1.0)

https://github.com/user-attachments/assets/944129de-6cde-452c-9e60-741618ff3739

## Features

For now, it only converts a folder of .note files into corresponding png/pdf files.

See [roadmap](#Roadmap) for more. See [changelog](/CHANGELOG.md) for all changes so far.

## Latest Changelog

- Add support for Nomad X2 devices.

- Add support for Horizontal orientations.

- Add cli flag (-device) and UI radio buttons (Select Device) to manually impose a device for all notes for the conversion.

## Install

Download the latest binaries from [release page](https://github.com/p-tupe/supernote-toolkit/releases/latest). Verify your download against checksum.txt included in each release. See OS-specific instructions below.

> You can also just `go install github.com/p-tupe/supernote-toolkit` if you have `go` installed.

### macOS

The app is unsigned. On first launch, macOS will block it.

```bash
xattr -d com.apple.quarantine "Supernote Toolkit.app"
```

Or: System Settings > Privacy & Security > Open Anyway.

### Linux

```bash
tar xf "Supernote Toolkit.tar.xz"
make user-install
```

### Windows

Extract the zip and run `Supernote Toolkit.exe`. Windows may show a SmartScreen warning since the app is unsigned -- click "More info" then "Run anyway".

## Build from source

Requires [Go](https://go.dev/dl/) 1.24+ and [Fyne](https://docs.fyne.io/started/) dependencies.

```bash
git clone https://github.com/p-tupe/supernote-toolkit.git
cd supernote-toolkit
go run .
```

## Use the CLI

```bash
git clone https://github.com/p-tupe/supernote-toolkit.git
cd supernote-toolkit
go install cmd/cli/supernote-toolkit.go
supernote-toolkit
```

## Roadmap

### In-progress

- [ ] Extract Text (from Realtime notes)
- [ ] Horizontal orientation
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
