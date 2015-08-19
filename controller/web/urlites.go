package web

import (
	"net/http"
	"strings"

	"gopkg.in/mgo.v2/bson"

	"github.com/jeevatkm/urlite/context"
	"github.com/jeevatkm/urlite/model"
	"github.com/jeevatkm/urlite/util"
	"github.com/zenazn/goji/web"

	log "github.com/Sirupsen/logrus"
	. "github.com/jeevatkm/urlite/controller"
)

func Urlites(a *context.App, c web.C, r *http.Request) *Response {
	u := GetUser(c)

	if u.IsAdmin() {
		AddData(c, Data{
			"IsUrlites": true,
			"Domains":   a.Domains,
		})
	} else {
		AddData(c, Data{
			"IsUserUrlites": true,
			"Domains":       a.Domains,
		})
	}

	content, err := a.Parse("urlites", c.Env)
	code := CheckError(err)

	AddData(c, Data{
		"Title":   "Your Urlites | urlite",
		"Content": ToHTML(content),
	})

	body, err := a.ParseF(c.Env)
	code = CheckError(err)

	return HTMLc(body, code)
}

func UrlitesData(a *context.App, c web.C, r *http.Request) *Response {
	dName := strings.TrimSpace(r.FormValue("dName"))
	log.Debugf("Domain Name: %v", dName)

	u := GetUser(c)
	if !util.Contains(u.Domains, dName) && !u.IsAdmin() {
		msg := "You do not have access to given domain"
		log.Errorf("%v: %v", u.Email, msg)
		return ErrForbidden(msg)
	}

	d, err := a.GetDomainDetail(dName)
	if err != nil {
		log.Errorf("Error: %v", err)
		return ErrBadRequest(err.Error())
	}

	page, db := ParsePagination(r), a.DB(&c)
	q := bson.M{}
	if u.IsAdmin() {
		if dName == "*" {
			q = bson.M{}
		} else {
			q = bson.M{"domain": dName}
		}
	} else {
		q = bson.M{"cb": u.ID.Hex()}
	}

	pageResult, _ := model.GetUrliteByPage(db, d.TrackName, q, page)
	body, err := MarshalJSON(pageResult)
	if err != nil {
		log.Errorf("JSON Marshal error: %q", err)
		return ErrInternalServer("Unable to get urlites")
	}

	return JSON(body)
}
