package controller

import (
	"io"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/jeevatkm/urlite/context"
	"github.com/zenazn/goji/web"

	log "github.com/Sirupsen/logrus"
)

const (
	TEXT_CONTENT = "text/plain"
	HTML_CONTENT = "text/html; charset=utf-8"
	JSON_CONTENT = "application/json"
)

type Response struct {
	ContentType string
	Body        string
	Code        int
	Redirect    string
}

type Data map[string]interface{}

type Handle struct {
	*context.App
	H func(*context.App, web.C, *http.Request) (*Response, error)
}

func (h Handle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.ServeHTTPC(web.C{}, w, r)
}

func (h Handle) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	res, err := h.H(h.App, c, r)
	body, code, contentType := res.Body, res.Code, res.ContentType

	if err != nil {
		log.Infof("HTTP %d: %q", code, err)
		switch code {
		case http.StatusNotFound:
			http.NotFound(w, r)
		// And if we wanted a friendlier error page:
		// err := ah.renderTemplate(w, "http_404.tmpl", nil)
		case http.StatusInternalServerError:
			http.Error(w, http.StatusText(code), code)
		default:
			http.Error(w, http.StatusText(code), code)
		}
	} else {
		if session, exists := c.Env["Session"]; exists {
			log.Debug("Saving sessions...")
			err := session.(*sessions.Session).Save(r, w)
			if err != nil {
				log.Errorf("Can't save session: %v", err)
			}
		}

		switch code {
		case http.StatusOK:
			w.Header().Set("Content-Type", contentType)
			io.WriteString(w, body)
		case http.StatusSeeOther, http.StatusFound:
			http.Redirect(w, r, res.Redirect, code)
		default:
			w.WriteHeader(code)
			io.WriteString(w, body)
		}
	}
}
