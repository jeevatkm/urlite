package api

import (
	"fmt"
	"net/http"

	"github.com/jeevatkm/urlite/context"
	"github.com/jeevatkm/urlite/model"
	"github.com/zenazn/goji/web"

	log "github.com/Sirupsen/logrus"
	. "github.com/jeevatkm/urlite/controller"
)

func Shorten(a *context.App, c web.C, r *http.Request) (*Response, error) {
	shortReq := &model.ShortenRequest{}
	err := DecodeJSON(r, &shortReq)
	if err != nil {
		log.Errorf("Unmarshal error: %q", err)

		errJson := cApiError("bad_request", "The request could not be understood by the urlite api due to bad syntax")
		return cResponse(errJson, http.StatusBadRequest), nil
	}

	domain, err := a.GetDomainDetail(shortReq.Domain)
	if err != nil {
		log.Errorf("Invalid domain: %v", err)

		errJson := cApiError("error", "Invalid domain")
		return cResponse(errJson, http.StatusBadRequest), nil
	}

	linkNum := a.GetDomainLinkNum(domain.Name)
	urliteId, err := a.GetUrliteID(shortReq.Domain, linkNum)
	if err != nil {
		log.Errorf("Unable to generate hashid for given number: %q", err)

		errJson := cApiError("error", "Unable to generate urlite id")
		return cResponse(errJson, http.StatusInternalServerError), nil
	}

	urlite := fmt.Sprintf("%v://%v/%v", domain.Scheme, domain.Name, urliteId)
	sr := &model.ShortenResponse{Urlite: urlite}
	result, err := MarshalJSON(sr)
	if err != nil {
		log.Errorf("JSON Marshal error: %q", err)

		errJson := cApiError("error", "Unable to generate urlite")
		return cResponse(errJson, http.StatusInternalServerError), nil
	}

	return cResponseH(result, http.StatusCreated, map[string]string{"Location": urlite}), nil
}
