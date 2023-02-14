DIR="builds/"

# remove builds dir if exists
if [ -d "$DIR" ]; then
    rm -r $DIR
fi

# compile for linux
# 64 bit
GOOS=linux GOARCH=amd64 go build -o $DIR/server-linux64 server.go
# 32 bit
GOOS=linux GOARCH=386 go build -o $DIR/server-linux32 server.go

echo "compiled for linux"

# compile for windows
# 64 bit
GOOS=windows GOARCH=amd64 go build -o $DIR/server-win64.exe server.go
# 32 bit
GOOS=windows GOARCH=386 go build -o $DIR/server-win32.exe server.go

echo "compiled for windows"

# compile for macos
# amd 64 bit
GOOS=darwin GOARCH=amd64 go build -o $DIR/server-macos-amd64 server.go
# arm 64 bit
GOOS=darwin GOARCH=arm64 go build -o $DIR/server-macos-arm64 server.go

echo "compiled for macos"

# compile for android
# amd 64 bit
CGO_ENABLED=0 GOOS=android GOARCH=amd64 go build -o $DIR/server-android-amd64 server.go
# arm 64 bit
CGO_ENABLED=0 GOOS=android GOARCH=arm64 go build -o $DIR/server-android-arm64 server.go

echo "compiled for android"
