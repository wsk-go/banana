package validator

import (
	"fmt"
	"testing"
)

type UserGender string

func (th UserGender) IsValid() bool {
	switch th {
	case UserGenderMan, UserGenderFemale:
		return true
	}
	return false
}

const (
	UserGenderMan    UserGender = "man"
	UserGenderFemale UserGender = "female"
)

type User struct {
	Name   string     `validate:"required" json:"name"`
	Gender UserGender `validate:"enum" json:"gender"`
}

func TestValidator(t *testing.T) {
	validator, err := New()
	if err != nil {
		t.Fatal(err)
	}

	var user = User{
		//Name: "jack",
	}

	err = validator.StructWithLocale(user)

	fmt.Printf("%+v", err)
}
