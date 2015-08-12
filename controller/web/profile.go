package web

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/jeevatkm/urlite/model"

	log "github.com/Sirupsen/logrus"
	"github.com/jeevatkm/urlite/util"

	"github.com/jeevatkm/urlite/context"
	"github.com/zenazn/goji/web"

	. "github.com/jeevatkm/urlite/controller"
)

func Profile(a *context.App, c web.C, r *http.Request) (*Response, error) {
	//u := GetUser(c)

	AddData(c, Data{
		"IsProfile": true,
	})

	content, err := a.Parse("profile", c.Env)
	code := CheckError(err)

	AddData(c, Data{
		"Title":   "Profile | urlite",
		"Content": ToHTML(content),
	})

	body, err := a.ParseF(c.Env)
	code = CheckError(err)

	return &Response{ContentType: HTML_CONTENT, Body: body, Code: code}, err

}

func ProfilePost(a *context.App, c web.C, r *http.Request) (*Response, error) {
	action := strings.TrimSpace(r.FormValue("a"))

	if action == "CP" { // Change password
		return handleChangePassword(a, c, r)
	}

	return handleProfileUpdate(a, c, r)
}

func handleChangePassword(a *context.App, c web.C, r *http.Request) (*Response, error) {
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
		return &Response{ContentType: JSON_CONTENT, Body: body, Code: http.StatusOK}, err
	}

	user, _ := model.GetUserByEmail(a.DB(), u.Email)
	user.Password = util.HashPassword(np)
	user.UpdatedTime = time.Now()
	user.UpdatedBy = user.ID.Hex()
	err := model.UpdateUser(a.DB(), user)
	if err != nil {
		log.Errorf("Error occurred while update: %v", err)
		body = `{
			"status":"error",
			"message": "Error occurred while updating user"
			}`
		return &Response{ContentType: JSON_CONTENT, Body: body, Code: http.StatusOK}, err
	}

	log.Debugf("Password updated successfully for %v", user.Email)
	body = `{
			"status":"success",
			"message": "Password updated successfully."
			}`

	return &Response{ContentType: JSON_CONTENT, Body: body, Code: http.StatusOK}, nil
}

func handleProfileUpdate(a *context.App, c web.C, r *http.Request) (*Response, error) {

	return &Response{ContentType: JSON_CONTENT, Body: "{}", Code: http.StatusOK}, nil
}
