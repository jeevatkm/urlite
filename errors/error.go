package errors

import (
	"encoding/json"
	"fmt"
)

type Error struct {
	Id      string                 `json:"id"`
	Message string                 `json:"message"`
	Values  map[string]interface{} `json:"values"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("Error[%s] %s", e.Id, e.Message)
}

func (e *Error) JSON() string {
	j, err := json.Marshal(e)
	if err != nil {
		return "{}"
	}

	return string(j)
}

func New(id, msg string) *Error {
	return &Error{Id: id, Message: msg}
}

func Newv(id string, msg string, values map[string]interface{}) *Error {
	return &Error{Id: id, Message: msg, Values: values}
}
