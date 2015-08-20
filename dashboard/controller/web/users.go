package web

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/jeevatkm/urlite/model"
	"github.com/jeevatkm/urlite/util"
	"github.com/jeevatkm/urlite/util/random"
	"github.com/zenazn/goji/web"

	log "github.com/Sirupsen/logrus"
	. "github.com/jeevatkm/urlite/dashboard/context"
	. "github.com/jeevatkm/urlite/dashboard/controller"
)

func Users(ctx *Context, c web.C, r *http.Request) *Response {
	AddData(c, Data{
		"IsUsers": true,
		"Domains": ctx.Domains,
	})

	content, err := ctx.Parse("users", c.Env)
	code := CheckError(err)

	AddData(c, Data{
		"Title":   "Users | urlite",
		"Content": ToHTML(content),
	})

	body, err := ctx.ParseF(c.Env)
	code = CheckError(err)

	return HTMLc(body, code)
}

func UsersData(ctx *Context, c web.C, r *http.Request) *Response {
	page, db := ParsePagination(r), ctx.DB(&c)
	pageResult, _ := model.GetUsersByPage(db, bson.M{}, page)

	body, err := MarshalJSON(pageResult)
	if err != nil {
		log.Errorf("JSON Marshal error: %q", err)
		body = `{"status":"error","message": "Unable to users list."}`
		return JSON(body)
	}

	return JSON(body)
}

func UsersValidate(ctx *Context, c web.C, r *http.Request) *Response {
	var response *Response

	// Email address exists check validation
	email := strings.TrimSpace(r.FormValue("uemail"))
	log.Debugf("uemail: %v", email)
	if len(email) > 0 {
		euser, _ := model.GetUserByEmail(ctx.DB(&c), email)
		if euser != nil {
			response = ErrValidation("Email address already exists")
		} else {
			response = JSON("{}")
		}
	}

	return response
}

func UsersPost(ctx *Context, c web.C, r *http.Request) *Response {
	email, db := strings.TrimSpace(r.FormValue("uemail")), ctx.DB(&c)
	euser, _ := model.GetUserByEmail(db, email)
	if euser != nil {
		m := fmt.Sprintf("User already exists with given email id: '%v'", email)
		log.Error(m)
		return ErrBadRequest(m)
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

	err := model.CreateUser(db, user)
	if err != nil {
		log.Errorf("Unable to create user: %q", err)
		return ErrInternalServer("Unable to create user, due to server issue")
	}

	userCount := model.GetActiveUserCount(db)
	msg := "Successfully added user: " + email + ", notification has been sent with user password."

	data := Data{
		"id":         ALERT_SUCCESS,
		"message":    msg,
		"user_count": userCount,
	}
	log.Info(msg)

	return PrepareJSON(data, GENERIC_ERROR_MSG)
}
