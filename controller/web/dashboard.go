package web

import (
	"net/http"

	"github.com/jeevatkm/urlite/context"
	"github.com/zenazn/goji/web"

	ctr "github.com/jeevatkm/urlite/controller"
	u "github.com/jeevatkm/urlite/util"
)

func Dashboard(a *context.App, c web.C, r *http.Request) (*ctr.Response, error) {
	content, err := a.Parse("dashboard", c)
	u.CheckError(err)

	u.AddData(c, ctr.Data{
		"IsDashboard": true,
		"Title":       "Dashboard - urlite",
		"Content":     u.ToHTML(content),
	})

	body, err := a.Parse("layout/base", c.Env)
	u.CheckError(err)

	return &ctr.Response{ContentType: ctr.HTML_CONTENT, Body: body, Code: http.StatusOK}, err
}
