package banana

import (
	"errors"
	"fmt"
	"github.com/JackWSK/banana/defines"
	"github.com/gofiber/fiber/v2"
	"reflect"
)

func DefaultApp() *fiber.App {
	return fiber.New()
}

type MethodMapping struct {
	method reflect.Method
}

type Config struct {
	Engine *fiber.App
}

type Banana struct {
	beans       []*defines.Bean
	controllers []*defines.Bean
	named       map[string]*defines.Bean
	typed       map[reflect.Type]*defines.Bean

	engine *fiber.App
}

func New(config Config) *Banana {
	if config.Engine == nil {
		config.Engine = DefaultApp()
	}
	return &Banana{
		engine: config.Engine,
	}
}

func (th *Banana) Engine() *fiber.App {
	return th.engine
}

// Import Configuration
func (th *Banana) Import(modules ...*defines.Configuration) error {
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

func (th *Banana) RegisterController(beans ...*defines.Bean) error {
	err := th.RegisterBean(beans...)
	if err != nil {
		return err
	}
	th.controllers = append(th.controllers, beans...)
	return nil
}

// RegisterBean register object
func (th *Banana) RegisterBean(beans ...*defines.Bean) error {

	if th.named == nil {
		th.named = make(map[string]*defines.Bean)
	}

	if th.typed == nil {
		th.typed = make(map[reflect.Type]*defines.Bean)
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

		bean.ReflectValue = reflectValue
		bean.ReflectType = reflectType
		th.beans = append(th.beans, bean)
	}

	return nil
}

func (th *Banana) Run(addr string) error {
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

func (th *Banana) inject() error {

	for _, bean := range th.beans {
		err := th.InjectOne(bean)
		if err != nil {
			return err
		}
	}

	return nil
}

func (th *Banana) callHook() error {

	for _, b := range th.beans {
		if setup, ok := b.Value.(defines.BeanLoaded); ok {
			setup.Loaded()
		}
	}

	return nil
}

func (th *Banana) handleMapping() error {
	for _, controller := range th.controllers {
		for i := 0; i < controller.ReflectType.NumMethod(); i++ {
			if isMappingMethod(controller.ReflectType.Method(i)) {
				method := controller.ReflectValue.Method(i)
				value := method.Call(nil)[0]
				if mapping, ok := value.Interface().(Mapping); ok {
					th.engine.Add(mapping.GetMethod(), mapping.GetPath(), mapping.GetHandler())
				}
			}
		}
	}

	return nil
}

func (th *Banana) InjectOne(b *defines.Bean) error {
	for i := 0; i < b.ReflectValue.Elem().NumField(); i++ {
		field := b.ReflectValue.Elem().Field(i)
		fieldType := field.Type()
		fieldTag := b.ReflectType.Elem().Field(i).Tag
		//fieldName := bean.ReflectType.Elem().Field(i).Name
		tag, err := parseTag(fieldTag)

		if err != nil {
			return fmt.Errorf(
				"unexpected tag format `%s` for field %s in type %s",
				string(fieldTag),
				b.ReflectType.Elem().Field(i).Name,
				b.ReflectType,
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
				b.ReflectType.Elem().Field(i).Name,
				b.ReflectType,
			)
		}

		var injectBean *defines.Bean
		if tag.Name == "" {
			if ib, ok := th.typed[fieldType]; ok {
				injectBean = ib
			} else {
				return fmt.Errorf(
					"inject bean not found for field %s in type %s",
					b.ReflectType.Elem().Field(i).Name,
					b.ReflectType,
				)
			}
		} else {
			if ib, ok := th.named[tag.Name]; ok {
				injectBean = ib
			} else {
				return fmt.Errorf(
					"inject bean not found for field %s in name %s",
					b.ReflectType.Elem().Field(i).Name,
					b.ReflectType,
				)
			}
		}

		field.Set(injectBean.ReflectValue)
	}

	return nil
}
