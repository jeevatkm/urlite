package api

import (
	"github.com/jeevatkm/urlite/model"
	"github.com/zenazn/goji/web"

	. "github.com/jeevatkm/urlite/controller"
)

func getApiUser(c web.C) *model.ApiUser {
	if au, ok := c.Env["ApiUser"]; ok {
		return au.(*model.ApiUser)
	}
	return nil // Never expected
}

func cResponse(body string, code int) *Response {
	return &Response{ContentType: JSON_CONTENT, Body: body, Code: code}
}

func cResponseH(body string, code int, hdr map[string]string) *Response {
	return &Response{ContentType: JSON_CONTENT, Body: body, Code: code, Headers: hdr}
}
