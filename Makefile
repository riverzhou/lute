all: lute

lute:
	GOOS=js GOARCH=wasm go build -ldflags="-s -w" .

clean:
	rm -rf lute.wasm

