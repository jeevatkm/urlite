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

func Dashboard(a *context.App, c web.C, r *http.Request) *Response {
	u := GetUser(c)

	if u.IsAdmin() {
		AddData(c, Data{
			"IsDashboard": true,
			"Domains":     a.Domains,
		})
	} else {
		AddData(c, Data{
			"IsUserDashboard": true,
			"Domains":         a.Domains,
		})
	}

	content, err := a.Parse("dashboard", c.Env)
	code := CheckError(err)

	AddData(c, Data{
		"Title":   "Dashboard | urlite",
		"Content": ToHTML(content),
	})

	body, err := a.ParseF(c.Env)
	code = CheckError(err)

	return HTMLc(body, code)
}

func DashboardData(a *context.App, c web.C, r *http.Request) *Response {
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

	pageResult, _ := model.GetUrliteByPage(db, d.CollName, q, page)
	body, err := MarshalJSON(pageResult)
	if err != nil {
		log.Errorf("JSON Marshal error: %q", err)
		return ErrInternalServer("Unable to get urlites")
	}

	return JSON(body)
}
