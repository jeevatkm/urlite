package middleware

import (
	"net/http"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
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
			c.Env["Config"] = a.Config
			c.Env["AppVersion"] = context.VERSION
			c.Env["DomainCount"] = len(a.Domains)
			//c.Env["AllLinkCount"] = a.AllLinkCount()

			h.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func Session(a *context.App) func(*web.C, http.Handler) http.Handler {
	return func(c *web.C, h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// No session for Path /api/* since it's token based
			if !isApiRoute(r) {
				session, err := a.Store.Get(r, "urlite-session")

				if err == nil { // No error we got the session
					c.Env["IsNewSession"] = session.IsNew
					c.Env["Session"] = session

					if session.IsNew {
						log.Debugf("New session created: %v", session.IsNew)
					}
				} else {
					log.Errorf("Could not be decoded, it's protected resources", err)
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
			if isApiRoute(r) {
				c.Env["ReqMode"] = "API"
			} else {
				c.Env["ReqMode"] = "WEB"
				session := c.Env["Session"].(*sessions.Session)

				if userId, exists := session.Values["User"]; exists {
					bsonId := userId.(bson.ObjectId)
					log.Debugf("User info exists in the session: %v", bsonId.Hex())

					user, err := model.GetUserById(a.DB(), bsonId)
					if err != nil {
						log.Warnf("Authentication error: %v", err)
						delete(c.Env, "User")
						c.Env["IsLoggedIn"] = false
					} else {
						log.Debugf("Assigning user into request context -> c.Env")
						c.Env["User"] = user
						c.Env["IsLoggedIn"] = true

						if user.IsAdmin() {
							c.Env["UserCount"] = model.GetActiveUserCount(a.DB())
						}
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
				http.Redirect(w, r, "/login?rt=/admin"+r.URL.Path, http.StatusFound)
			}
		}

		return http.HandlerFunc(fn)
	}
}

func MediaTypeCheck(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" || r.Method == "PUT" {
			ct := strings.ToLower(strings.TrimSpace(r.Header.Get("Content-Type")))
			if ct != "application/json" {
				w.WriteHeader(http.StatusUnsupportedMediaType)
				return
			}
		}

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func ApiAuth(a *context.App) func(*web.C, http.Handler) http.Handler {
	return func(c *web.C, h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/stats" {
				auth := r.Header.Get("Authorization")
				if !strings.HasPrefix(auth, "Bearer ") {
					log.Error("User bearer is not provided")
					authRequired(w)
					return
				}

				// Validating Bearer
				bearer := auth[7:]
				u, err := model.GetUserByBearer(a.DB(), &bearer)
				if err != nil {
					log.Errorf("User is not exists for bearer '%v', error: %v", bearer, err)
					authRequired(w)
					return
				}

				// Update last api access
				u.ApiIPAddress = r.RemoteAddr
				u.LastApiAccessed = time.Now()
				c.Env["User"] = u
				go func(db *mgo.Database, u *model.User) {
					err := model.UpdateUserLastApiAccess(db, u)
					if err != nil {
						log.Errorf("Unable to update last api access for user '%v'", u.ID.Hex())
					} else {
						log.Debugf("Last api access update completed for '%v'", u.ID.Hex())
					}
				}(a.DB(), u)
			}

			h.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func isApiRoute(r *http.Request) bool {
	return strings.HasPrefix(r.URL.Path, "/api")
}

func authRequired(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Bearer realm="urlite-api"`)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(`{"id":"unauthorized", "message": "Provide valid credential"}`))
}
