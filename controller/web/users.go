package web

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/jeevatkm/urlite/context"
	"github.com/jeevatkm/urlite/model"
	"github.com/jeevatkm/urlite/util"
	"github.com/jeevatkm/urlite/util/random"
	"github.com/zenazn/goji/web"

	log "github.com/Sirupsen/logrus"
	. "github.com/jeevatkm/urlite/controller"
)

func Users(a *context.App, c web.C, r *http.Request) *Response {
	AddData(c, Data{
		"IsUsers": true,
		"Domains": a.Domains,
	})

	content, err := a.Parse("users", c.Env)
	code := CheckError(err)

	AddData(c, Data{
		"Title":   "Users | urlite",
		"Content": ToHTML(content),
	})

	body, err := a.ParseF(c.Env)
	code = CheckError(err)

	return HTMLc(body, code)
}

func UsersData(a *context.App, c web.C, r *http.Request) *Response {

	page := ParsePagination(r)
	pageResult, _ := model.GetUsersByPage(a.DB(), bson.M{}, page)
	body, err := MarshalJSON(pageResult)
	if err != nil {
		log.Errorf("JSON Marshal error: %q", err)
		body = `{"status":"error","message": "Unable to users list."}`
		return JSON(body)
	}

	return JSON(body)
}

func UsersPost(a *context.App, c web.C, r *http.Request) *Response {
	email := r.FormValue("uemail")
	euser, _ := model.GetUserByEmail(a.DB(), email)
	if euser != nil {
		m := fmt.Sprintf("User already exists with given email id: '%v'", email)
		log.Error(m)
		SetErrorAlert(c, m)
		return Users(a, c, r)
	}

	udomains, upermissions := r.FormValue("sUserDomains"), r.FormValue("sUserPermissions")
	password := random.GenerateUserPassword(8)
	user := &model.User{Email: email,
		Password:    util.HashPassword(password),
		IsActive:    true,
		Permissions: strings.Split(upermissions, ","),
		Bearer:      random.GenerateBearerToken(),
		Domains:     strings.Split(udomains, ","),
		CreatedBy:   GetUser(c).ID.Hex(),
		CreatedTime: time.Now()}

	err := model.CreateUser(a.DB(), user)
	if err != nil {
		log.Errorf("Unable to create user: %q", err)
		SetErrorAlert(c, "Unable to create user: "+email)
		return Users(a, c, r)
	}

	c.Env["UserCount"] = model.GetActiveUserCount(a.DB())
	msg := "Successfully added user: " + email + ", notification has been sent with user password."
	SetSuccessAlert(c, msg)
	log.Info(msg)
	return Users(a, c, r)
}
