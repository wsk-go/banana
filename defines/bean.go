package defines

import "reflect"

type Bean struct {
	// the object you try to register
	Value any

	// the name you register in nest
	Name string

	// type
	ReflectType reflect.Type

	// value
	ReflectValue reflect.Value
}
