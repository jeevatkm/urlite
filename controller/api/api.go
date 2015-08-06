package api

import . "github.com/jeevatkm/urlite/controller"

func cResponse(body string, code int) *Response {
	return &Response{ContentType: JSON_CONTENT, Body: body, Code: code}
}

func cResponseH(body string, code int, hdr map[string]string) *Response {
	return &Response{ContentType: JSON_CONTENT, Body: body, Code: code, Headers: hdr}
}
