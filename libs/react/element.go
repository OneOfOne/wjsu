//+build js,wasm

package react

import (
	"strconv"
	"sync/atomic"
	"syscall/js"

	"github.com/OneOfOne/wjsu"
)

type Element interface {
	elem()
}

type reactElement struct{ js.Value }

func (reactElement) elem() {}

// E is a shorthand and wrapper for React.createElement
// tag can either be a Component, a func() Component or a string,
// if tag is nil it defaults to React.Fragment.
// if a key isn't specifically set in properties, it'll be automatically generated.
// usage:
// - E("div")
// - E("h1", "this is cool")
// - E(&MyComp{}, E("div"))
// - E(func() Component { return &MyComp{something: 1} }, O("style", O("color", "red"), E("h1", "something"))
func E(tag interface{}, childrenAndMaybeProps ...interface{}) Element {
	var props wjsu.Object
	if len(childrenAndMaybeProps) > 0 {
		if p, ok := childrenAndMaybeProps[0].(wjsu.Object); ok {
			props = p
			childrenAndMaybeProps = childrenAndMaybeProps[1:]
		}
	}

	if !props.Valid() {
		props = wjsu.NewObject(false)
	}

	if !props.Has("key") {
		props.Set("key", "auto:"+strconv.Itoa(int(atomic.AddInt32(&eIdx, 1))))
	}

	children := childrenAndMaybeProps[:0]
	for _, c := range childrenAndMaybeProps {
		if ele := makeElement(c, true); ele != js.Null() {
			children = append(children, ele)
		}
	}

	if cc := props.Get("children"); cc.IsArray() {
		for i, a := 0, cc.Array(); i < a.Len(); i++ {
			children = append(children, a.Get(i))
		}

		props.Set("children", js.Undefined())
	}

	if tag == nil {
		tag = Fragment()
	} else {
		tag = makeElement(tag, false)
	}

	return reactElement{RawCreateElement(tag, props, children)}
}

func makeElement(v interface{}, invoke bool) interface{} {
	switch any := v.(type) {
	case Component:
		return createCtor(func() Component { return any }, invoke)
	case func() Component:
		return createCtor(any, invoke)
	case StatelessComponent:
		return wrapFunc(any, invoke)
	case func() Element:
		return wrapFunc(func(wjsu.Object) Element { return any() }, invoke)
	default:
		return v
	}
}
