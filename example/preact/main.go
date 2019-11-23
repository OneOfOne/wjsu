//+build js,wasm

package main

import (
	"syscall/js"

	// this helps with stutering, for example wjsu.Console.Log, 100% optional.
	. "github.com/OneOfOne/wjsu"
)

func main() {
	Initialize()
	preact := js.Global().Get("preact")
	Console.Log(preact)
	comp := preact.Get("Component").New()
	comp.Set("render", xrender)
	h1 := H("h1", O("style", O("color", "red")), comp)
	Console.Log(h1, Document.Body(), comp)
	Render(comp, Document.Body())
}

var xrender = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	Console.Log("blah", this, args)
	return H("h1", O("style", O("color", "red")), "blah")
})
