package zap

import (
	"github.com/JackWSK/go-nest/defines"
	"github.com/JackWSK/go-nest/nest"
)

func Configuration(config LoggerConfig) *defines.Configuration {
	return &defines.Configuration{
		Beans: []*nest.Bean{
			{
				Value: NewLogger(config),
				Name:  "",
			},
		},
	}
}
