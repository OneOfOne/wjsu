# jsenv [![GoDoc](https://godoc.org/github.com/OneOfOne/jsenv?status.svg)](https://godoc.org/github.com/OneOfOne/jsenv)

This is a helper package for [`syscall/js`](https://godoc.org/syscall/js).

Inpsired by [honnef.co/go/js/dom](https://github.com/dominikh/go-js-dom).

## Install

    go get github.com/OneOfOne/jsenv

## Features

* Object.
* Array (partial).
* console (Log, Warn and Error).
* wrapper for loading external scripts.
* wrapper for js.ValueOf that doesn't panic and handles slices / maps.

## Usage

```go
	import (
		. "github.com/OneOfOne/jsenv"
	)

	func main() {
		Console.Log("hi", []int{1, -1}, map[string]float64 {"a": 1, "z":-1})

		a := NewArray(0)
		a.Push("a", 1, 2.2)
		na := a.Concat([]int{1, 2, 3})
	}
```

## TODO

* Proper documentation.

## License

This project is released under the Apache v2. licence. See [LICENCE](LICENCE) for more details.
