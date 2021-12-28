package web

import (
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

func NewResponse(code int, content interface{}) interface{} {
	if code < 300 {
		return Response{strconv.FormatInt(int64(code), 10), content}
	}

	return Error{strconv.FormatInt(int64(code), 10), content}
}

func ResponseUnauthorized() Error {
	return Error{strconv.FormatInt(int64(http.StatusUnauthorized), 10), "no tiene permisos para realizar la petición solicitada"}
}
