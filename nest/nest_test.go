package nest

import (
	"testing"
)

type User struct {
	Name string
}

type TestBean struct {
	User  *User `inject:""`
	User2 *User `inject:"user2"`
}

func TestRegister(t *testing.T) {
	var nest Nest

	testBean := &TestBean{}

	err := nest.Register(&Bean{
		Value: testBean,
	}, &Bean{
		Value: &User{Name: "user"},
	}, &Bean{
		Value: &User{Name: "user2"},
		Name:  "user2",
	})

	if err != nil {
		t.Fatal(err)
	}

	err = nest.Run()
	if err != nil {
		t.Fatal(err)
	}

}
