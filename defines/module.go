package defines

import "reflect"

type Application interface {
	GetBeanByType(t reflect.Type) any
	GetBeanByName(name string) any
}

type ModuleFunc func(Application) (*Configuration, error)

type Configuration struct {
	Controllers []*Bean
	Beans       []*Bean
}
