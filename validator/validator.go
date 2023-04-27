package validator

import (
	"github.com/JackWSK/banana/errors"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
	zhtranslations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
)

var translationMapping = map[Language]func(v *validator.Validate, trans ut.Translator) error{
	LanguageZH: zhtranslations.RegisterDefaultTranslations,
	LanguageEN: entranslations.RegisterDefaultTranslations,
}

type Language string

const (
	// LanguageZH chinese
	LanguageZH Language = "zh"
	// LanguageEN english
	LanguageEN Language = "en"
)

func bindTranslators(validate *validator.Validate, languages ...Language) *ut.UniversalTranslator {
	var uni *ut.UniversalTranslator
	var zhTranslator ut.Translator
	var enTranslator ut.Translator

	zhTrans := zh.New()
	enTrans := en.New()

	// 第一个参数是默认翻译
	uni = ut.New(zhTrans, zhTrans, enTrans)

	if trans, ok := uni.GetTranslator("zh"); ok {
		zhTranslator = trans
	}

	if trans, ok := uni.GetTranslator("en"); ok {
		enTranslator = trans
	}

	err := zhtranslations.RegisterDefaultTranslations(validate, zhTranslator)

	if err != nil {
		panic(err)
	}

	err = entranslations.RegisterDefaultTranslations(validate, enTranslator)

	if err != nil {
		panic(err)
	}

	return uni
}

type StructValidator interface {
	ValidateStruct(obj interface{}) error
}

type Validator struct {
	validate *validator.Validate
	ut       *ut.UniversalTranslator
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

	ut := bindTranslators(validate)
	return &Validator{validate: validate, ut: ut}, nil
}

// ValidateStruct receives any kind of type, but only performed struct or pointer to struct type.
func (th *Validator) ValidateStruct(obj any) error {
	value := reflect.ValueOf(obj)
	valueType := value.Kind()
	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	if valueType == reflect.Struct {
		if err := th.validate.Struct(obj); err != nil {
			ve := err.(validator.ValidationErrors)
			for _, vee := range ve {
				trans, found := th.ut.GetTranslator("zh")
				if !found {
					trans = th.ut.GetFallback()
				}
				message := vee.Translate(trans)
				return errors.NewValidationError(message)
			}
		}
	}
	return nil
}
