package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/jeevatkm/urlite/context"
	"github.com/jeevatkm/urlite/model"
	"github.com/jeevatkm/urlite/util"
	"github.com/zenazn/goji/web"

	log "github.com/Sirupsen/logrus"
	. "github.com/jeevatkm/urlite/controller"
)

func Urlite(a *context.App, c web.C, r *http.Request) *Response {
	shortReq := &model.ShortenRequest{}
	if err := DecodeJSON(r, &shortReq); err != nil {
		log.Errorf("Unmarshal error: %q", err)
		return errUnmarshal()
	}
	c.Env["shortReq"] = shortReq
	return HandleUrlite(a, c, r)
}

func HandleUrlite(a *context.App, c web.C, r *http.Request) *Response {
	shortReq := c.Env["shortReq"].(*model.ShortenRequest)
	if !shortReq.IsValid() {
		msg := "Either Long URL or Domain is not provided"
		log.Error(msg)
		return errBadRequest(msg)
	}

	u := GetUser(c)
	if !util.Contains(u.Domains, shortReq.Domain) && !u.IsAdmin() {
		msg := "You do not have access to given domain"
		log.Errorf("%v: %v", u.Email, msg)
		return errForbidden(msg)
	}

	domain, err := a.GetDomainDetail(shortReq.Domain)
	if err != nil {
		log.Errorf("Invalid domain: %v", err)
		return errInvalidDomain()
	}

	urlite := ""
	urliteId := strings.TrimSpace(shortReq.CustomName)

	if len(urliteId) > 0 { // Custom name mode
		log.Debug("Custom name mode")

		exists, _ := model.GetUrlite(a.DB(), domain.CollName, &urliteId)
		if exists != nil {
			log.Errorf("Given custom name is unavailable: %v", err)
			return errConflict("Given custom name [" + urliteId + "] is unavailable")
		}
		a.IncDomainCustomLink(domain.Name)
	} else { // Hash generate mode
		log.Debug("Hash generate mode")

		linkNum := a.GetDomainLinkNum(domain.Name)
		turliteId, err := a.GetUrliteID(domain.Name, linkNum)
		if err != nil {
			log.Errorf("Unable to generate hashid for number[%d]: %q", linkNum, err)
			return errGenerateUrlite()
		}

		urliteId = turliteId
	}

	urlite = domain.ComposeUrlite(&urliteId)
	ul := &model.Urlite{ID: urliteId,
		Urlite:      urlite,
		LongUrl:     strings.TrimSpace(shortReq.LongUrl),
		CreatedBy:   u.ID.Hex(),
		CreatedTime: time.Now()}
	err = model.CreateUrlite(a.DB(), domain.CollName, ul)
	if err != nil {
		log.Errorf("Unable to insert new urlite into db: %q", err)
		return errGenerateUrlite()
	}

	sr := &model.ShortenResponse{Urlite: urlite}
	result, err := MarshalJSON(sr)
	if err != nil {
		log.Errorf("JSON Marshal error: %q", err)
		return errGenerateUrlite()
	}

	return JSONch(result, http.StatusCreated, map[string]string{"Location": urlite})
}
