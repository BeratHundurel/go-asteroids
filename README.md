# Go Asteroids

## Wasm build command:
  - $env:GOOS="js"; $env:GOARCH="wasm"; go build -o wasm/goasteroids.wasm .
  
## Winres command:
  - go-winres patch go-asteroids.exe