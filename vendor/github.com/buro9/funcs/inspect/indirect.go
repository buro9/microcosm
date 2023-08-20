package inspect

import "reflect"

// indirect returns the item at the end of indirection, and a bool to indicate if it's nil.
//
// indirect is taken from https://golang.org/src/text/template/exec.go
func indirect(v reflect.Value) (rv reflect.Value, isNil bool) {
	for ; v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface; v = v.Elem() {
		if v.IsNil() {
			return v, true
		}
	}
	return v, false
}
