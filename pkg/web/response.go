package web

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"
)

type Response struct {
	Code string      `json:"code"`
	Data interface{} `json:"data,omitempty"`
}

type Error struct {
	Code  string      `json:"code"`
	Error interface{} `json:"error,omitempty"`
}

type ErrorInValidations struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Message string `json:"message"`
}

func Simple(err error) []ErrorInValidations {
	var validatorErrors validator.ValidationErrors
	var errorsReturned []ErrorInValidations
	if errors.As(err, &validatorErrors) { // err.(validator.ValidationErrors)
		for _, f := range validatorErrors {
			errorsReturned = append(errorsReturned, ErrorInValidations{Field: f.Field(), Tag: f.Tag(), Message: fmt.Sprintf("el campo '%v' no cumple con la validación '%v'", f.Field(), f.Tag())})
		}
	}
	return errorsReturned
}

func NewResponse(code int, content interface{}) interface{} {
	if code < 300 {
		return Response{strconv.FormatInt(int64(code), 10), content}
	}

	return Error{strconv.FormatInt(int64(code), 10), content}
}

func ResponseUnauthorized() Error {
	return Error{strconv.FormatInt(int64(http.StatusUnauthorized), 10), "no tiene permisos para realizar la petición solicitada"}
}
