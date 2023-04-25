package nest

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"reflect"
)

func DefaultApp() *fiber.App {
	return fiber.New()
}

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

type MethodMapping struct {
	method reflect.Method
}

type Config struct {
	Engine *fiber.App
}

type Nest struct {
	beans       []*Bean
	controllers []*Bean
	named       map[string]*Bean
	typed       map[reflect.Type]*Bean

	engine *fiber.App
}

func New(config Config) *Nest {
	if config.Engine == nil {
		config.Engine = DefaultApp()
	}
	return &Nest{
		engine: config.Engine,
	}
}

func (th *Nest) Engine() *fiber.App {
	return th.engine
}

// Import Module
func (th *Nest) Import(modules ...*Module) error {
	for _, module := range modules {
		if len(module.Beans) > 0 {
			err := th.RegisterBean(module.Beans...)
			if err != nil {
				return err
			}
		}

		if len(module.Controllers) > 0 {
			err := th.RegisterController(module.Controllers...)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (th *Nest) RegisterController(beans ...*Bean) error {
	err := th.RegisterBean(beans...)
	if err != nil {
		return err
	}
	th.controllers = append(th.controllers, beans...)
	return nil
}

// RegisterBean register object
func (th *Nest) RegisterBean(beans ...*Bean) error {

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
		th.beans = append(th.beans, bean)
	}

	return nil
}

func (th *Nest) Run(addr string) error {
	err := th.inject()
	if err != nil {
		return err
	}

	err = th.callHook()
	if err != nil {
		return err
	}

	err = th.handleMapping()

	if err != nil {
		return err
	}

	return th.engine.Listen(addr)

}

func (th *Nest) inject() error {

	for _, bean := range th.beans {
		err := th.InjectOne(bean)
		if err != nil {
			return err
		}
	}

	return nil
}

func (th *Nest) callHook() error {

	for _, bean := range th.beans {
		if setup, ok := bean.Value.(BeanLoaded); ok {
			setup.Loaded()
		}
	}

	return nil
}

func (th *Nest) handleMapping() error {
	for _, controller := range th.controllers {
		for i := 0; i < controller.reflectType.NumMethod(); i++ {
			if isMappingMethod(controller.reflectType.Method(i)) {
				method := controller.reflectValue.Method(i)
				value := method.Call(nil)[0]
				if mapping, ok := value.Interface().(Mapping); ok {
					th.engine.Add(mapping.GetMethod(), mapping.GetPath(), mapping.GetHandler())
				}
			}
		}
	}

	return nil
}

func (th *Nest) InjectOne(bean *Bean) error {
	for i := 0; i < bean.reflectValue.Elem().NumField(); i++ {
		field := bean.reflectValue.Elem().Field(i)
		fieldType := field.Type()
		fieldTag := bean.reflectType.Elem().Field(i).Tag
		//fieldName := bean.reflectType.Elem().Field(i).Name
		tag, err := parseTag(fieldTag)

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

		var injectBean *Bean
		if tag.Name == "" {
			if ib, ok := th.typed[fieldType]; ok {
				injectBean = ib
			} else {
				return fmt.Errorf(
					"inject bean not found for field %s in type %s",
					bean.reflectType.Elem().Field(i).Name,
					bean.reflectType,
				)
			}
		} else {
			if ib, ok := th.named[tag.Name]; ok {
				injectBean = ib
			} else {
				return fmt.Errorf(
					"inject bean not found for field %s in name %s",
					bean.reflectType.Elem().Field(i).Name,
					bean.reflectType,
				)
			}
		}

		field.Set(injectBean.reflectValue)
	}

	return nil
}
