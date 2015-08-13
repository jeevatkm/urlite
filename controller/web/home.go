package web

import (
	"net/http"

	"github.com/jeevatkm/urlite/context"
	"github.com/jeevatkm/urlite/controller/api"
	"github.com/jeevatkm/urlite/model"
	"github.com/zenazn/goji/web"

	. "github.com/jeevatkm/urlite/controller"
)

func Home(a *context.App, c web.C, r *http.Request) *Response {
	content, err := a.Parse("home", c.Env)
	code := CheckError(err)

	AddData(c, Data{
		"IsHome":  true,
		"Title":   "Home | urlite",
		"Content": ToHTML(content),
	})

	body, err := a.ParseF(c.Env)
	code = CheckError(err)

	return HTMLc(body, code)
}

func Urlite(a *context.App, c web.C, r *http.Request) *Response {
	shortReq := &model.ShortenRequest{LongUrl: r.FormValue("longUrl"),
		Domain:     r.FormValue("selectedDomain"),
		CustomName: r.FormValue("customLinkName")}
	c.Env["shortReq"] = shortReq

	return api.HandleUrlite(a, c, r)
}
