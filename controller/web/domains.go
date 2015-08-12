package web

import (
	"fmt"
	"net/http"

	"github.com/jeevatkm/urlite/context"
	"github.com/jeevatkm/urlite/model"
	"github.com/jeevatkm/urlite/util/random"
	"github.com/zenazn/goji/web"

	log "github.com/Sirupsen/logrus"
	. "github.com/jeevatkm/urlite/controller"
)

func Domains(a *context.App, c web.C, r *http.Request) *Response {
	domains, _ := model.GetAllDomain(a.DB())
	AddData(c, Data{
		"IsDomains": true,
		"Domains":   domains,
	})

	content, err := a.Parse("domains", c.Env)
	code := CheckError(err)

	AddData(c, Data{
		"Title":   "Domains | urlite",
		"Content": ToHTML(content),
	})

	body, err := a.ParseF(c.Env)
	code = CheckError(err)

	return HTMLc(body, code)
}

func DomainsPost(a *context.App, c web.C, r *http.Request) *Response {
	dName := r.FormValue("dName")
	ed, _ := model.GetDomain(a.DB(), &dName)
	if ed != nil {
		m := fmt.Sprintf("Given domain '%v' already exists", dName)
		log.Error(m)
		SetErrorAlert(c, m)
		return Domains(a, c, r)
	}

	dScheme, dColl, dStatsColl := r.FormValue("dScheme"), r.FormValue("dColl"), r.FormValue("dStatsColl")
	domain := &model.Domain{Name: dName,
		Scheme:        dScheme,
		Salt:          random.DomainSalt(),
		CollName:      dColl,
		StatsCollName: dStatsColl,
		CreatedBy:     GetUser(c).ID.Hex()}

	err := model.CreateDomain(a.DB(), domain)
	if err != nil {
		log.Errorf("Unable to create domain: %q", err)
		SetErrorAlert(c, "Unable to create domain: "+dName)
		return Domains(a, c, r)
	}

	d, _ := model.GetDomain(a.DB(), &dName)
	a.AddDomain(d)
	c.Env["DomainCount"] = len(a.Domains)
	msg := "Successfully added domain: " + dName
	SetSuccessAlert(c, msg)
	log.Info(msg)
	return Domains(a, c, r)
}
