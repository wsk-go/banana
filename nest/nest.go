package nest

import (
	"errors"
	"fmt"
	"reflect"
)

type Bean struct {
	// the object you try to register
	Value any

	// the name you register in nest
	Name string

	// type
	reflectType reflect.Type
	// value
	reflectValue reflect.Value
}

type Nest struct {
	named map[string]*Bean
	typed map[reflect.Type]*Bean
}

// Register register object
func (th *Nest) Register(beans ...*Bean) error {

	if th.named == nil {
		th.named = make(map[string]*Bean)
	}

	if th.typed == nil {
		th.typed = make(map[reflect.Type]*Bean)
	}

	for _, bean := range beans {
		reflectType := reflect.TypeOf(bean.Value)
		reflectValue := reflect.ValueOf(bean.Value)

		if !isStructPtr(reflectType) {
			return errors.New(fmt.Sprintf("Bean[%s] value must be ptr", reflectType.Name()))
		}

		if bean.Name == "" {
			if _, ok := th.typed[reflectType]; ok {
				return errors.New(fmt.Sprintf("the type %s duplicate", reflectType.String()))
			}
			th.typed[reflectType] = bean
		} else {
			if _, ok := th.named[bean.Name]; ok {
				return errors.New(fmt.Sprintf("the name %s duplicate", bean.Name))
			}
			th.named[bean.Name] = bean
		}

		bean.reflectValue = reflectValue
		bean.reflectType = reflectType
	}

	return nil
}

func (th *Nest) Run() error {
	th.Inject()

	return nil
}

func (th *Nest) Inject() {

}

func (th *Nest) InjectOne(bean *Bean) error {
	for i := 0; i < bean.reflectValue.Elem().NumField(); i++ {
		field := bean.reflectValue.Elem().Field(i)
		fieldType := field.Type()
		fieldTag := bean.reflectType.Elem().Field(i).Tag
		//fieldName := bean.reflectType.Elem().Field(i).Name
		tag, err := parseTag(string(fieldTag))

		if err != nil {
			return fmt.Errorf(
				"unexpected tag format `%s` for field %s in type %s",
				string(fieldTag),
				bean.reflectType.Elem().Field(i).Name,
				bean.reflectType,
			)
		}

		// Skip fields without a tag.
		if tag == nil {
			continue
		}

		// Cannot be used with unexported fields.
		if !field.CanSet() {
			return fmt.Errorf(
				"inject requested on unexported field %s in type %s",
				bean.reflectType.Elem().Field(i).Name,
				bean.reflectType,
			)
		}

		var injectBean any
		if tag.Name == "" {
			if ib, ok := th.typed[fieldType]; ok {
				injectBean = ib
			}
			return fmt.Errorf(
				"inject bean not found for field %s in type %s",
				bean.reflectType.Elem().Field(i).Name,
				bean.reflectType,
			)
		} else {
			if ib, ok := th.named[tag.Name]; ok {
				injectBean = ib
			}
		}

		field.Set(reflect.ValueOf(injectBean))
	}

	return nil
}
