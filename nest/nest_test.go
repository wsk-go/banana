package nest

import (
	"testing"
)

type TestBean struct {
}

func TestRegister(t *testing.T) {
	var nest Nest

	err := nest.Register(&Bean{
		Value: &TestBean{},
	}, &Bean{
		Value: &TestBean{},
		Name:  "hello",
	})

	if err != nil {
		t.Fatal(err)
	}
}
