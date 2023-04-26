package zaplogger

import (
	"github.com/JackWSK/banana/defines"
)

func Configuration(config LoggerConfig) *defines.Configuration {
	return &defines.Configuration{
		Beans: []*defines.Bean{
			{
				Value: NewLogger(config),
				Name:  "",
			},
		},
	}
}
