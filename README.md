# Go Asteroids

## Wasm build command:

- $env:GOOS="js"; $env:GOARCH="wasm"; go build -o wasm/goasteroids.wasm .

## Native build command:

- go build -o go-asteroids.exe .

## Winres command:

- go-winres patch go-asteroids.exe

## Docker build command:

- docker build -t go-asteroids .
- docker run --rm -p 4000:4000 go-asteroids
