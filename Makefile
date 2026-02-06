.PHONY: test

test:
	go run . ./test/A5X2/Standard.note
	go run . ./test/A5X2/Artifacts.note
	go run . ./test/A5X2/Realtime.note
