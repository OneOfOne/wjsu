#!/bin/sh
set -e
set -x

P="${1:-.}"

if [ ! -d "$P" ]; then
	echo "usage: $0 [path]"
	exit 1
fi

[ -f "${P}/index.html" ] || cp -v "$(dirname $0)/index.html" "${P}/"
cp -v "$(go env GOPATH)/src/github.com/tinygo-org/tinygo/targets/wasm_exec.js" "${P}/"
cp -v "$(dirname $0)/wasm.js" "${P}/"

cd "${P}"

# goimports -w .
TGO=$(go env GOBIN)/tinygo
env PATH=/bin $TGO build -target wasm -tags go1.13 -o main.wasm .
