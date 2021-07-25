#!/bin/sh
set -e
set -x

P="${1:-.}"

if [ ! -d "$P" ]; then
	echo "usage: $0 [path]"
	exit 1
fi

export PATH=$(go env GOBIN):/bin

[ -f "${P}/index.html" ] || cp -v "$(dirname $0)/index.html" "${P}/"
cp -v "$(tinygo env TINYGOROOT)/targets/wasm_exec.js" "${P}/"
cp -v "$(dirname $0)/wasm.js" "${P}/"

cd "${P}"

# goimports -w .
tinygo build -target wasm -o main.wasm .
