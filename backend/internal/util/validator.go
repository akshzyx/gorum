package util

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(s interface{}) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	verrs := err.(validator.ValidationErrors)
	msg := ""

	for _, e := range verrs {
		msg += fmt.Sprintf("field '%s' failed on the '%s' rule; ", e.Field(), e.Tag())
	}

	return fmt.Errorf("%s", msg)
}
