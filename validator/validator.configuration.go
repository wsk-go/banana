package validator

import "github.com/JackWSK/banana/defines"

func Configuration() defines.ModuleFunc {
	return func(application defines.Application) (*defines.Configuration, error) {
		validator, err := NewValidator()
		if err != nil {
			panic(err)
		}
		return &defines.Configuration{
			Beans: []*defines.Bean{
				{
					Value: validator,
					Name:  "",
				},
			},
		}, nil
	}
}
