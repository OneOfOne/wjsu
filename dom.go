//+build js,wasm

package jsenv

import (
	"syscall/js"
)

var (
	document = js.Global().Get("document")
	head     = document.Get("head")

	Document doc
)

func RawDocument() js.Value { return document }

type Event interface {
	PreventDefault()

	StopPropagation()
	StopImmediatePropagation()
}

type event struct {
	js.Value
}

func (e event) PreventDefault()  { e.Call("preventDefault") }
func (e event) StopPropagation() { e.Call("stopPropagation") }

func (e event) StopImmediatePropagation() { e.Call("stopImmediatePropagation") }

type HTMLElement struct {
	v js.Value
}

func toString(v js.Value) string {
	if v.Type() == js.TypeString {
		return v.String()
	}

	return ""
}

func (e HTMLElement) GetAttribute(key string) string {
	return toString(e.v.Call("getAttribute", key))
}

func (e HTMLElement) SetAttribute(key, value string) {
	e.v.Call("setAttribute", key, value)
}

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

func (e HTMLElement) AddEventListener(evt string, cb func(evt Event), opts *ListenerOptions) {
	var (
		jcb js.Func
	)

	jcb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if opts != nil && opts.Once {
			jcb.Release()
		}
		cb(event{args[0]})
		return nil
	})

	e.v.Call("addEventListener", evt, jcb, opts)
}

func (e HTMLElement) Object() Object { return Object{v: e.v} }

func (e HTMLElement) JSValue() js.Value { return e.v }

type doc struct{}

func (doc) CreateElement(typ string, props map[string]string) HTMLElement {
	ele := HTMLElement{v: document.Call("createElement", typ)}
	for k, v := range props {
		ele.SetAttribute(k, v)
	}

	return ele
}

func (doc) Title() string {
	return toString(document.Get("title"))
}

func (doc) SetTitle(title string) {
	document.Set("title", title)
}

func AddScript(src string, deferProp bool) <-chan struct{} {
	var (
		ch = make(chan struct{})
		cb js.Func
	)

	cb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		cb.Release()
		close(ch)
		return nil
	})

	ele := document.Call("createElement", "script")
	ele.Set("src", src)
	ele.Set("onload", cb)

	if deferProp {
		ele.Call("setAttribute", "defer", "")
	}

	head.Call("appendChild", ele)

	return ch
}

func LoadScripts(srcs ...string) <-chan struct{} {
	chs := make([]<-chan struct{}, 0, len(srcs))
	ch := make(chan struct{})
	for _, src := range srcs {
		chs = append(chs, AddScript(src, true))
	}

	go func() {
		for _, ch := range chs {
			<-ch
		}
		close(ch)
	}()

	return ch
}
