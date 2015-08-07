package web

import (
	"net/http"
	"time"

	"gopkg.in/mgo.v2"

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

	// Update last login successful
	user.LoginIPAddress = r.RemoteAddr
	user.LastLoggedIn = time.Now()
	go func(db *mgo.Database, u *model.User) {
		err := model.UpdateUserLastLogin(db, u)
		if err != nil {
			log.Errorf("Unable to update last login for user '%v'", u.ID.Hex())
		} else {
			log.Debugf("Last login update completed for '%v'", u.ID.Hex())
		}
	}(a.DB(), user)

	if user.IsAdmin() {
		log.Debugf("Admin logged in '%s'", user.Email)
		return &Response{Redirect: "/admin/dashboard", Code: http.StatusFound}, nil
	}

	return &Response{Redirect: "/", Code: http.StatusSeeOther}, nil
}
