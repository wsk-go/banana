package pointer

import "reflect"

// IsNil we use reflection to determine if a value is nil
// Checking whether something is nil is done by checking the value and type inside the interface
// It is only considered nil if both the value and type are nil.
// for example
// var a []string = nil
// var b any = a
// b == nil (false) Because the type of b is []string, it cannot be determined as nil
// Generally, we determine if variable is nil by examining its value
func IsNil(i any) bool {
	if i == nil {
		return true
	}

	value := reflect.ValueOf(i)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer:
		fallthrough
	case reflect.Interface, reflect.Slice:
		return value.IsNil()
	}
	return false
}

// IsZero we use reflection to determine if a value is zero
func IsZero(i any) bool {
	value := reflect.ValueOf(i)
	return value.IsZero()
}
