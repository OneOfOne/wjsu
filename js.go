//+build js,wasm

package jsenv

import (
	"syscall/js"
)

var (
	object = js.Global().Get("Object")

	nullObj  = Object{v: js.Null(), ro: true}
	undefObj = Object{v: js.Undefined(), ro: true}

	_ js.Wrapper = (*Object)(nil)
)

func Null() Object      { return nullObj }
func Undefined() Object { return undefObj }

func RawObject() js.Value { return object }

// IsNull checks if an o is Null or Undefined
func IsNull(o js.Wrapper) bool {
	return o == js.Null() || o == js.Undefined()
}

// O is a shortcut for NewObject(false).SetMulti(keyVals...)
func O(keyVals ...interface{}) Object {
	return NewObject(false).SetMulti(keyVals...)
}

func NewObject(ro bool) Object {
	return Object{v: object.New(), ro: ro}
}

func ObjectsFromJS(ro bool, values []js.Value) []Object {
	out := make([]Object, 0, len(values))
	for _, v := range values {
		out = append(out, Object{v: v, ro: ro})
	}
	return out
}

// Object wraps a js.Value
type Object struct {
	v  js.Value
	ro bool
}

func (o Object) SetMulti(keyVals ...interface{}) Object {
	if len(keyVals)%2 != 0 {
		panic("NewState: must supply pairs of (string, value)")
	}

	for i := 0; i < len(keyVals); i += 2 {
		o.Set(keyVals[i].(string), keyVals[i+1])
	}

	return o
}

func (o Object) Has(k string) bool {
	// have to use this so it wouldn't trigger react special properity checks
	return o.v.Call("hasOwnProperty", k).Bool()
}

func (o Object) Get(k string) Object {
	if !o.Valid() {
		return nullObj
	}

	return Object{v: o.v.Get(k)}
}

func (o Object) GetString(k string) string {
	return o.Get(k).String()
}

func (o Object) GetNumber(k string) float64 {
	return o.Get(k).Number()
}

// Set wraps js.Value.Set with (some) special handling for func(), func(Object), func(...Object), func(...Object) Object
// it doesn't panic
func (o Object) Set(k string, v interface{}) {
	if !o.Valid() {
		panic("nil object")
	}

	switch rv := v.(type) {
	case func(...Object) Object:
		v = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			return rv(ObjectsFromJS(false, args)...)
		})
	case func(...Object):
		v = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			rv(ObjectsFromJS(false, args)...)
			return nil
		})
	case func():
		v = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			rv()
			return nil
		})
	case func(Object):
		v = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			rv(Object{v: args[0]})
			return nil
		})
	case func(Object) Object:
		v = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			return rv(Object{v: args[0]})
		})
	default:
		v = ValueOf(v)
	}

	o.v.Set(k, v)
}

func (o Object) IsReadOnly() bool { return o.ro }
func (o Object) ReadOnly() Object { return Object{v: o.v, ro: true} }

func (o Object) IsString() bool { return o.v.Type() == js.TypeString }
func (o Object) String() string {
	return toString(o.v)
}

func (o Object) IsNumber() bool { return o.v.Type() == js.TypeNumber }
func (o Object) Number() float64 {
	if o.IsNumber() {
		return o.v.Float()
	}
	return 0
}

func (o Object) IsArray() bool { return IsArray(o.v) }
func (o Object) Array() Array  { return Array{o.v} }

func (o Object) Copy() Object {
	return Object{v: object.Call("assign", object.New(), o.v), ro: false}
}

func (o Object) Valid() bool {
	return !IsNull(o)
}

func (o Object) JSValue() js.Value { return o.v }
