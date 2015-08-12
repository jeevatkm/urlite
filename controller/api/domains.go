package api

import (
	"net/http"
	"strings"

	"github.com/jeevatkm/urlite/context"
	"github.com/jeevatkm/urlite/model"
	"github.com/jeevatkm/urlite/util"
	"github.com/zenazn/goji/web"

	log "github.com/Sirupsen/logrus"
	. "github.com/jeevatkm/urlite/controller"
)

func Domains(a *context.App, c web.C, r *http.Request) *Response {
	u := GetUser(c)
	dName := strings.TrimSpace(c.URLParams["name"])

	if len(dName) > 0 {
		return handleDomainInfo(a, u, dName)
	}

	return handleDomains(a, u)
}

func handleDomains(a *context.App, u *model.User) *Response {
	var domains []string
	if u.IsAdmin() {
		for k, _ := range a.LinkState {
			domains = append(domains, k)
		}
	} else {
		domains = u.Domains
	}

	result, err := MarshalJSON(Data{"domains": domains})
	if err != nil {
		log.Errorf("JSON Marshal error: %q", err)
		return errInternalServer("Unable to get user associated domains")
	}

	return JSON(result)
}

func handleDomainInfo(a *context.App, u *model.User, dName string) *Response {
	if util.Contains(u.Domains, dName) || u.IsAdmin() {
		result, err := MarshalJSON(a.Domains[dName])
		if err != nil {
			log.Errorf("JSON Marshal error: %q", err)
			return errInternalServer("Unable to get domain information")
		}

		return JSON(result)
	}

	return errForbidden("You do not have access to given domain")
}
