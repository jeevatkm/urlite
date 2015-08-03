package middleware

import (
	"net/http"
	"strings"

	"gopkg.in/mgo.v2/bson"

	log "github.com/Sirupsen/logrus"

	"github.com/gorilla/sessions"
	"github.com/jeevatkm/urlite/context"
	"github.com/jeevatkm/urlite/model"
	"github.com/zenazn/goji/web"
)

func AppInfo(a *context.App) func(*web.C, http.Handler) http.Handler {
	return func(c *web.C, h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			c.Env["AppName"] = a.Config.AppName
			c.Env["AppVersion"] = context.VERSION
			c.Env["OrgName"] = a.Config.Owner.Org
			h.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func Session(a *context.App) func(*web.C, http.Handler) http.Handler {
	return func(c *web.C, h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// No session for Path /api/* since it's token based
			//if !strings.HasPrefix(r.URL.Path, "/api") {
			if !isApiRoute(r) {
				session, err := a.Store.Get(r, "urlite-session")

				if err == nil { // No error we got the session
					c.Env["IsNewSession"] = session.IsNew
					c.Env["Session"] = session

					if session.IsNew {
						log.Debugf("New session created: %v", session.IsNew)
					}
				} else {
					log.Errorf("Could not be decoded", err)
				}
			}

			h.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func Auth(a *context.App) func(*web.C, http.Handler) http.Handler {
	return func(c *web.C, h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if !isApiRoute(r) {
				session := c.Env["Session"].(*sessions.Session)

				if userId, exists := session.Values["User"]; exists {
					bsonId := userId.(bson.ObjectId)
					log.Debugf("User info exists in the session: %v", bsonId.Hex())

					user, err := model.GetUserById(a.DBDefault(), bsonId)
					if err != nil {
						log.Warnf("Authentication error: %v", err)
						delete(c.Env, "User")
						c.Env["IsLoggedIn"] = false
					} else {
						log.Debugf("Assigning user into request context -> c.Env")
						c.Env["User"] = user
						c.Env["IsLoggedIn"] = true
					}
				}
			}

			h.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func AdminAuth(a *context.App) func(*web.C, http.Handler) http.Handler {
	return func(c *web.C, h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if u, exists := c.Env["User"]; exists {
				if u.(*model.User).IsAdmin() {
					h.ServeHTTP(w, r)
				} else {
					log.Warnf("User[%v] doesn't have admin rights, denying access to /admin/*", u.(*model.User).Email)
					http.Error(w, "The supplied credentials were valid but still were not enough to grant access", http.StatusForbidden)
				}
			} else {
				log.Error("User not loggedin, sending to login page")
				http.Redirect(w, r, "/login", http.StatusFound)
			}
		}

		return http.HandlerFunc(fn)
	}
}

func ApiAuth(a *context.App) func(*web.C, http.Handler) http.Handler {
	return func(c *web.C, h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			log.Infof("I'm in API route: %v", r.URL.Path)
			h.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func isApiRoute(r *http.Request) bool {
	return strings.HasPrefix(r.URL.Path, "/api")
}
