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
	liteReq := &model.UrliteRequest{}
	if err := DecodeJSON(r, &liteReq); err != nil {
		log.Errorf("Unmarshal error: %q", err)
		return ErrBadRequest("The request could not be understood by the urlite api due to bad syntax")
	}
	c.Env["liteReq"] = liteReq
	return HandleUrlite(a, c, r)
}

func HandleUrlite(a *context.App, c web.C, r *http.Request) *Response {
	liteReq := c.Env["liteReq"].(*model.UrliteRequest)
	if !liteReq.IsValid() {
		msg := "Either Long URL or Domain is not provided"
		log.Error(msg)
		return ErrBadRequest(msg)
	}

	u := GetUser(c)
	if !util.Contains(u.Domains, liteReq.Domain) && !u.IsAdmin() {
		msg := "You do not have access to given domain"
		log.Errorf("%v: %v", u.Email, msg)
		return ErrForbidden(msg)
	}

	domain, err := a.GetDomainDetail(liteReq.Domain)
	if err != nil {
		log.Errorf("Invalid domain: %v", err)
		return ErrValidation("Invalid domain")
	}

	urlite, db := "", a.DB(&c)
	urliteId := strings.TrimSpace(liteReq.CustomName)

	if len(urliteId) > 0 { // Custom name mode
		log.Debug("Custom name mode")

		exists, _ := model.GetUrlite(db, domain.CollName, &urliteId)
		if exists != nil {
			log.Errorf("Given custom name is unavailable: %v", err)
			return ErrConflict("Given custom name [" + urliteId + "] is unavailable")
		}
		a.IncDomainCustomLink(domain.Name)
	} else { // Hash generate mode
		log.Debug("Hash generate mode")

		linkNum := a.GetDomainLinkNum(domain.Name)
		turliteId, err := a.GetUrliteID(domain.Name, linkNum)
		if err != nil {
			log.Errorf("Unable to generate hashid for number[%d]: %q", linkNum, err)
			return ErrInternalServer("Unable to generate urlite")
		}

		urliteId = turliteId
	}

	urlite = domain.ComposeUrlite(&urliteId)
	ul := &model.Urlite{ID: urliteId,
		Urlite:      urlite,
		LongUrl:     strings.TrimSpace(liteReq.LongUrl),
		Domain:      domain.Name,
		CreatedBy:   u.ID.Hex(),
		CreatedTime: time.Now()}
	err = model.CreateUrlite(db, domain.CollName, ul)
	if err != nil {
		log.Errorf("Unable to insert new urlite into db: %q", err)
		return ErrInternalServer("Unable to generate urlite")
	}

	sr := &model.UrliteResponse{Urlite: urlite}
	result, err := MarshalJSON(sr)
	if err != nil {
		log.Errorf("JSON Marshal error: %q", err)
		return ErrInternalServer("Unable to generate urlite")
	}

	return JSONch(result, http.StatusCreated, map[string]string{"Location": urlite})
}
