package controller

import (
	"net/http"

	"github.com/jeevatkm/urlite/gateway/context"
	"github.com/jeevatkm/urlite/model"
	"github.com/zenazn/goji/web"

	log "github.com/Sirupsen/logrus"
)

type Handle struct {
	*context.Context
	H func(*context.Context, web.C, http.ResponseWriter, *http.Request)
}

func (h Handle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.ServeHTTPC(web.C{}, w, r)
}

func (h Handle) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	h.H(h.Context, c, w, r)
}

func Home(ctx *context.Context, c web.C, w http.ResponseWriter, r *http.Request) {
	domain := model.Domain{}
	log.Infof("%q", domain)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("gateway home is here"))
}

func Urlite(ctx *context.Context, c web.C, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(c.URLParams["urlite"]))
}
