package validator

import (
	"github.com/JackWSK/banana"
)

func Configuration() banana.ConfigurationFunc {
	return func(application banana.Application) (*banana.Configuration, error) {
		validator, err := NewValidator()
		if err != nil {
			panic(err)
		}
		return &banana.Configuration{
			Beans: []*banana.Bean{
				{
					Value: validator,
					Name:  "",
				},
			},
		}, nil
	}
}
