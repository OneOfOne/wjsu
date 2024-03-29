//go:build go1.14
// +build go1.14

package wjsu

import "syscall/js"

// IsNull checks if an o is Null or Undefined
func IsNull(o js.Wrapper) bool {
	return o == nil || o.JSValue().IsNull() || o.JSValue().IsUndefined()
}
