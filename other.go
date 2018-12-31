//+build js,wasm

package jsenv

import (
	"fmt"
	"reflect"
	"syscall/js"
)

var (
	console = js.Global().Get("console")

	Console con
)

func RawConsole() js.Value { return console }

type con struct{}

func (con) Log(args ...interface{})   { console.Call("log", safeLogArgs(args)...) }
func (con) Warn(args ...interface{})  { console.Call("warn", safeLogArgs(args)...) }
func (con) Error(args ...interface{}) { console.Call("error", safeLogArgs(args)...) }

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
			Console.Error(fmt.Sprintf("unsupported type: %T (%v)", x, x))
		}
	}
	return
}

func safeLogArgs(in []interface{}) []interface{} {
	out := in[:0]
	for _, v := range in {
		out = append(out, ValueOf(v))
	}
	return out
}

func tryCall(fn func()) (err string) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Sprint(r)
		}
	}()
	fn()
	return
}
