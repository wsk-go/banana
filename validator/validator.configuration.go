package validator

import "github.com/JackWSK/banana/defines"

func Configuration() *defines.Configuration {
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
	}
}
