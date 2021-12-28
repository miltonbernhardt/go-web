package web

import (
	"net/http"
	"strconv"
)

type Response struct {
	Code  string      `json:"code"`
	Data  interface{} `json:"data,omitempty"`
	Error interface{} `json:"error,omitempty"`
}

func NewResponse(code int, content interface{}) Response {
	if code < 300 {
		return Response{strconv.FormatInt(int64(code), 10), content, nil}
	}

	return Response{strconv.FormatInt(int64(code), 10), nil, content}
}

func ResponseUnauthorized() Response {
	return Response{strconv.FormatInt(int64(http.StatusUnauthorized), 10), nil, "no tiene permisos para realizar la petición solicitada"}
}
