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
	"strings"
)

type Local string

const (
	// 中国简体
	LocalZH = "zh"
	// 英文
	LocalEN = "en"
)

func bindTranslators(validate *validator.Validate) *ut.UniversalTranslator {
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

func (th *Validator) FindTranslator(lan string) (t ut.Translator, found bool) {
	s := strings.Split(lan, "-")
	if len(s) > 0 {
		return th.ut.FindTranslator(s[0])
	} else {
		return th.ut.FindTranslator("zh")
	}

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
			return errors.WithStack(err)
		}
	}
	return nil
}

// Engine returns the underlying validator engine which powers the default
// Validator instance. This is useful if you want to register custom validations
// or struct level validations. See validator GoDoc for more info -
// https://godoc.org/gopkg.in/go-playground/validator.v8
func (th *Validator) Engine() interface{} {
	return th.validate
}
