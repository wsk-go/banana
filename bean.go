package banana

import "reflect"

type Bean struct {
	// the object you try to register
	Value any

	// the name you register in banana
	Name string

	// type
	reflectType reflect.Type

	// value
	reflectValue reflect.Value

	// has injected
	injected bool
}
