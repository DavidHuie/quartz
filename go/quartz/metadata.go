package quartz

import "reflect"

type structMetadata struct {
	NameToMethodMetadata map[string]*methodMetadata
	TargetStruct         interface{} `json:"-"`
}

func newStructMetadata(targetStruct interface{}) *structMetadata {
	return &structMetadata{
		make(map[string]*methodMetadata),
		targetStruct,
	}
}

type methodMetadata struct {
	Method         reflect.Method `json:"-"`
	ArgumentToType map[string]string
}

// structFieldToType creates a mapping of field name to string
// representation of the field name's type.
func structFieldToType(t reflect.Type) map[string]string {
	fieldToType := make(map[string]string)
	for i := 0; i < t.NumField(); i++ {
		fieldToType[t.Field(i).Name] = t.Field(i).Type.String()
	}
	return fieldToType
}
