//+build js,wasm

package react

import (
	"syscall/js"

	"github.com/OneOfOne/wjsu"
)

var (
	funcIdx int32
	ctorIdx int32
	eIdx    int32
)

type (
	StatelessComponent = func(props wjsu.Object) Element
)

func RawReact() js.Value    { return js.Global().Get("React") }
func RawReactDOM() js.Value { return js.Global().Get("ReactDOM") }

func RawCreateElement(args ...interface{}) js.Value {
	return RawReact().Call("createElement", args...)
}

func RawReactComponent() js.Value {
	return RawReact().Get("Component")
}

func HasReact() bool { return wjsu.IsNull(RawReact()) }

// Ref is a React.ref
type Ref struct {
	v js.Value
}

func (r Ref) Current() js.Value { return r.v.Get("current") }
func (r Ref) JSValue() js.Value { return r.v }

func CreateRef() Ref { return Ref{RawReact().Call("createRef")} }

func Fragment() Element { return ReactElement{RawReact().Get("Fragment")} }
