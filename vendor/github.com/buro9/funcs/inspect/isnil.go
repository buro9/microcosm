package inspect

import "reflect"

// IsNil returns true if the provided variable is nil
func IsNil(v interface{}) bool {
	if v == nil {
		return true
	}

	vo := reflect.ValueOf(v)
	switch vo.Kind() {
	case reflect.String:
		// always ok
	default:
		if vo.IsNil() || !vo.IsValid() {
			return true
		}
	}

	return false
}
