package banana

import "reflect"

type Bean struct {
	// the object you try to register
	Value any

	// the name you register in banana
	Name string

	// type
	ReflectType reflect.Type

	// value
	ReflectValue reflect.Value

	// has injected
	Injected bool
}
