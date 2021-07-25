//go:build js && wasm
// +build js,wasm

package main

import (
	// this helps with stutering, for example wjsu.Console.Log, 100% optional.
	. "github.com/OneOfOne/wjsu"
)

func main() {
	changeColor := func(co ComponentObject) {
		Console.Log(co.JSValue())
		state := co.State()
		color := state.Get("color").String()
		if color == "red" {
			color = "green"
		} else {
			color = "red"
		}
		state.Set("color", color)
		co.SetState(state)
	}

	c := &Component{
		Name: "uwutm8",
		Init: func(o ComponentObject) {
			Console.Log("constructor", o.Props(), o.Context())
			o.SetState(O("color", "green"))
		},
		DidMount:    func(o ComponentObject) { Console.Log("didMount") },
		WillUnmount: func(o ComponentObject) { Console.Log("willUmount") },
		ShouldUpdate: func(o ComponentObject, nextProps, nextState Object) bool {
			Console.Log("shouldUpdate", nextProps, nextState)
			return true
		},
		DidCatch: func(o ComponentObject, err Object) {
			Console.Log("didCatch", err)
		},
		Render: func(o ComponentObject) interface{} {
			Console.Log("render", o.State())
			color := o.State().GetString("color")
			return H("h1", O("style", O("color", color), "onclick", func() { changeColor(o) }), "hi")
		},
	}

	Render(c, Document.Body())
	select {}
}
