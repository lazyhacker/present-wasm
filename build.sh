#!/bin/sh


if [ -x present.wasm ]; then
    rm present.wasm
fi

GOROOT=$HOME/go-wasm
GOOS=js GOARCH=wasm $HOME/go-wasm/bin/go build -o present.wasm present.go dir.go play.go play_http.go
