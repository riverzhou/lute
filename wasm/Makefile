
all: lute 

lute:
	GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o lute.wasm .

clean:
	rm -rf lute.wasm

