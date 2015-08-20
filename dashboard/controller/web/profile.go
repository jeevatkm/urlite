package web

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/jeevatkm/urlite/model"
	"github.com/jeevatkm/urlite/util"
	"github.com/zenazn/goji/web"

	log "github.com/Sirupsen/logrus"
	. "github.com/jeevatkm/urlite/dashboard/context"
	. "github.com/jeevatkm/urlite/dashboard/controller"
)

func Profile(ctx *Context, c web.C, r *http.Request) *Response {
	//u := GetUser(c)

	AddData(c, Data{
		"IsProfile": true,
	})

	content, err := ctx.Parse("profile", c.Env)
	code := CheckError(err)

	AddData(c, Data{
		"Title":   "Profile | urlite",
		"Content": ToHTML(content),
	})

	body, err := ctx.ParseF(c.Env)
	code = CheckError(err)

	return HTMLc(body, code)

}

func ProfilePost(ctx *Context, c web.C, r *http.Request) *Response {
	action := strings.TrimSpace(r.FormValue("a"))

	if action == "CP" { // Change password
		return handleChangePassword(ctx, c, r)
	}

	return handleProfileUpdate(ctx, c, r)
}

func handleChangePassword(ctx *Context, c web.C, r *http.Request) *Response {
	u := GetUser(c)
	log.Debugf("Action change password for %v", u.Email)
	ep, np := r.FormValue("existingPassword"), r.FormValue("newPassword")
	body := ""

	if !util.ComparePassword(u.Password, ep) {
		err := errors.New("Existing password is incorrect")
		log.Errorf("%v", err)
		body = `{
			"status":"error",
			"message": "Existing password is incorrect"
			}`
		return JSON(body)
	}

	db := ctx.DB(&c)
	user, _ := model.GetUserByEmail(db, u.Email)
	user.Password = util.HashPassword(np)
	user.UpdatedTime = time.Now()
	user.UpdatedBy = user.ID.Hex()
	err := model.UpdateUser(db, user)
	if err != nil {
		log.Errorf("Error occurred while update: %v", err)
		body = `{
			"status":"error",
			"message": "Error occurred while updating user"
			}`
		return JSON(body)
	}

	log.Debugf("Password updated successfully for %v", user.Email)
	body = `{
			"status":"success",
			"message": "Password updated successfully."
			}`

	return JSON(body)
}

func handleProfileUpdate(ctx *Context, c web.C, r *http.Request) *Response {

	return JSON("{}")
}
