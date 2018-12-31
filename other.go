//+build js,wasm

package wjsu

import (
	"reflect"
	"syscall/js"
)

var (
	console = js.Global().Get("console")

	Console ConsoleImpl
)

func RawConsole() js.Value { return console }

type ConsoleImpl struct{}

func (ConsoleImpl) Log(args ...interface{})   { console.Call("log", safeArgs(args)...) }
func (ConsoleImpl) Warn(args ...interface{})  { console.Call("warn", safeArgs(args)...) }
func (ConsoleImpl) Error(args ...interface{}) { console.Call("error", safeArgs(args)...) }

func ValueOf(x interface{}) (o Object) {
	if err := tryCall(func() { o = Object{v: js.ValueOf(x)} }); err != "" {
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

		default:
			Console.Error("unsupported type:", rv.Type().String())
		}
	}
	return
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

func safeArgs(in []interface{}) []interface{} {
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
