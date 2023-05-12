package validator

import (
	"github.com/go-playground/locales"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

var defaultValidator, _ = New()

func AddTranslation(translator locales.Translator, register func(*validator.Validate, ut.Translator) error) error {
	return defaultValidator.addTranslation(translator, register)
}

func Struct(obj any) error {
	return defaultValidator.Struct(obj)
}

func Structs[T any](objs []T) error {
	for _, obj := range objs {
		err := defaultValidator.Struct(obj)
		if err != nil {
			return err
		}
	}
	return nil
}

func StructWithLocale(obj any, locale ...string) error {
	return defaultValidator.StructWithLocale(obj, locale...)
}

func Var(field any, tag string) error {
	return defaultValidator.Var(field, tag)
}

func VarWithLocale(field any, tag string, locale ...string) error {
	return defaultValidator.VarWithLocale(field, tag, locale...)
}

func Default() *Validator {
	return defaultValidator
}
