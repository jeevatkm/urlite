package api

import (
	"net/http"
	"strings"

	"github.com/jeevatkm/urlite/util"
	"github.com/zenazn/goji/web"

	. "github.com/jeevatkm/urlite/dashboard/context"
	. "github.com/jeevatkm/urlite/dashboard/controller"
)

func Domains(ctx *Context, c web.C, r *http.Request) *Response {
	u := GetUser(c)
	dName := strings.TrimSpace(c.URLParams["name"])

	if len(dName) > 0 {
		if util.Contains(u.Domains, dName) || u.IsAdmin() {
			return PrepareJSON(ctx.Domains[dName], "Unable to get domain information")
		}

		return ErrForbidden("You do not have access to given domain")
	}

	var domains []string
	if u.IsAdmin() {
		for k, _ := range ctx.LinkState {
			domains = append(domains, k)
		}
	} else {
		domains = u.Domains
	}

	return PrepareJSON(Data{"domains": domains}, "Unable to get user associated domains")
}
