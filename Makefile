VERSION ?= $(shell git describe --tags --always --dirty)
LDFLAGS = -ldflags "-X=github.com/p-tupe/supernote-toolkit/cmd/app.Version=$(VERSION)"

.PHONY: test build

test:
	go run . ./test/A5X2/Standard.note
	go run . ./test/A5X2/Artifacts.note
	go run . ./test/A5X2/Realtime.note

build:
	$(MAKE) clean-dist
	$(MAKE) build-linux
	$(MAKE) build-windows
	$(MAKE) build-macos
	$(MAKE) build-android
	$(MAKE) move-dist
	$(MAKE) flatten-dist
	$(MAKE) checksum
	$(MAKE) clean-tmp

build-linux:
	fyne-cross linux -pull -arch=arm,arm64,amd64 $(LDFLAGS)

build-windows:
	fyne-cross windows -pull -arch=amd64,arm64 $(LDFLAGS)

build-macos:
	go build $(LDFLAGS) -o supernote-toolkit ./cmd/app
	fyne package --os darwin --exe supernote-toolkit
	rm -f supernote-toolkit
	mkdir -p fyne-cross/dist/darwin-arm64
	hdiutil create -volname "Supernote Toolkit" \
		-srcfolder "Supernote Toolkit.app" \
		-ov -format UDZO "fyne-cross/dist/darwin-arm64/Supernote Toolkit.dmg"

build-android:
	fyne-cross android -pull $(LDFLAGS)

checksum:
	rm -f dist/checksum.txt
	fd . -t f dist | while read f; do sha256 -b $$f >> dist/checksum.txt; done

flatten-dist:
	mv "dist/darwin-arm64/Supernote Toolkit.dmg" dist/darwin-arm64.dmg
	mv "dist/linux-amd64/Supernote Toolkit.tar.xz" dist/linux-amd64.tar.xz
	mv "dist/linux-arm64/Supernote Toolkit.tar.xz" dist/linux-arm64.tar.xz
	mv "dist/linux-arm/Supernote Toolkit.tar.xz" dist/linux-arm.tar.xz
	mv "dist/windows-amd64/Supernote Toolkit.exe.zip" dist/windows-amd64.exe.zip
	mv "dist/windows-arm64/Supernote Toolkit.exe.zip" dist/windows-arm64.exe.zip
	mv "dist/android/Supernote Toolkit.apk" dist/android.apk
	rm -rf dist/darwin-arm64 dist/linux-amd64 dist/linux-arm64 dist/linux-arm dist/windows-amd64 dist/windows-arm64 dist/android

move-dist:
	mv fyne-cross/dist ./

clean-dist:
	rm -rf dist

clean-tmp:
	rm -rf "Supernote Toolkit.app"
	rm -rf fyne-cross
