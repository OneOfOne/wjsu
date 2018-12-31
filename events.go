//+build js,wasm

package jsenv

import (
	"syscall/js"
)

type ListenerOptions struct {
	Capture bool
	Once    bool
	Passive bool
}

func (lo *ListenerOptions) JSValue() js.Value {
	if lo == nil {
		return js.Null()
	}
	o := object.New()
	if lo.Capture {
		o.Set("capture", true)
	}
	if lo.Once {
		o.Set("once", true)
	}
	if lo.Passive {
		o.Set("passive", true)
	}
	return o
}
