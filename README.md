# Denaro WASM miner

## Installation

```bash
git clone https://github.com/geiccobs/denaro-wasm-miner
cd denaro-wasm-miner
```

### Compiling by source
[Install golang first](https://go.dev/doc/install)
#### WASM
You can skip this if you want to use pre-built wasm binary (resources/main.wasm).
```bash
cd go
GOOS=js GOARCH=wasm go build -o ../resources/main.wasm
```
#### Server
You can skip this if you want to use pre-built server binary ([download latest binary](https://github.com/geiccobs/denaro-pool-miner/releases/latest)).
```bash
go build server.go
```

## Usage

Start the server using `./server` or connect to an already existing server.  
Notice: memorize the port echoed by the server for the next step.

Open `http://localhost:3010` in your web browser (where "3010" is the port saw before),

Set the parameters according to your own taste and once ready click `Start`.