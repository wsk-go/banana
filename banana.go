package banana

import (
	"fmt"
	"github.com/wsk-go/banana/errors"
	"github.com/wsk-go/banana/hook"
	"github.com/wsk-go/banana/utils/stream"
	"reflect"
)

type MethodMapping struct {
	method reflect.Method
}

type Config struct {
	Engine Engine
}

type MiddlewareFunc func(ctx Context, application *Banana) error

type Middleware interface {
	Handle(ctx Context, application *Banana) error
}

type Banana struct {
	beans       []*Bean
	controllers []*Bean
	named       map[string]*Bean
	typed       map[reflect.Type]*Bean

	engine Engine
}

func New(config Config) *Banana {
	return &Banana{
		engine: config.Engine,
	}
}

func (th *Banana) Engine() Engine {
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

func (th *Banana) Use(fs ...MiddlewareFunc) {
	for _, f := range fs {
		th.engine.Use(func(ctx Context) error {
			return f(ctx, th)
		})
	}
}

func (th *Banana) RegisterMiddleware(middlewares ...Middleware) error {
	for _, middleware := range middlewares {
		th.engine.Use(func(ctx Context) error {
			return middleware.Handle(ctx, th)
		})
	}
	// convert to beans and register
	beans := stream.Map(stream.Of(middlewares), func(in Middleware) *Bean {
		return &Bean{
			Value: in,
			Name:  "",
		}
	})
	return th.RegisterBean(beans...)
}

func (th *Banana) Import(modules ...ConfigurationFunc) error {
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

func (th *Banana) RegisterController(beans ...*Bean) error {
	err := th.RegisterBean(beans...)
	if err != nil {
		return err
	}
	th.controllers = append(th.controllers, beans...)
	return nil
}

// RegisterBean register object
func (th *Banana) RegisterBean(beans ...*Bean) error {

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
	beans := stream.Of(th.beans).Filter(func(bean *Bean) bool {
		return bean.injected == false
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

func (th *Banana) inject(beans []*Bean) error {

	for _, bean := range beans {
		err := th.injectOne(bean)
		if err != nil {
			return err
		}
		bean.injected = true
	}

	return nil
}

func (th *Banana) callHook(beans []*Bean) error {

	var configurations []ConfigurationFunc
	for _, b := range beans {
		if setup, ok := b.Value.(hook.BeanLoaded); ok {
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
		for i := 0; i < controller.reflectType.NumMethod(); i++ {
			if isMappingMethod(controller.reflectType.Method(i)) {
				method := controller.reflectValue.Method(i)
				value := method.Call(nil)[0]
				if mapping, ok := value.Interface().(Mapping); ok {
					handler := mapping.GetHandler()
					th.engine.Add(mapping.GetMethod(), mapping.GetPath(), func(context Context) error {
						if len(mapping.GetRequiredQuery()) > 0 {
							for _, key := range mapping.GetRequiredQuery() {
								if v := context.Query(key); v == "" {
									return errors.New(fmt.Sprintf("%s not exists in query", key))
								}
							}
						}

						return handler(context)
					})
				}
			}
		}
	}

	return nil
}

func (th *Banana) injectOne(b *Bean) error {
	for i := 0; i < b.reflectValue.Elem().NumField(); i++ {
		field := b.reflectValue.Elem().Field(i)
		fieldType := field.Type()
		fieldTag := b.reflectType.Elem().Field(i).Tag
		//fieldName := bean.reflectType.Elem().Field(i).Name
		tag, err := parseTag(fieldTag)

		if err != nil {
			return fmt.Errorf(
				"unexpected tag format `%s` for field %s in type %s",
				string(fieldTag),
				b.reflectType.Elem().Field(i).Name,
				b.reflectType,
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
				b.reflectType.Elem().Field(i).Name,
				b.reflectType,
			)
		}

		var injectBean *Bean
		if tag.Name == "" {
			if ib, ok := th.typed[fieldType]; ok {
				injectBean = ib
			} else {
				return fmt.Errorf(
					"inject bean not found for field %s in type %s",
					b.reflectType.Elem().Field(i).Name,
					b.reflectType,
				)
			}
		} else {
			if ib, ok := th.named[tag.Name]; ok {
				injectBean = ib
			} else {
				return fmt.Errorf(
					"inject bean not found for field %s in name %s",
					b.reflectType.Elem().Field(i).Name,
					b.reflectType,
				)
			}
		}

		field.Set(injectBean.reflectValue)
	}

	return nil
}

func MustGetBeanByType[T any](application Application) T {
	if v, ok := GetBeanByType[T](application); ok {
		return v
	}
	panic(errors.New(fmt.Sprintf("%T bean not found", *new(T))))
}

func GetBeanByType[T any](application Application) (T, bool) {
	var t T
	if v := application.GetBeanByType(reflect.TypeOf(t)); v != nil {
		return v.(T), true
	}
	return t, false
}

func MustGetBeanByName[T any](application Application, name string) T {
	if v, ok := GetBeanByName[T](application, name); ok {
		return v
	}

	panic(errors.New(fmt.Sprintf("bean with name [%s] not found", name)))
}

func GetBeanByName[T any](application Application, name string) (T, bool) {
	var t T
	if v := application.GetBeanByName(name); v != nil {
		if v2, ok := v.(T); ok {
			return v2, ok
		}
	}
	return t, false
}
