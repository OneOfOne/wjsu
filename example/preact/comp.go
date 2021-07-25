package main

import (
	"strings"
	"syscall/js"

	. "github.com/OneOfOne/wjsu"
)

const comp = `
(({Component}) => {
class %s extends Component {
	constructor(props, ctx) {
		super(props, ctx)
		if (this._ctor) this._ctor(this);
	}
	render() {
		return "nil render func";
	}
}
return %s;
})(preact);
`

func createComp(name string) js.Value {
	return Window().Call("eval", strings.ReplaceAll(comp, "%s", name))
}

var _ = Console

type Component struct {
	Name string

	Init         func(o ComponentObject)
	DidMount     func(o ComponentObject)
	WillUnmount  func(o ComponentObject)
	ShouldUpdate func(o ComponentObject, nextProps, nextState Object) bool
	DidCatch     func(o ComponentObject, err Object)
	Render       func(o ComponentObject) interface{}

	comp js.Value
}

func (c *Component) Clone(name string) *Component {
	cp := *c
	cp.Name = name
	cp.comp = js.Undefined()
	return &cp
}

func (c *Component) JSValue() js.Value {
	if !IsNull(c.comp) {
		return c.comp
	}

	if c.Name == "" {
		panic("component must have a name")
	}

	if c.Render == nil {
		panic("render can't be nil")
	}

	comp := createComp(c.Name)
	proto := comp.Get("prototype")

	if c.Init != nil {
		proto.Set("_ctor", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			c.Init(ComponentObject{ObjectOf(this), c.Name})
			return nil
		}))
	}

	if c.DidMount != nil {
		proto.Set("componentDidMount", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			c.DidMount(ComponentObject{ObjectOf(this), c.Name})
			return nil
		}))
	}

	if c.WillUnmount != nil {
		proto.Set("componentWillUnmount", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			c.WillUnmount(ComponentObject{ObjectOf(this), c.Name})
			return nil
		}))
	}

	if c.ShouldUpdate != nil {
		proto.Set("shouldComponentUpdate", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			return c.ShouldUpdate(ComponentObject{ObjectOf(this), c.Name}, ObjectOf(args[0], true), ObjectOf(args[1], true))
		}))
	}

	if c.DidCatch != nil {
		proto.Set("componentDidCatch", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			c.DidCatch(ComponentObject{ObjectOf(this), c.Name}, ObjectOf(args[0], true))
			return nil
		}))
	}

	proto.Set("render", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if c.Render != nil {
			return ValueOf(c.Render(ComponentObject{ObjectOf(this), c.Name})).JSValue()
		}
		return "no render func"
	}))
	c.comp = comp
	return c.comp
}

type ComponentObject struct {
	Object
	name string
}

func (c ComponentObject) Name() string { return c.name }
func (c ComponentObject) Context() Object {
	if c.Valid() {
		return c.Get("context").ReadOnly()
	}
	return Undefined()
}

func (c ComponentObject) Props() Object {
	if c.Valid() {
		return c.Get("props").ReadOnly()
	}
	return Undefined()
}

func (c ComponentObject) State() Object {
	if c.Valid() {
		return c.Get("state").Copy()
	}
	return Undefined()
}

func (c ComponentObject) SetState(o Object) {
	if c.Valid() {
		c.CallByName("setState", o)
	}
}
