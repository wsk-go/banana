package validator

import (
	"github.com/JackWSK/banana/errors"
	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
	"reflect"
)

var enumTranslationText = map[string]string{
	"en": "{0} is invalid！",
	"zh": "{0} 不合法!",
}

type StructValidator interface {
	Struct(obj interface{}) error
}

type Validator struct {
	validate    *validator.Validate
	uTranslator *ut.UniversalTranslator
}

type Enum interface {
	IsValid() bool
}

func ValidateEnum(fl validator.FieldLevel) bool {
	value := fl.Field().Interface().(Enum)
	return value.IsValid()
}

func NewValidator() (*Validator, error) {
	validate := validator.New()

	// register enum
	err := validate.RegisterValidation("enum", ValidateEnum)
	if err != nil {
		return nil, err
	}

	// register tag
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		param := fld.Tag.Get("param")

		if param == "" {
			param = fld.Tag.Get("json")

			if param == "" {
				param = fld.Name
			}
		}

		return param
	})

	v := &Validator{validate: validate}
	err = v.AddTranslation(en.New(), entranslations.RegisterDefaultTranslations)
	return v, err
}

// Struct receives any kind of type, but only performed struct or pointer to struct type.
func (th *Validator) Struct(obj any) error {
	value := reflect.ValueOf(obj)
	valueType := value.Kind()
	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	if valueType == reflect.Struct {
		return th.validate.Struct(obj)
	}
	return nil
}

func (th *Validator) Translate(err error, locale ...string) error {
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, vee := range ve {
			if th.uTranslator != nil {
				var trans ut.Translator
				var found bool
				for _, s := range locale {
					trans, found = th.uTranslator.GetTranslator(s)
					if found {
						break
					}
				}

				if !found {
					trans = th.uTranslator.GetFallback()
				}

				message := vee.Translate(trans)
				return errors.NewValidationError(message)
			} else {
				return errors.NewValidationError(vee.Error())
			}
		}
	} else {
		return ve
	}
}

func (th *Validator) AddTranslation(translator locales.Translator, register func(*validator.Validate, ut.Translator) error) error {

	if th.uTranslator == nil {
		th.uTranslator = ut.New(translator)
	}

	err := th.uTranslator.AddTranslator(translator, true)
	if err != nil {
		return err
	}
	if trans, ok := th.uTranslator.GetTranslator(translator.Locale()); ok {
		text := enumTranslationText[translator.Locale()]
		if text == "" {
			text = enumTranslationText["en"]
		}
		err = th.validate.RegisterTranslation("enum", trans, func(ut ut.Translator) error {
			return ut.Add("enum", text, true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("enum", fe.Field())
			return t
		})

		if err != nil {
			return err
		}

		return register(th.validate, trans)
	}
	return nil
}
