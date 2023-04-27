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
	Gender UserGender `validate:"enum"`
}

func TestValidator(t *testing.T) {
	validator, err := NewValidator()
	if err != nil {
		t.Fatal(err)
	}

	var user = User{
		Name: "jack",
	}

	err = validator.ValidateStruct(user)
	fmt.Printf("%+v", err)
}
