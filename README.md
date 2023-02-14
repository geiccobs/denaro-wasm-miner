# Denaro WASM miner

## Installation

```bash
git clone https://github.com/geiccobs/denaro-wasm-miner
cd denaro-wasm-miner
```

### Compiling by source

You can skip this if you want to use pre-built wasm binary.  
[Install golang first](https://go.dev/doc/install)
```bash
cd go
go mod tidy
GOOS=js GOARCH=wasm go build -o ../resources/main.wasm
```

## Usage

Start the server using `go run server.go` or connect to an already existing server.  
Notice: memorize the port echoed by the server for the next step.

Open `http://localhost:3010` in your web browser.

Set the parameters according to your own taste and once ready click `Start`.