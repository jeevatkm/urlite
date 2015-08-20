package web

import (
	"net/http"

	"github.com/jeevatkm/urlite/dashboard/controller/api"
	"github.com/jeevatkm/urlite/model"
	"github.com/zenazn/goji/web"

	. "github.com/jeevatkm/urlite/dashboard/context"
	. "github.com/jeevatkm/urlite/dashboard/controller"
)

func Home(ctx *Context, c web.C, r *http.Request) *Response {
	content, err := ctx.Parse("home", c.Env)
	code := CheckError(err)

	AddData(c, Data{
		"IsHome":  true,
		"Title":   "Home | urlite",
		"Content": ToHTML(content),
	})

	body, err := ctx.ParseF(c.Env)
	code = CheckError(err)

	return HTMLc(body, code)
}

func Urlite(ctx *Context, c web.C, r *http.Request) *Response {
	c.Env["liteReq"] = &model.UrliteRequest{LongUrl: r.FormValue("longUrl"),
		Domain:     r.FormValue("selectedDomain"),
		CustomName: r.FormValue("customLinkName")}

	return api.HandleUrlite(ctx, c, r)
}
