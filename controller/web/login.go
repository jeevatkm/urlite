package web

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/jeevatkm/urlite/context"
	"github.com/jeevatkm/urlite/model"
	"github.com/zenazn/goji/web"

	log "github.com/Sirupsen/logrus"
	ctr "github.com/jeevatkm/urlite/controller"
	u "github.com/jeevatkm/urlite/util"
)

func Login(a *context.App, c web.C, r *http.Request) (*ctr.Response, error) {
	session := c.Env["Session"].(*sessions.Session)

	content, err := a.Parse("login", ctr.Data{
		"Flash": session.Flashes("authmsg"),
	})
	u.CheckError(err)

	u.AddData(c, ctr.Data{
		"IsLogin": true,
		"Title":   "Login - urlite",
		"Content": u.ToHTML(content),
	})

	body, err := a.Parse("layout/base", c.Env)
	u.CheckError(err)

	return &ctr.Response{ContentType: ctr.HTML_CONTENT, Body: body, Code: http.StatusOK}, err
}

func LoginPost(a *context.App, c web.C, r *http.Request) (*ctr.Response, error) {
	email, password := r.FormValue("email"), r.FormValue("password")
	user, err := model.AuthenticateUser(a.DBDefault(), email, password)

	session := c.Env["Session"].(*sessions.Session)
	if err != nil {
		session.AddFlash("Incorrect Email/Password", "authmsg")
		return Login(a, c, r)
	}

	session.Values["User"] = user.ID

	if user.IsAdmin() {
		log.Debugf("Admin logged in '%s'", user.Email)
		return &ctr.Response{Redirect: "/admin/", Code: http.StatusFound}, nil
	}

	return &ctr.Response{Redirect: "/", Code: http.StatusSeeOther}, nil
}
