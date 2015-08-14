package api

import (
	"net/http"

	"github.com/jeevatkm/urlite/context"
	"github.com/zenazn/goji/web"

	log "github.com/Sirupsen/logrus"
	. "github.com/jeevatkm/urlite/controller"
)

type Domain struct {
	Name         string `json:"name"`
	Total        int64  `json:"total"`
	Urlite       int64  `json:"urlite"`
	CustomUrlite int64  `json:"custom_urlite"`
}

type UrliteStats struct {
	TotalUrlite int64     `json:"total_urlite"`
	DomainCount int       `json:"domain_count"`
	Domains     []*Domain `json:"domains"`
}

func Stats(a *context.App, c web.C, r *http.Request) *Response {
	stats := &UrliteStats{}
	stats.Domains = []*Domain{}
	var all int64

	for _, v := range a.Domains {
		all += v.LinkCount + v.CustomLinkCount
		stats.Domains = append(stats.Domains, &Domain{Name: v.Name,
			Total:        v.LinkCount + v.CustomLinkCount,
			Urlite:       v.LinkCount,
			CustomUrlite: v.CustomLinkCount})
	}

	stats.TotalUrlite = all
	stats.DomainCount = len(a.Domains)
	result, err := MarshalJSON(stats)
	if err != nil {
		log.Errorf("JSON Marshal error: %q", err)
		return errInternalServer("Unable to generate urlite stats")
	}

	return JSON(result)
}
