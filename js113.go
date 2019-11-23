// +build !go1.14

package wjsu

import "syscall/js"

// IsNull checks if an o is Null or Undefined
func IsNull(o js.Wrapper) bool {
	return o.JSValue() == js.Null()|| o.JSValue() == js.Undefined()
}
