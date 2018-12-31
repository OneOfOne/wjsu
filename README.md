# wjsu [![GoDoc](https://godoc.org/github.com/OneOfOne/wjsu?status.svg)](https://godoc.org/github.com/OneOfOne/wjsu)

wjsu (aka Webassembly Javascript Utils) is a helper package for [`syscall/js`](https://godoc.org/syscall/js).

Inpsired by [honnef.co/go/js/dom](https://github.com/dominikh/go-js-dom).

## Install

    GOOS=js GOARCH=wasm go get -u github.com/OneOfOne/wjsu

* **note**: the package requires Go 1.12 for the `syscall/js` changes.

## Features

* Object.
* Array (partial).
* console (Log, Warn and Error).
* wrapper for loading external scripts.
* wrapper for js.ValueOf that doesn't panic and handles slices / maps.

## Usage

x

### main.go

```go
//+build js,wasm

package main

import (
	"time"

	// this helps with stutering, for example wjsu.Console.Log, 100% optional.
	. "github.com/OneOfOne/wjsu"
)

func main() {
	Console.Log("hello console", []int{1, -1}, map[string]float64{"a": 1, "z": -1})
	Document.QuerySelector("div#loader").SetInnerHTML("hello world")

	app := Document.QuerySelector("app")
	for t := range time.Tick(time.Second) {
		app.SetInnerHTML("time (UTC): " + t.UTC().String())
	}
}
```

### index.html

```html
<html>
	<head>
		<title>WASM Loader</title>
		<meta charset="utf-8">
		<script>
			const goInit = () => {
				const go = new Go();
				WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
					go.run(result.instance);
				});
			}
		</script>
		<script defer src="wasm_exec.js" onload="goInit()"></script>
	</head>
	<body>
		<div id="loader">Loading...</div>
		<app />
	</body>
</html>
```

### serve.go

```go
//+build !js

package main

import (
	"net/http"
	"os/exec"
	"time"
	"log"
)

func main() {
	time.AfterFunc(time.Second/2, func() {
		exec.Command("xdg-open", "http://localhost:8080").Start()
	})
	log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir("."))))
}
```

### build & run

```sh
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
GOOS=js GOARCH=wasm go build -o main.wasm
go run serve.go
```

## VSCode settings

```json
{
	"go.toolsEnvVars": {
		"GOOS": "js",
		"GOARCH": "wasm"
	}
}
```

## TODO

* Proper documentation.

## License

This project is released under the Apache v2. licence. See [LICENCE](LICENCE) for more details.
