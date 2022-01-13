package web

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strings"
)

type Response struct {
	Data interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Status  int         `json:"-"`
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Fields  interface{} `json:"fields,omitempty"`
}

func response(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}

func Success(c *gin.Context, status int, data interface{}) {
	response(c, status, Response{Data: data})
}

func Error(c *gin.Context, status int, format string, args ...interface{}) {
	err := ErrorResponse{
		Code:    StatusMsg(status),
		Message: fmt.Sprintf(format, args...),
		Status:  status,
	}

	response(c, status, err)
}

func ValidationError(c *gin.Context, status int, err error) {
	errorResponse := ErrorResponse{
		Code:    StatusMsg(status),
		Message: InvalidFields,
		Fields:  printValidationError(err),
		Status:  status,
	}

	response(c, status, errorResponse)
}

func printValidationError(err error) interface{} {
	type validationError struct {
		Field   string `json:"field"`
		Tag     string `json:"tag"`
		Message string `json:"message"`
	}

	var validatorErrors validator.ValidationErrors
	var errorsReturned []validationError
	if errors.As(err, &validatorErrors) {
		for _, f := range validatorErrors {
			errorsReturned = append(errorsReturned,
				validationError{
					Field:   f.Field(),
					Tag:     f.Tag(),
					Message: FieldMissing,
				},
			)
		}
	}
	return errorsReturned
}

func StatusMsg(status int) string {
	return strings.ReplaceAll(strings.ToLower(http.StatusText(status)), " ", "_")
}
