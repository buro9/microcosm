package inspect

import "reflect"

// IsSet returns whether a given array, channel, slice, or map has a key
// defined.
//
// Copyright The Hugo Authors and covered by both an MIT license for
// the original code, and an Apache license for later modifications.
// https://github.com/spf13/hugo/blob/master/tpl/template_funcs.go
func IsSet(a interface{}, key interface{}) bool {
	av := reflect.ValueOf(a)
	kv := reflect.ValueOf(key)

	switch av.Kind() {
	case reflect.Array, reflect.Chan, reflect.Slice:
		if int64(av.Len()) > kv.Int() {
			return true
		}
	case reflect.Map:
		if kv.Type() == av.Type().Key() {
			return av.MapIndex(kv).IsValid()
		}
	}

	return false
}
