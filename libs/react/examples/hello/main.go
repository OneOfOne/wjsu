//+build js,wasm

package main

import (
	"time"

	. "github.com/OneOfOne/wjsu"
	R "github.com/OneOfOne/wjsu/libs/react"
)

func H1(props Object) R.Element {
	if props.Has("children") {
		return R.E("h1", props, "with props: ")
	}
	return R.E("h1", "no props :(")
}

type Hello struct {
	R.DebugComponent
}

func (h *Hello) Ctor() (state Object) {
	h.SetHandler("onClick", h.toggleColor)
	return O("color", "blue", "text", "Hello")
}

func (h *Hello) toggleColor() {
	Console.Log("toggleColor state:", h.State)

	o := h.State.Copy()
	if o.Get("color").String() == "red" {
		o.Set("color", "blue")
		o.Set("text", "Hello")
	} else {
		o.Set("color", "red")
		o.Set("text", "Bye")
	}

	h.SetState(o, nil)
}

func (h *Hello) Render() R.Element {
	s := h.State

	return R.E(nil, Null(),
		H1,
		R.E(H1, "w00t"),
		R.E("h1", O("style", O("color", s.Get("color"))),
			s.GetString("text")+": "+time.Now().UTC().String(),
		),
		R.E("button", O("onClick", h.Handler("onClick")), "toggle"),
	)
}

type App struct {
	R.ComponentDef
}

func (App) Name() string { return "CoolApp" }
func (App) Render() R.Element {
	return R.E(nil, Null(),
		R.E("h1", Null(), "Touchy the button"),
		R.E(&Hello{}, Null()),
	)
}

func main() {
	Document.SetTitle("Hello World")
	Document.QuerySelector("body").RemoveChild(Document.QuerySelector("#loader"))
	R.Serve(&App{}, "app")
}
