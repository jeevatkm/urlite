package controller

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/jeevatkm/urlite/dashboard/context"
	"github.com/jeevatkm/urlite/dashboard/errors"
	"github.com/jeevatkm/urlite/model"
	"github.com/zenazn/goji/web"

	log "github.com/Sirupsen/logrus"
)

const (
	TEXT_CONTENT      = "text/plain"
	HTML_CONTENT      = "text/html; charset=utf-8"
	JSON_CONTENT      = "application/json; charset=utf-8"
	ALERT_SUCCESS     = "success"
	ALERT_INFO        = "info"
	ALERT_WARN        = "warning"
	ALERT_ERROR       = "error"
	GENERIC_ERROR_MSG = "Oops! something went wrong"
)

type Response struct {
	ContentType string
	Body        string
	Code        int
	Redirect    string
	Headers     map[string]string
}

type Data map[string]interface{}

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

type Handle struct {
	*context.Context
	H func(*context.Context, web.C, *http.Request) *Response
}

func (h Handle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.ServeHTTPC(web.C{}, w, r)
}

func (h Handle) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	res := h.H(h.Context, c, r)
	body, code, contentType := res.Body, res.Code, res.ContentType

	if session, exists := c.Env["Session"]; exists {
		log.Debug("Saving sessions")
		err := session.(*sessions.Session).Save(r, w)
		if err != nil {
			log.Errorf("Can't save session: %v", err)
			code = http.StatusInternalServerError
		}
	}

	if contentType == JSON_CONTENT {
		w.Header().Set("Content-Type", contentType)

		for k, v := range res.Headers {
			w.Header().Set(k, v)
		}

		w.WriteHeader(code)
		io.WriteString(w, body)
	} else {
		switch code {
		case http.StatusOK: // http.StatusBadRequest
			w.Header().Set("Content-Type", contentType)
			w.WriteHeader(code)
			io.WriteString(w, body)
		case http.StatusSeeOther, http.StatusFound:
			http.Redirect(w, r, res.Redirect, code)
		case http.StatusInternalServerError:
			http.Error(w, http.StatusText(code), code)
		default:
			log.Error("Unable to render output, will do something")
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, GENERIC_ERROR_MSG)
		}
	}
}

func isApiRequest(c web.C) bool {
	return c.Env["ReqMode"] == "API"
}

func GetSession(c web.C) *sessions.Session {
	return c.Env["Session"].(*sessions.Session)
}

func GetUser(c web.C) *model.User {
	return c.Env["User"].(*model.User)
}

func SetErrorAlert(c web.C, m string) {
	c.Env["AlertMsg"] = m
	c.Env["AlertType"] = ALERT_ERROR
}

func SetSuccessAlert(c web.C, m string) {
	c.Env["AlertMsg"] = m
	c.Env["AlertType"] = ALERT_SUCCESS
}

func ParsePagination(r *http.Request) *model.Pagination {
	// "<URL>?sort=last_api_accessed&order=desc&limit=2&offset=0"
	// "<URL>?sort=last_api_accessed&order=desc&limit=2&offset=4"
	log.Debugf("Query params: %v", r.URL.RawQuery)

	limit, _ := strconv.Atoi(r.FormValue("limit"))
	if limit == 0 {
		limit = 10
	}

	offset, _ := strconv.Atoi(r.FormValue("offset"))

	sort := strings.TrimSpace(r.FormValue("sort"))
	if len(sort) == 0 {
		sort = "_id"
	}

	order := strings.TrimSpace(r.FormValue("order"))
	if len(order) == 0 {
		order = "asc"
	}
	if order == "desc" {
		sort = "-" + sort
	}

	return &model.Pagination{Sort: sort, Order: order, Limit: limit, Offset: offset}
}

func ToHTML(s string) template.HTML {
	return template.HTML(s)
}

func AddData(c web.C, data map[string]interface{}) {
	for k, v := range data {
		c.Env[k] = v
	}
}

func DecodeJSON(req *http.Request, v interface{}) error {
	decoder := json.NewDecoder(req.Body)
	return decoder.Decode(v)
}

func MarshalJSON(v interface{}) (string, error) {
	j, err := json.Marshal(v)
	if err != nil {
		return "{}", nil
	}

	return string(j), err
}

func HTMLch(body string, code int, hdr map[string]string) *Response {
	return &Response{ContentType: HTML_CONTENT, Body: body, Code: code, Headers: hdr}
}

func HTMLc(body string, code int) *Response {
	return HTMLch(body, code, nil)
}

func HTML(body string) *Response {
	return HTMLc(body, http.StatusOK)
}

func JSONch(body string, code int, hdr map[string]string) *Response {
	return &Response{ContentType: JSON_CONTENT, Body: body, Code: code, Headers: hdr}
}

func JSONh(body string, hdr map[string]string) *Response {
	return JSONch(body, http.StatusOK, hdr)
}

func JSONc(body string, code int) *Response {
	return JSONch(body, code, nil)
}

func JSON(body string) *Response {
	return JSONc(body, http.StatusOK)
}

func PrepareJSONc(v interface{}, code int, errMsg string) *Response {
	result, err := MarshalJSON(v)
	if err != nil {
		log.Errorf("JSON Marshal error: %q", err)
		return ErrInternalServer(errMsg)
	}

	return JSONc(result, code)
}

func PrepareJSON(v interface{}, errMsg string) *Response {
	return PrepareJSONc(v, http.StatusOK, errMsg)
}

func ErrBadRequest(msg string) *Response {
	err := errors.New("bad_request", msg)
	return JSONc(err.JSON(), http.StatusBadRequest)
}

func ErrInternalServer(msg string) *Response {
	err := errors.New("server_error", msg)
	return JSONc(err.JSON(), http.StatusInternalServerError)
}

func ErrConflict(msg string) *Response {
	err := errors.New("server_error", msg)
	return JSONc(err.JSON(), http.StatusConflict)
}

func ErrForbidden(msg string) *Response {
	err := errors.New("forbidden", msg)
	return JSONc(err.JSON(), http.StatusForbidden)
}

func ErrValidation(msg string) *Response {
	err := errors.New("validation", msg)
	return JSONc(err.JSON(), http.StatusBadRequest)
}

func CheckError(err error) int {
	if err != nil {
		log.Errorf("Error: %v", err)
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

func CheckErrorp(err error, p string) int {
	if err != nil {
		log.Errorf("%v: %v", p, err)
		return http.StatusInternalServerError
	}
	return http.StatusOK
}
