package utils

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func Validate(obj any) (int, string) {
	validate := validator.New()

	if err := validate.Struct(obj); err != nil {
		var ve validator.ValidationErrors
		var msg = ""

		if errors.As(err, &ve) {
			msg = "data not valid : " + ve[0].Field()
		}

		return http.StatusBadRequest, msg
	}

	return 0, ""
}
