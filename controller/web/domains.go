package web

import (
	"fmt"
	"net/http"
	"strings"

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

func DomainsValidate(a *context.App, c web.C, r *http.Request) *Response {
	var response *Response

	// Domain name validation
	dName := strings.TrimSpace(r.FormValue("dName"))
	log.Debugf("dName: %v", dName)
	if len(dName) > 0 {
		if _, ok := a.Domains[dName]; ok {
			response = ErrBadRequest("Domain already exists")
		} else {
			response = JSON("{}")
		}
	}

	// Collection name validation
	dColl := strings.TrimSpace(r.FormValue("dColl"))
	log.Debugf("dColl: %v", dColl)
	if len(dColl) > 0 {
		res := a.CheckDomainCollName(dColl)
		if res {
			log.Debugf("Collection name [%v] exists, will be suggestion one", dColl)
			// Hoping not to reach more than 20 in numbers, might need a revisit
			for i := 1; i <= 20; i++ {
				nDColl := fmt.Sprintf("%v%d", dColl, i)
				res = a.CheckDomainCollName(nDColl)
				if !res {
					body := `{"id":"bad_request","message": "Collection already exists","suggested_name": "` + nDColl + `"}`
					response = JSONc(body, http.StatusBadRequest)
					break
				}
			}
		} else {
			response = JSON("{}")
		}
	}

	return response
}

func DomainsPost(a *context.App, c web.C, r *http.Request) *Response {
	dName := strings.TrimSpace(r.FormValue("dName"))
	if _, ok := a.Domains[dName]; ok {
		m := fmt.Sprintf("Given domain '%v' already exists", dName)
		log.Error(m)
		return ErrBadRequest(m)
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
		return ErrInternalServer("Unable to add domain due to server issue")
	}

	d, _ := model.GetDomain(a.DB(), &dName)
	a.AddDomain(d)
	msg := "Successfully added domain: " + dName
	log.Info(msg)

	body := `{"id":"success","message": "` + msg + `"}`
	return JSON(body)
}
