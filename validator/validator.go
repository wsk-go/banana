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

//func addTranslation(translator locales.Translator,
//	uTranslator *ut.UniversalTranslator,
//	register func(trans ut.Translator) error) error {
//	err := uTranslator.AddTranslator(translator, true)
//	if err != nil {
//		return err
//	}
//	if trans, ok := uTranslator.GetTranslator(translator.Locale()); ok {
//		return register(trans)
//	}
//	return nil
//}
//
//type Language string
//
//func bindTranslators(validate *validator.Validate, defaultLanguage Language, languages ...Language) *ut.UniversalTranslator {
//	var uni *ut.UniversalTranslator
//
//	var zhTranslator ut.Translator
//	var enTranslator ut.Translator
//
//	zhTrans := zh.New()
//	enTrans := en.New()
//
//	// 第一个参数是默认翻译
//	uni = ut.New(zhTrans, zhTrans, enTrans)
//
//	if trans, ok := uni.GetTranslator("zh"); ok {
//		zhTranslator = trans
//	}
//
//	if trans, ok := uni.GetTranslator("en"); ok {
//		enTranslator = trans
//	}
//
//	err := zhtranslations.RegisterDefaultTranslations(validate, zhTranslator)
//
//	if err != nil {
//		panic(err)
//	}
//
//	err = entranslations.RegisterDefaultTranslations(validate, enTranslator)
//
//	if err != nil {
//		panic(err)
//	}
//
//	return uni
//}

var enumTranslationText = map[string]string{
	"en": "{0} is invalid！",
	"zh": "{0} 不合法!",
}

type StructValidator interface {
	ValidateStruct(obj interface{}) error
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

// ValidateStruct receives any kind of type, but only performed struct or pointer to struct type.
func (th *Validator) ValidateStruct(obj any, local ...string) error {
	value := reflect.ValueOf(obj)
	valueType := value.Kind()
	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	if valueType == reflect.Struct {
		if err := th.validate.Struct(obj); err != nil {
			ve := err.(validator.ValidationErrors)
			for _, vee := range ve {
				if th.uTranslator != nil {
					var trans ut.Translator
					var found bool
					for _, s := range local {
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
		}
	}
	return nil
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
