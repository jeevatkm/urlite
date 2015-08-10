package web

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jeevatkm/urlite/util"

	"github.com/jeevatkm/urlite/context"
	"github.com/jeevatkm/urlite/model"
	"github.com/jeevatkm/urlite/util/random"
	"github.com/zenazn/goji/web"

	log "github.com/Sirupsen/logrus"
	. "github.com/jeevatkm/urlite/controller"
)

func Users(a *context.App, c web.C, r *http.Request) (*Response, error) {
	users, _ := model.GetAllUsers(a.DB())
	AddData(c, Data{
		"IsUsers": true,
		"Users":   users,
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

	return &Response{ContentType: HTML_CONTENT, Body: body, Code: code}, err
}

func UsersPost(a *context.App, c web.C, r *http.Request) (*Response, error) {
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
	log.Debugf("User password: %v", password) // SHOULD BE REMOVED
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
