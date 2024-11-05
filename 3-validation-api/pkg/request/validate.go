package request

import "github.com/go-playground/validator/v10"

func Validate[T any](body T) error {
	validate := validator.New()
	err := validate.Struct(body)
	return err
}
