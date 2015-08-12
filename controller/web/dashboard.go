package web

import (
	"net/http"

	"github.com/jeevatkm/urlite/context"
	"github.com/zenazn/goji/web"

	. "github.com/jeevatkm/urlite/controller"
)

func Dashboard(a *context.App, c web.C, r *http.Request) *Response {
	u := GetUser(c)

	if u.IsAdmin() {
		AddData(c, Data{
			"IsDashboard": true,
		})
	} else {
		AddData(c, Data{
			"IsUserDashboard": true,
		})
	}

	content, err := a.Parse("dashboard", c.Env)
	code := CheckError(err)

	AddData(c, Data{
		"Title":   "Dashboard | urlite",
		"Content": ToHTML(content),
	})

	body, err := a.ParseF(c.Env)
	code = CheckError(err)

	return HTMLc(body, code)
}
