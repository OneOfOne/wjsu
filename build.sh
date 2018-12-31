#!/bin/sh
set -e
P="${1:-.}"

if [ ! -d "$P" ]; then
	echo "usage: $0 [path]"
	exit 1
fi

[ -f "${P}/index.html" ] || cp -v "$(dirname $0)/index.html" "${P}/"
cp -v "$(go env GOROOT)/misc/wasm/wasm_exec.js" "${P}/"

cd "${P}"

export GOOS=js
export GOARCH=wasm

goimports -w .
go build -ldflags="-s -w" -o main.wasm
