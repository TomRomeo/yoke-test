compile:
	GOOS=wasip1 GOARCH=wasm go build -o example.wasm ./cmd/music-assistant/main.go

