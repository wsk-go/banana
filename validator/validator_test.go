package validator

import (
	"fmt"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	entranslations "github.com/go-playground/validator/v10/translations/en"
	zhtranslations "github.com/go-playground/validator/v10/translations/zh"
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
	validator, err := NewValidator()
	if err != nil {
		t.Fatal(err)
	}

	err = validator.AddTranslation(zh.New(), zhtranslations.RegisterDefaultTranslations)
	if err != nil {
		t.Fatal(err)
	}

	err = validator.AddTranslation(en.New(), entranslations.RegisterDefaultTranslations)
	if err != nil {
		t.Fatal(err)
	}

	var user = User{
		Name: "jack",
	}

	err = validator.ValidateStruct(user, "en")
	fmt.Printf("%+v", err)
}
