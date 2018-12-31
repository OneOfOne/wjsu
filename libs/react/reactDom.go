//+build js,wasm

package react

import (
	"github.com/OneOfOne/wjsu"
)

func Render(comp interface{}, ele string) {
	RawReactDOM().Call("render", E(comp), wjsu.Document.QuerySelector(ele))
}

// Serve blocks for the lifetime of the app
func Serve(app Component, targetElement string) {
	if wjsu.IsNull(RawReact()) {
		panic("React is not found")
	}
	if wjsu.IsNull(RawReactDOM()) {
		panic("ReactDOM is not found")
	}

	Render(func() Component { return app }, "app")
	select {}
}
