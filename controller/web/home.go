package web

import (
	"net/http"

	"github.com/jeevatkm/urlite/context"
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
