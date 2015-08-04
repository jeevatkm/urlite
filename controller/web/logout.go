package web

import (
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/sessions"
	"github.com/jeevatkm/urlite/context"
	"github.com/zenazn/goji/web"

	log "github.com/Sirupsen/logrus"
	ctr "github.com/jeevatkm/urlite/controller"
)

func Logout(a *context.App, c web.C, r *http.Request) (*ctr.Response, error) {
	session, exists := c.Env["Session"]

	if exists {
		log.Debugf("Session exists, is user loggedin?")
		s := session.(*sessions.Session)

		if uid, loggedIn := s.Values["User"]; loggedIn {
			log.Debugf("Yes user loggedin and Hex ID is %v", uid.(bson.ObjectId).Hex())

			// Removing user from session
			delete(s.Values, "User")
		} else {
			log.Debug("User not loggedin.")
		}
	} else {
		log.Error("User session is not exists, incorrect invocation")
	}

	return &ctr.Response{Redirect: "/", Code: http.StatusSeeOther}, nil
}
