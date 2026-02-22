# ---- Build stage ----
FROM golang:latest AS builder

WORKDIR /app

# Copy module files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the full source
COPY . .

# Compile the game to WebAssembly
RUN GOOS=js GOARCH=wasm go build -o wasm/goasteroids.wasm .

# Copy the wasm_exec.js that matches this exact Go version.
# The path changed between Go versions so we try both locations.
RUN cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" wasm/wasm_exec.js 2>/dev/null || \
    cp "$(go env GOROOT)/lib/wasm/wasm_exec.js"  wasm/wasm_exec.js

# Build the static file server (serve.go) natively for the current platform.
# CGO_ENABLED=0 ensures a fully static binary compatible with Alpine's musl libc.
RUN CGO_ENABLED=0 go build -o server ./wasm/

# ---- Runtime stage ----
FROM alpine:latest

WORKDIR /app

# Copy only what is needed at runtime
COPY --from=builder /app/wasm/goasteroids.wasm ./wasm/
COPY --from=builder /app/wasm/wasm_exec.js      ./wasm/
COPY --from=builder /app/wasm/index.html        ./wasm/
COPY --from=builder /app/wasm/main.html         ./wasm/
COPY --from=builder /app/server                 ./server

# The file server (serve.go) uses http.Dir("."), so we run it
# from inside the wasm directory so it can find all static files.
WORKDIR /app/wasm

EXPOSE 4000

CMD ["/app/server"]
