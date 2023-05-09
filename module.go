package banana

import (
	"reflect"
)

type Application interface {
	GetBeanByType(t reflect.Type) any
	GetBeanByName(name string) any
}

type ConfigurationFunc func(Application) (*Configuration, error)

type Configuration struct {
	Controllers []*Bean
	Beans       []*Bean
}
