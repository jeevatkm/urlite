package web

import (
	"net/http"

	"github.com/jeevatkm/urlite/context"
	"github.com/jeevatkm/urlite/model"
	"github.com/zenazn/goji/web"

	. "github.com/jeevatkm/urlite/controller"
)

func Users(a *context.App, c web.C, r *http.Request) (*Response, error) {
	users, _ := model.GetAllUsers(a.DB())
	AddData(c, Data{
		"IsUsers": true,
		"Users":   users,
	})

	content, err := a.Parse("users", c.Env)
	code := CheckError(err)

	AddData(c, Data{
		"Title":   "Users | urlite",
		"Content": ToHTML(content),
	})

	body, err := a.ParseF(c.Env)
	code = CheckError(err)

	return &Response{ContentType: HTML_CONTENT, Body: body, Code: code}, err
}
