package main

import (
	"syscall/js"

	. "github.com/OneOfOne/wjsu"
)

var (
	preact js.Value
)

func init() {
	preact = js.Global().Get("preact")
}

func H(name string, attrs Object, children ...interface{}) js.Wrapper {
	if len(children) > 0 {
		attrs.Set("children", children)
	}
	return preact.Call("h", name, attrs)
}

func Render(node, parent js.Wrapper) {
	preact.Call("render", node, parent)
}
