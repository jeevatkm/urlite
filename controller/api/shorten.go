package api

import (
	"net/http"

	"github.com/jeevatkm/urlite/context"
	"github.com/jeevatkm/urlite/model"
	"github.com/zenazn/goji/web"

	log "github.com/Sirupsen/logrus"
	. "github.com/jeevatkm/urlite/controller"
)

type ApiError struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

func Shorten(a *context.App, c web.C, r *http.Request) (*Response, error) {
	code := http.StatusOK
	var resBody interface{}

	shortReq := &model.ShortenRequest{}
	err := DecodeJSON(r, &shortReq)

	if err != nil {
		log.Errorf("Unmarshal error: %q", err)
		code = http.StatusBadRequest
		resBody = &ApiError{Id: "bad_request", Message: "The request could not be understood by the urlite api due to bad syntax."}
		goto E
	}

	log.Infof("Received request %q", shortReq)
	log.Infof("Domain details %q", a.Domains["sample.com"])

	resBody = &model.ShortenResponse{Urlite: "http://test.com/sdjgfdh"}

	// sreq := &model.ShortenRequest{
	// 	LongUrl:  "http://test.com/sdjgfdhfgdhgfhdgfdgf",
	// 	Domain:   "sample.com",
	// 	Secure:   true,
	// 	Password: "pass"}

	// res, err := json.Marshal(res)
	// if err != nil {
	// 	log.Errorf("JSON marshal error: %s", err)
	// }
E:
	jres, err := MarshalJSON(resBody)
	// if err != nil {
	// 	log.Errorf("Error occurred", err)
	// 	resBody = &ApiError{Id: "error", Message: "Unable to generated the urlite"}
	// 	jres, err := MarshalJSON(resBody)
	// }

	return &Response{ContentType: JSON_CONTENT, Body: jres, Code: code}, nil
}
