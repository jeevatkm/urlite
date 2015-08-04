package web

import (
	"net/http"

	"github.com/jeevatkm/urlite/context"
	"github.com/zenazn/goji/web"

	. "github.com/jeevatkm/urlite/controller"
)

func Dashboard(a *context.App, c web.C, r *http.Request) (*Response, error) {
	content, err := a.Parse("dashboard", c)
	code := CheckError(err)

	AddData(c, Data{
		"IsDashboard": true,
		"Title":       "Dashboard - urlite",
		"Content":     ToHTML(content),
	})

	body, err := a.ParseF(c.Env)
	code = CheckError(err)

	return &Response{ContentType: HTML_CONTENT, Body: body, Code: code}, err
}
