package zaplogger

import (
	"github.com/JackWSK/banana/defines"
)

func Configuration(config LoggerConfig) defines.ModuleFunc {

	return func(application defines.Application) (*defines.Configuration, error) {
		return &defines.Configuration{
			Beans: []*defines.Bean{
				{
					Value: NewLogger(config),
					Name:  "",
				},
			},
		}, nil
	}

}
