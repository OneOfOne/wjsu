//go:build js && wasm
// +build js,wasm

package wjsu

import (
	"reflect"
	"syscall/js"
)

type jsFuncType = func(js.Value, []js.Value) interface{}

var Console ConsoleImpl

func RawConsole() js.Value { return console }

type ConsoleImpl struct{}

func (ConsoleImpl) Log(args ...interface{})   { console.Call("log", ArgsToObjects(args)...) }
func (ConsoleImpl) Warn(args ...interface{})  { console.Call("warn", ArgsToObjects(args)...) }
func (ConsoleImpl) Error(args ...interface{}) { console.Call("error", ArgsToObjects(args)...) }

func ValueOf(x interface{}) (o Object) {
	if err := tryCall(func() { o = Object{v: js.ValueOf(x)} }); err == "" {
		return o
	}
	rv := reflect.Indirect(reflect.ValueOf(x))

	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		a := NewArray(rv.Len())
		for i := 0; i < rv.Len(); i++ {
			a.Set(i, ValueOf(rv.Index(i).Interface()))
		}
		o = a.Object()

	case reflect.Map:
		o = NewObject(false)
		for it := rv.MapRange(); it.Next(); {
			o.Set(it.Key().String(), ValueOf(it.Value().Interface()))
		}
	case reflect.Func:
		switch fn := rv.Interface().(type) {
		case jsFuncType:
			jfn := js.FuncOf(fn)
			o = Object{v: jfn, ro: true}
		default:
			panic("unsupported func type: " + rv.Type().String() + ", use func(js.Value, []js.Value) interface{}")
		}
	default:
		panic("unsupported type: " + rv.Type().String() + ", maybe use a pointer?")
	}

	return o
}

func WrapFunc(fn func(), once bool) (jfn js.Func) {
	jfn = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if once {
			jfn.Release()
		}
		fn()
		return nil
	})

	return
}

func ArgsToObjects(in []interface{}) []interface{} {
	out := in[:0]
	for _, v := range in {
		out = append(out, ValueOf(v))
	}
	return out
}

func tryCall(fn func()) (err string) {
	defer func() {
		if r := recover(); r != nil {
			err, _ = r.(string)
		}
	}()
	fn()
	return
}

func toString(v js.Value) string {
	if v.Type() == js.TypeString {
		return v.String()
	}

	return ""
}
