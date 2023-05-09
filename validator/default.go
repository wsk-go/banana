package validator

import (
	"github.com/go-playground/locales"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

var defaultValidator, _ = NewValidator()

func AddTranslation(translator locales.Translator, register func(*validator.Validate, ut.Translator) error) error {
	return defaultValidator.AddTranslation(translator, register)
}

func Struct(obj any) error {
	return defaultValidator.Struct(obj)
}

func StructWithLocale(obj any, locale ...string) {

}

func Default() *Validator {
	return defaultValidator
}
