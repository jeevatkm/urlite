package web

import (
	"net/http"

	"github.com/jeevatkm/urlite/context"
	"github.com/jeevatkm/urlite/model"
	"github.com/zenazn/goji/web"

	log "github.com/Sirupsen/logrus"
	. "github.com/jeevatkm/urlite/controller"
)

func Login(a *context.App, c web.C, r *http.Request) (*Response, error) {
	session := GetSession(c)

	content, err := a.Parse("login", Data{
		"Flash": session.Flashes("authmsg"),
	})
	code := CheckError(err)

	AddData(c, Data{
		"IsLogin": true,
		"Title":   "Login | urlite",
		"Content": ToHTML(content),
	})

	body, err := a.ParseF(c.Env)
	code = CheckError(err)

	return &Response{ContentType: HTML_CONTENT, Body: body, Code: code}, err
}

func LoginPost(a *context.App, c web.C, r *http.Request) (*Response, error) {
	email, password := r.FormValue("email"), r.FormValue("password")
	user, err := model.AuthenticateUser(a.DB(), email, password)

	session := GetSession(c)
	if err != nil {
		session.AddFlash("Incorrect Email/Password", "authmsg")
		return Login(a, c, r)
	}

	session.Values["User"] = user.ID

	if user.IsAdmin() {
		log.Debugf("Admin logged in '%s'", user.Email)
		return &Response{Redirect: "/admin/dashboard", Code: http.StatusFound}, nil
	}

	return &Response{Redirect: "/", Code: http.StatusSeeOther}, nil
}
