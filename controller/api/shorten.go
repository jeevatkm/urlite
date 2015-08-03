package api

import (
	"encoding/json"
	"net/http"

	"github.com/jeevatkm/urlite/context"
	"github.com/jeevatkm/urlite/model"
	"github.com/zenazn/goji/web"

	log "github.com/Sirupsen/logrus"
	ctr "github.com/jeevatkm/urlite/controller"
)

func Shorten(a *context.App, c web.C, r *http.Request) (*ctr.Response, error) {
	sreq := &model.ShortenRequest{
		LongUrl:  "http://test.com/sdjgfdhfgdhgfhdgfdgf",
		Domain:   "sample.com",
		Secure:   true,
		Password: "pass"}

	res, err := json.Marshal(sreq)
	if err != nil {
		log.Errorf("JSON marshal error: %s", err)
	}

	return &ctr.Response{ContentType: ctr.JSON_CONTENT, Body: string(res), Code: http.StatusOK}, err
}
