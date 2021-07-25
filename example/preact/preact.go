package main

import (
	"reflect"
	"syscall/js"

	. "github.com/OneOfOne/wjsu"
)

var Preact js.Value

func init() {
	Preact = js.Global().Get("preact")
}

type PreactElement interface {
	Element() js.Value
}

type preactNode struct {
	js.Value
}

func (pn preactNode) Element() js.Value { return pn.Value }

func H(tagOrComponent, props interface{}, children ...interface{}) PreactElement {
	n := preactNode{Preact.Call("h", tagOrComponent, props, ArgsToObjects(children))}
	return n
}

func Render(node PreactElement, parent interface{}) {
	switch p := parent.(type) {
	case HTMLElement:
		// do nothing
	case string:
		parent = Document.Body().QuerySelector(p)
	case nil:
		parent = Document.Body()
	default:
		panic("unknown parent type: " + reflect.TypeOf(parent).String())
	}
	Preact.Call("render", node.Element(), parent)
}
