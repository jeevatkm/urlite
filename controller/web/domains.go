package web

import (
	"net/http"

	"github.com/jeevatkm/urlite/context"
	"github.com/jeevatkm/urlite/model"
	"github.com/zenazn/goji/web"

	. "github.com/jeevatkm/urlite/controller"
)

func Domains(a *context.App, c web.C, r *http.Request) (*Response, error) {
	domains, _ := model.GetAllDomain(a.DB())
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

	return &Response{ContentType: HTML_CONTENT, Body: body, Code: code}, err
}
