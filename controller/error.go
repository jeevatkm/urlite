package controller

import (
	"net/http"
)

type Error struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

func cError(id, msg string) string {
	errJson, _ := MarshalJSON(Error{Id: id, Message: msg})
	return errJson
}

// Error helper Method
func ErrBadRequest(m string) *Response {
	return JSONc(cError("bad_request", m), http.StatusBadRequest)
}

func ErrInternalServer(m string) *Response {
	return JSONc(cError("error", m), http.StatusInternalServerError)
}

func ErrConflict(m string) *Response {
	return JSONc(cError("already_exists", m), http.StatusConflict)
}

func ErrForbidden(m string) *Response {
	return JSONc(cError("forbidden", m), http.StatusForbidden)
}

// Functional Errors
func ErrUnmarshal() *Response {
	return ErrBadRequest("The request could not be understood by the urlite api due to bad syntax")
}

func ErrInvalidDomain() *Response {
	return JSONc(cError("bad_request", "Invalid domain"), http.StatusBadRequest)
}

func ErrGenerateUrlite() *Response {
	return ErrInternalServer("Unable to generate urlite")
}
