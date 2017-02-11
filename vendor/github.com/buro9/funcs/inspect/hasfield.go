package inspect

import (
	reflections "gopkg.in/oleiade/reflections.v1"
)

// HasField checks if the provided field name is part of a struct.
//
// The stuct can be a structure or pointer to structure.
func HasField(s interface{}, fieldName string) bool {
	has, _ := reflections.HasField(s, fieldName)
	return has
}
