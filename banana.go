package banana

import (
	"errors"
	"fmt"
	"github.com/JackWSK/banana/defines"
	"github.com/JackWSK/banana/utils/stream"
	"reflect"
)

type MethodMapping struct {
	method reflect.Method
}

type Config struct {
	Engine defines.Engine
}

type Banana struct {
	beans       []*defines.Bean
	controllers []*defines.Bean
	named       map[string]*defines.Bean
	typed       map[reflect.Type]*defines.Bean

	engine defines.Engine
}

func New(config Config) *Banana {
	return &Banana{
		engine: config.Engine,
	}
}

func (th *Banana) Engine() defines.Engine {
	return th.engine
}

func (th *Banana) GetBeanByType(t reflect.Type) any {
	if bean, ok := th.typed[t]; ok {
		return bean.Value
	}

	return nil
}

func (th *Banana) GetBeanByName(name string) any {
	if bean, ok := th.named[name]; ok {
		return bean.Value
	}

	return nil
}

func (th *Banana) Import(modules ...defines.ModuleFunc) error {
	for _, module := range modules {

		configuration, err := module(th)
		if err != nil {
			return err
		}

		if len(configuration.Beans) > 0 {
			err := th.RegisterBean(configuration.Beans...)
			if err != nil {
				return err
			}
		}

		if len(configuration.Controllers) > 0 {
			err := th.RegisterController(configuration.Controllers...)
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
	err := th.prepareBeans()

	if err != nil {
		return err
	}

	err = th.handleMapping()

	if err != nil {
		return err
	}

	return th.engine.Listen(addr)
}

func (th *Banana) prepareBeans() error {
	beans := stream.Of(th.beans).Filter(func(bean *defines.Bean) bool {
		return bean.Injected == false
	}).ToList()

	err := th.inject(beans)
	if err != nil {
		return err
	}

	err = th.callHook(beans)
	if err != nil {
		return err
	}
	return nil
}

func (th *Banana) inject(beans []*defines.Bean) error {

	for _, bean := range beans {
		err := th.InjectOne(bean)
		if err != nil {
			return err
		}
		bean.Injected = true
	}

	return nil
}

func (th *Banana) callHook(beans []*defines.Bean) error {

	var configurations []defines.ModuleFunc
	for _, b := range beans {
		if setup, ok := b.Value.(defines.BeanLoaded); ok {
			setup.Loaded()
		}

		// continue configuration
		//if beanConfiguration, ok := b.Value.(defines.BeanConfiguration); ok {
		//	configurations = append(configurations, beanConfiguration.Configuration())
		//}
	}

	if len(configurations) > 0 {
		err := th.Import(configurations...)
		if err != nil {
			return err
		}
		err = th.prepareBeans()
		if err != nil {
			return err
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

func MustGetBeanByType[T any](application defines.Application) T {
	if v, ok := GetBeanByType[T](application); ok {
		return v
	}
	panic(errors.New(fmt.Sprintf("%T bean not found", *new(T))))
}

func GetBeanByType[T any](application defines.Application) (T, bool) {
	var t T
	if v := application.GetBeanByType(reflect.TypeOf(t)); v != nil {
		return v.(T), true
	}
	return t, false
}

func MustGetBeanByName[T any](application defines.Application, name string) T {
	if v, ok := GetBeanByName[T](application, name); ok {
		return v
	}

	panic(errors.New(fmt.Sprintf("bean with name [%s] not found", name)))
}

func GetBeanByName[T any](application defines.Application, name string) (T, bool) {
	var t T
	if v := application.GetBeanByName(name); v != nil {
		if v2, ok := v.(T); ok {
			return v2, ok
		}
	}
	return t, false
}
