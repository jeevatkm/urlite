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
	domains, _ := model.GetAllDomain(a.DB(&c))
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
	dTrack := strings.TrimSpace(r.FormValue("dTrack"))
	log.Debugf("dTrack: %v", dTrack)
	if len(dTrack) > 0 {
		res := a.CheckDomainTrackName(dTrack)
		if res {
			log.Debugf("Collection name [%v] exists, will be suggestion one", dTrack)
			// Hoping not to reach more than 20 in numbers, might need a revisit
			for i := 1; i <= 20; i++ {
				nDTrack := fmt.Sprintf("%v%d", dTrack, i)
				res = a.CheckDomainTrackName(nDTrack)
				if !res {
					data := Data{
						"id":             "validation",
						"message":        "Collection already exists",
						"suggested_name": nDTrack,
					}

					response = PrepareJSONc(data, http.StatusBadRequest, GENERIC_ERROR_MSG)
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

	dScheme, dTrack := r.FormValue("dScheme"), r.FormValue("dTrack")
	domain := &model.Domain{Name: dName,
		Scheme:    dScheme,
		Salt:      random.DomainSalt(),
		TrackName: dTrack,
		CreatedBy: GetUser(c).ID.Hex()}
	db := a.DB(&c)

	err := model.CreateDomain(db, domain)
	if err != nil {
		log.Errorf("Unable to create domain: %q", err)
		return ErrInternalServer("Unable to add domain due to server issue")
	}

	d, _ := model.GetDomain(db, &dName)
	a.AddDomain(d)
	msg := "Successfully added domain: " + dName
	data := Data{
		"id":      ALERT_SUCCESS,
		"message": msg,
	}
	log.Info(msg)

	return PrepareJSON(data, GENERIC_ERROR_MSG)
}
