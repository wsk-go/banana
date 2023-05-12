package banana

import (
	"fmt"
	"reflect"
	"testing"
)

func TestMappingAssign(t *testing.T) {
	m := RequestMapping{
		Method:        "",
		Path:          "",
		Handler:       nil,
		RequiredQuery: nil,
	}

	fmt.Println(reflect.TypeOf(m).Implements(mappingType))
}
