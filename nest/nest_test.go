package nest

import (
	"fmt"
	"testing"
)

type User struct {
	Name string
}

type TestBean struct {
	User  *User `json:"user" inject:""`
	User2 *User `inject:"user2"`
}

func (t *TestBean) Loaded() {
	fmt.Println("loaded")
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
