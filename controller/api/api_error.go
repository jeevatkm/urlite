package api

import (
	"net/http"

	. "github.com/jeevatkm/urlite/controller"
)

type ApiError struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

func cApiError(id, msg string) string {
	errJson, _ := MarshalJSON(ApiError{Id: id, Message: msg})
	return errJson
}

// Error helper Method
func errBadRequest(m string) *Response {
	return cResponse(cApiError("bad_request", m), http.StatusBadRequest)
}

func errInternalServer(m string) *Response {
	return cResponse(cApiError("error", m), http.StatusInternalServerError)
}

func errConflict(m string) *Response {
	return cResponse(cApiError("already_exists", m), http.StatusConflict)
}

func errForbidden(m string) *Response {
	return cResponse(cApiError("forbidden", m), http.StatusForbidden)
}

// Functional Errors
func errUnmarshal() *Response {
	return errBadRequest("The request could not be understood by the urlite api due to bad syntax")
}

func errInvalidDomain() *Response {
	return cResponse(cApiError("bad_request", "Invalid domain"), http.StatusBadRequest)
}

func errGenerateUrlite() *Response {
	return errInternalServer("Unable to generate urlite")
}
