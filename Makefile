.PHONY: test build

test:
	go run . ./test/A5X2/Standard.note
	go run . ./test/A5X2/Artifacts.note
	go run . ./test/A5X2/Realtime.note

build-linux:
	fyne-cross linux

build-windows:
	fyne-cross windows

build-macos:
	fyne package --os darwin
	hdiutil create -volname "Supernote Toolkit" \
		-srcfolder "Supernote Toolkit.app" \
		-ov -format UDZO "dist/darwin-arm64/Supernote Toolkit.dmg"
	rm "Supernote Toolkit.app"

build-android:
	fyne-cross android

build:
	$(MAKE) build-linux
	$(MAKE) build-windows
	$(MAKE) build-macos
	$(MAKE) build-android
