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

type Node interface {
	Node() js.Wrapper
}

type preactNode js.Value

func (pn preactNode) Node() js.Wrapper { return js.Value(pn) }

func H(tagOrComponent, props interface{}, children ...interface{}) Node {
	n := preactNode(Preact.Call("h", tagOrComponent, props, ArgsToObjects(children)))
	return n
}

func Render(node Node, parent interface{}) {
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
	Preact.Call("render", node.Node().JSValue(), parent)
}
