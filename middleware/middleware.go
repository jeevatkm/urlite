package middleware

import (
	"net/http"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/sessions"
	"github.com/jeevatkm/urlite/context"
	"github.com/jeevatkm/urlite/model"
	"github.com/zenazn/goji/web"

	log "github.com/Sirupsen/logrus"
)

func AppInfo(a *context.App) func(*web.C, http.Handler) http.Handler {
	return func(c *web.C, h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			c.Env["Config"] = a.Config
			c.Env["AppVersion"] = context.VERSION
			c.Env["DomainCount"] = len(a.Domains)

			h.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func Database(a *context.App) func(*web.C, http.Handler) http.Handler {
	return func(c *web.C, h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			dbs := a.DBSession.Clone()
			defer dbs.Close()
			c.Env["DB"] = dbs.DB(a.Config.DB.DBName)

			h.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func Session(a *context.App) func(*web.C, http.Handler) http.Handler {
	return func(c *web.C, h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// No session for Path /api/* since it's goes with bearer token
			if !isRoute(r, a.NoSessionRoute...) {
				session, err := a.Store.Get(r, "urlite-session")
				if err == nil { // if err is nil, it means we got the session
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
			if isRoute(r, "/api") {
				c.Env["ReqMode"] = "API"
			} else if !isRoute(r, a.Config.Security.PublicPath...) {
				c.Env["ReqMode"] = "WEB"
				session := c.Env["Session"].(*sessions.Session)

				if userId, exists := session.Values["User"]; exists {
					bsonId := userId.(bson.ObjectId)
					db := a.DB(c)
					log.Debugf("User info exists in the session: %v", bsonId.Hex())

					user, err := model.GetUserById(db, bsonId)
					if err != nil {
						log.Warnf("Authentication error: %v", err)
						delete(c.Env, "User")
						c.Env["IsLoggedIn"] = false
					} else {
						log.Debugf("Assigning user into request context -> c.Env")
						c.Env["User"] = user
						c.Env["IsLoggedIn"] = true

						if user.IsAdmin() {
							c.Env["UserCount"] = model.GetActiveUserCount(db)
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
			if !strings.HasPrefix(r.URL.Path, "/stats") {
				auth := r.Header.Get("Authorization")
				if !strings.HasPrefix(auth, "Bearer ") {
					log.Error("User bearer is not provided")
					authRequired(w)
					return
				}

				// Validating Bearer
				bearer := auth[7:]
				db := a.DB(c)
				u, err := model.GetUserByBearer(db, &bearer)
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
				}(db, u)
			}

			h.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func isRoute(r *http.Request, rts ...string) bool {
	for _, rt := range rts {
		if strings.HasPrefix(r.URL.Path, rt) {
			return true
		}
	}
	return false
}

func authRequired(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Bearer realm="urlite-api"`)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(`{"id":"unauthorized", "message": "Provide valid credential"}`))
}
