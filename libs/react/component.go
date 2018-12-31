//+build js,wasm

package react

import (
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"syscall/js"

	"github.com/OneOfOne/wjsu"
)

type ComponentDef struct {
	Props wjsu.Object
	State wjsu.Object

	p  js.Value
	hm map[string]js.Func
}

func (c *ComponentDef) SetState(state wjsu.Object, cb func()) {
	c.State = state.ReadOnly()

	if cb == nil {
		c.p.Call("setState", state)
		return
	}

	c.p.Call("setState", state, wjsu.WrapFunc(cb, true))
}

func (c *ComponentDef) ForceUpdate(cb func()) {
	if cb == nil {
		c.p.Call("forceUpdate")
		return
	}

	c.p.Call("forceUpdate", wjsu.WrapFunc(cb, true))
}

func (c *ComponentDef) SetHandler(name string, h func()) js.Func {
	jfn := wjsu.WrapFunc(h, false)

	if c.hm == nil {
		c.hm = map[string]js.Func{}
	}

	if ofn, ok := c.hm[name]; ok {
		// if we're re-assigning it, release the old version
		ofn.Release()
	}

	c.hm[name] = jfn

	return jfn
}

func (c *ComponentDef) Handler(name string) js.Func {
	return c.hm[name]
}

func (ComponentDef) Ctor() wjsu.Object { return wjsu.Null() }

func (ComponentDef) ComponentWillMount()   {}
func (ComponentDef) ComponentDidMount()    {}
func (ComponentDef) ComponentWillUpdate()  {}
func (ComponentDef) ComponentDidUpdate()   {}
func (ComponentDef) ComponentWillUnmount() {}

func (c *ComponentDef) def() *ComponentDef { return c }

func (c *ComponentDef) cleanup() {
	for k, jfn := range c.hm {
		jfn.Release()
		delete(c.hm, k)
	}
}

func (c *ComponentDef) JSValue() js.Value { return c.p }

func (c *ComponentDef) init(p js.Value, props wjsu.Object) {
	c.p = p
	c.Props = props.ReadOnly()
}

type DebugComponent struct{ ComponentDef }

func (c *DebugComponent) Ctor() wjsu.Object {
	wjsu.Console.Log("ctor", c.p)
	return wjsu.Null()
}

func (c *DebugComponent) ComponentWillMount()   { wjsu.Console.Log("will mount", c.p) }
func (c *DebugComponent) ComponentDidMount()    { wjsu.Console.Log("did mount", c.p) }
func (c *DebugComponent) ComponentWillUpdate()  { wjsu.Console.Log("will update", c.p) }
func (c *DebugComponent) ComponentDidUpdate()   { wjsu.Console.Log("did update", c.p) }
func (c *DebugComponent) ComponentWillUnmount() { wjsu.Console.Log("will unmount", c.p) }

type Component interface {
	Ctor() (state wjsu.Object)
	ComponentWillMount()
	ComponentDidMount()
	ComponentWillUpdate()
	ComponentDidUpdate()
	ComponentWillUnmount()

	Render() Element

	def() *ComponentDef
}

type ComponentWithName interface {
	Component

	Name() string
}

func createCtor(ctor func() Component, invoke bool) js.Value {
	var (
		fn   js.Func
		init bool

		reactComp = RawReactComponent()
	)

	fn = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		props := wjsu.ObjectOf(args[0], true)
		obj := reactComp.New(props)
		comp := ctor()
		comp.def().init(obj, props)

		if !init {
			if c, ok := comp.(ComponentWithName); ok {
				fn.Set("displayName", c.Name())
			} else {
				fn.Set("displayName", reflect.TypeOf(comp).Elem().Name())
			}

			init = true
		}

		if state := comp.Ctor(); state.Valid() {
			obj.Set("state", state)
			comp.def().State = state.ReadOnly()
		}

		obj.Set("render", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			return comp.Render()
		}))

		obj.Set("componentWillMount", wjsu.WrapFunc(comp.ComponentWillMount, false))
		obj.Set("componentDidMount", wjsu.WrapFunc(comp.ComponentDidMount, false))
		obj.Set("componentWillUpdate", wjsu.WrapFunc(comp.ComponentWillUpdate, false))
		obj.Set("componentDidUpdate", wjsu.WrapFunc(comp.ComponentDidUpdate, false))
		obj.Set("componentWillUnmount", wjsu.WrapFunc(func() {
			comp.def().cleanup()
			comp.ComponentWillUnmount()
		}, false))

		return obj
	})

	if invoke {
		o := wjsu.O("key", "ctor:"+strconv.Itoa(int(atomic.AddInt32(&ctorIdx, 1))))
		return RawCreateElement(fn, o)
	}

	return fn.JSValue()
}

func wrapFunc(ctor StatelessComponent, invoke bool) js.Value {
	var (
		fn   js.Func
		init bool
	)

	fn = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		props := wjsu.ObjectOf(args[0], true)

		if !init {
			name := runtime.FuncForPC(reflect.ValueOf(ctor).Pointer()).Name()
			if idx := strings.LastIndexByte(name, '.'); idx > -1 {
				name = name[idx+1:]
			}
			fn.Set("displayName", name)
			init = true
		}

		return ctor(props)
	})

	if invoke {
		o := wjsu.O("key", "func:"+strconv.Itoa(int(atomic.AddInt32(&funcIdx, 1))))
		return RawCreateElement(fn, o)
	}

	return fn.JSValue()
}
