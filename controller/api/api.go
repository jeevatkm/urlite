package api

import (
	. "github.com/jeevatkm/urlite/controller"
)

type ApiError struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

func cApiError(id, msg string) string {
	ae := ApiError{Id: id, Message: msg}
	errJson, _ := MarshalJSON(ae)

	return errJson
}

func cResponse(body string, code int) *Response {
	return &Response{ContentType: JSON_CONTENT, Body: body, Code: code}
}

func cResponseH(body string, code int, h map[string]string) *Response {
	return &Response{ContentType: JSON_CONTENT, Body: body, Code: code, Headers: h}
}
