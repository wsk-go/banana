package zap

import (
	"github.com/JackWSK/go-nest/defines"
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
