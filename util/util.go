package util

import (
	"html/template"

	log "github.com/Sirupsen/logrus"
	"github.com/zenazn/goji/web"
)

func Contains(a []string, s string) bool {
	for _, v := range a {
		if v == s {
			return true
		}
	}

	return false
}

func ToHTML(s string) template.HTML {
	return template.HTML(s)
}

func AddData(c web.C, data map[string]interface{}) {
	for k, v := range data {
		c.Env[k] = v
	}
}

func CheckError(err error) {
	if err != nil {
		log.Errorf("Error: %v", err)
	}
}

func CheckErrorp(err error, p string) {
	if err != nil {
		log.Errorf("%v: %v", p, err)
	}
}
