//+build js,wasm

package wjsu

import (
	"syscall/js"
)

var (
	document = js.Global().Get("document")
	head     = document.Get("head")

	Document = HTMLDocument{HTMLElement{document}}
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

func (e HTMLElement) Class() string {
	return toString(e.v.Get("className"))
}

func (e HTMLElement) SetClass(s string) {
	e.v.Set("className", s)
}

func (e HTMLElement) ID() string {
	return toString(e.v.Get("id"))
}

func (e HTMLElement) SetID(s string) {
	e.v.Set("id", s)
}

func (e HTMLElement) HasAttribute(s string) bool {
	return e.v.Call("hasAttribute", s).Bool()
}

func (e HTMLElement) HasAttributeNS(ns string, name string) bool {
	return e.v.Call("hasAttributeNS", ns, name).Bool()
}

func (e HTMLElement) GetAttribute(name string) string {
	return toString(e.v.Call("getAttribute", name))
}

func (e HTMLElement) GetAttributeNS(ns, name string) string {
	return toString(e.v.Call("getAttributeNS", ns, name))
}

func (e HTMLElement) SetAttribute(name, value string) {
	e.v.Call("setAttribute", name, value)
}

func (e HTMLElement) SetAttributeNS(ns, name, value string) {
	e.v.Call("setAttributeNS", ns, name, value)
}

func (e HTMLElement) AppendChild(child HTMLElement) {
	e.v.Call("appendChild", child.v)
}

func (e HTMLElement) SetInnerHTML(html string) {
	e.v.Set("innerHTML", html)
}

// unless opts.Once is true, this call will leak a js.Func when the element gets destroyed
func (e HTMLElement) AddEventListener(evt string, cb func(evt Event), opts *ListenerOptions) {
	var jcb js.Func

	jcb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if opts != nil && opts.Once {
			jcb.Release()
		}
		cb(event{args[0]})
		return nil
	})

	e.v.Call("addEventListener", evt, jcb, opts)
}

func (e HTMLElement) On(evt string, cb func(evt Event), once bool) {
	var jcb js.Func

	jcb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if once {
			jcb.Release()
		}
		cb(event{args[0]})
		return nil
	})

	e.v.Set("on"+evt, jcb)
}

func (e HTMLElement) QuerySelector(q string) HTMLElement {
	return HTMLElement{e.v.Call("querySelector", q)}
}

// TODO NodeList
func (e HTMLElement) QuerySelectorAll(q string) []HTMLElement {
	// do NodeList
	a := e.v.Call("querySelectorAll", q)
	if a.Type() != js.TypeObject {
		return nil
	}
	ln := a.Length()
	out := make([]HTMLElement, ln)
	for i := 0; i < ln; i++ {
		out[i].v = a.Index(i)
	}
	return out

}

func (e HTMLElement) RawGet(k string) js.Value {
	return e.v.Get(k)
}

func (e HTMLElement) RawSet(k string, v interface{}) {
	e.v.Set(k, v)
}

func (e HTMLElement) RawCall(method string, args ...interface{}) js.Value {
	return e.v.Call(method, args...)
}

func (e HTMLElement) Object() Object { return Object{v: e.v} }

func (e HTMLElement) JSValue() js.Value { return e.v }

type HTMLDocument struct {
	HTMLElement
}

func (HTMLDocument) CreateElement(typ string, props map[string]string) HTMLElement {
	ele := HTMLElement{v: document.Call("createElement", typ)}
	for k, v := range props {
		ele.SetAttribute(k, v)
	}

	return ele
}

func (HTMLDocument) Title() string {
	return toString(document.Get("title"))
}

func (HTMLDocument) SetTitle(title string) {
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
