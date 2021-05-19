package controllers

import (
	"html/template"
	"goodadvice/v1/models"
	"net/http"
	"regexp"
)

type logoutController struct {
	logoutIDPattern *regexp.Regexp
}

// entry point from front.go
func newLogOutController() *logoutController {
	return &logoutController{
		logoutIDPattern: regexp.MustCompile(`^/login/(\d+)/?`),
	}
}

var	logouttpl = template.Must(template.ParseFiles("htmlpages/logout.html"))

func (loc logoutController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	models.LogOut(w,r)
	//fmt.Fprint(w, "<h1>Thanks for visiting!</h1>")
	logouttpl.Execute(w, nil)
}
