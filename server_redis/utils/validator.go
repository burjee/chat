package utils

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var alphanumeric validator.Func = func(fl validator.FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if ok {
		matched, err := regexp.MatchString(`^[a-zA-Z0-9_]+$`, s)
		return err == nil && matched
	}
	return false
}

func RegisterValidation(validate *validator.Validate) {
	validate.RegisterValidation("alphanumeric", alphanumeric)
}
