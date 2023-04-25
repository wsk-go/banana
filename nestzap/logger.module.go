package nestzap

import (
	"github.com/JackWSK/go-nest/nest"
)

func Module(config LoggerConfig) *nest.Module {
	return &nest.Module{
		Beans: []*nest.Bean{
			{
				Value: NewLogger(config),
				Name:  "",
			},
		},
	}
}
