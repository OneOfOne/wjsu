#!/bin/sh
set -e
P="${1:-.}"

if [ ! -d "$P" ]; then
	echo "usage: $0 [path]"
	exit 1
fi

[ -f "${P}/index.html" ] || cp -v index.html "${P}/"
cp -v "$(go env GOROOT)/misc/wasm/wasm_exec.js" "${P}/"

cd "${P}"

env GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o main.wasm
