package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/jeevatkm/urlite/context"
	"github.com/jeevatkm/urlite/model"
	"github.com/zenazn/goji/web"

	log "github.com/Sirupsen/logrus"
	. "github.com/jeevatkm/urlite/controller"
)

const ()

func Shorten(a *context.App, c web.C, r *http.Request) (*Response, error) {
	shortReq := &model.ShortenRequest{}
	if err := DecodeJSON(r, &shortReq); err != nil {
		log.Errorf("Unmarshal error: %q", err)
		return errUnmarshal(), nil
	}

	domain, err := a.GetDomainDetail(shortReq.Domain)
	if err != nil {
		log.Errorf("Invalid domain: %v", err)
		return errInvalidDomain(), nil
	}

	urlite := ""
	urliteId := strings.TrimSpace(shortReq.CustomName)

	if len(urliteId) > 0 { // Custom name mode
		log.Debug("Custom name mode")

		exists, _ := model.GetUrlite(a.DB(), domain.CollName, &urliteId)
		if exists != nil {
			log.Errorf("Given custom name is unavailable: %v", err)
			return errConflict("Given custom name [" + urliteId + "] is unavailable"), nil
		}
		a.IncDomainCustomLink(domain.Name)
	} else { // Hash generate mode
		log.Debug("Hash generate mode")

		linkNum := a.GetDomainLinkNum(domain.Name)
		turliteId, err := a.GetUrliteID(domain.Name, linkNum)
		if err != nil {
			log.Errorf("Unable to generate hashid for number[%d]: %q", linkNum, err)
			return errGenerateUrlite(), nil
		}

		urliteId = turliteId
	}

	urlite = domain.ComposeUrlite(&urliteId)
	u := GetUser(c)

	ul := &model.Urlite{ID: urliteId,
		Urlite:      urlite,
		LongUrl:     strings.TrimSpace(shortReq.LongUrl),
		CreatedBy:   u.ID.Hex(),
		CreatedTime: time.Now()}
	err = model.CreateUrlite(a.DB(), domain.CollName, ul)
	if err != nil {
		log.Errorf("Unable to insert new urlite into db: %q", err)
		return errGenerateUrlite(), nil
	}

	sr := &model.ShortenResponse{Urlite: urlite}
	result, err := MarshalJSON(sr)
	if err != nil {
		log.Errorf("JSON Marshal error: %q", err)
		return errGenerateUrlite(), nil
	}

	return cResponseH(result, http.StatusCreated, map[string]string{"Location": urlite}), nil
}
