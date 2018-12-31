//+build js,wasm

package wjsu

import "syscall/js"

var (
	array = js.Global().Get("Array")

	_ js.Wrapper = (*Array)(nil)
)

func RawArray() js.Value { return array }

func IsArray(v js.Value) bool {
	return array.Call("isArray", v).Bool()
}

func NewArray(len int) Array {
	return Array{v: array.New(len)}
}

type Array struct {
	v js.Value
}

func (a Array) Len() int { return a.v.Length() }

func (a Array) Get(i int) Object { return Object{v: a.v.Index(i)} }

func (a Array) Set(i int, v interface{}) { a.v.SetIndex(i, ValueOf(v)) }

func (a Array) Push(args ...interface{}) int {
	for i := range args {
		args[i] = ValueOf(args[i]).v
	}
	return a.v.Call("push", args...).Int()
}

func (a Array) Concat(args ...interface{}) Array {
	return a
}

func (a Array) Object() Object { return Object{v: a.v} }

func (a Array) JSValue() js.Value { return a.v }
