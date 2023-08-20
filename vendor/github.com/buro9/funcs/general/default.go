package general

import (
	"fmt"
	"reflect"
	"time"
)

// Default checks whether a given value is set and returns a default value if it
// is not.  "Set" in this context means non-zero for numeric types and times;
// non-zero length for strings, arrays, slices, and maps;
// any boolean or struct value; or non-nil for any other types.
//
// Copyright The Hugo Authors
// License Apache 2
// https://github.com/spf13/hugo/blob/master/LICENSE.md
// https://github.com/spf13/hugo/blob/master/tpl/template_funcs.go
func Default(dflt interface{}, given ...interface{}) (interface{}, error) {
	// given is variadic because the following construct will not pass a piped
	// argument when the key is missing:  {{ index . "key" | default "foo" }}
	// The Go template will complain that we got 1 argument when we expectd 2.

	if len(given) == 0 {
		return dflt, nil
	}
	if len(given) != 1 {
		return nil, fmt.Errorf("wrong number of args for default: want 2 got %d", len(given)+1)
	}

	g := reflect.ValueOf(given[0])
	if !g.IsValid() {
		return dflt, nil
	}

	set := false

	switch g.Kind() {
	case reflect.Bool:
		set = true
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
		set = g.Len() != 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		set = g.Int() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		set = g.Uint() != 0
	case reflect.Float32, reflect.Float64:
		set = g.Float() != 0
	case reflect.Complex64, reflect.Complex128:
		set = g.Complex() != 0
	case reflect.Struct:
		switch actual := given[0].(type) {
		case time.Time:
			set = !actual.IsZero()
		default:
			set = true
		}
	default:
		set = !g.IsNil()
	}

	if set {
		return given[0], nil
	}

	return dflt, nil
}
