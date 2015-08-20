package api

import (
	"net/http"
	"strings"

	"github.com/zenazn/goji/web"

	. "github.com/jeevatkm/urlite/dashboard/context"
	. "github.com/jeevatkm/urlite/dashboard/controller"
)

func Stats(ctx *Context, c web.C, r *http.Request) *Response {
	// Individual domain stats
	dName := strings.TrimSpace(c.URLParams["name"])
	if len(dName) > 0 {
		if domain, ok := ctx.Domains[dName]; ok {
			info := &Domain{Name: domain.Name,
				Total:        domain.LinkCount + domain.CustomLinkCount,
				Urlite:       domain.LinkCount,
				CustomUrlite: domain.CustomLinkCount}
			return PrepareJSON(info, "Unable to generate urlite stats")
		}
		return ErrValidation("Invalid domain")
	}

	// For all domains stats
	stats := &UrliteStats{}
	stats.Domains = []*Domain{}
	var all int64

	for _, v := range ctx.Domains {
		all += v.LinkCount + v.CustomLinkCount
		stats.Domains = append(stats.Domains, &Domain{Name: v.Name,
			Total:        v.LinkCount + v.CustomLinkCount,
			Urlite:       v.LinkCount,
			CustomUrlite: v.CustomLinkCount})
	}

	stats.TotalUrlite = all
	stats.DomainCount = len(ctx.Domains)

	return PrepareJSON(stats, "Unable to generate urlite stats")
}
