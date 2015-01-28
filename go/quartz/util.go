package quartz

import "reflect"

// Given a struct type, creates a mapping of field name
// to string representation of the field name's type.
func structFieldToType(t reflect.Type) map[string]string {
	fieldToType := make(map[string]string)
	for i := 0; i < t.NumField(); i++ {
		fieldToType[t.Field(i).Name] = t.Field(i).Type.String()
	}
	return fieldToType
}
