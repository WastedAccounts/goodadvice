package controllers

import (
	"fmt"
	"goodadvice/v1/models"
	"html/template"
	"net/http"
	"regexp"
)

// loginController - handles values after .com/
type loginController struct {
	loginIDPattern *regexp.Regexp
}

// newLoginController - entry point from front.go
func newLoginController() *loginController {
	return &loginController{
		loginIDPattern: regexp.MustCompile(`^/login/(\d+)/?`),
	}
}

var logintpl = template.Must(template.ParseFiles("htmlpages/login.html", "htmlpages/templates/headerguest.html", "htmlpages/templates/footerguest.html"))

func (lc loginController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/login" {
		switch r.Method {
		case http.MethodGet:
			lc.pageLoad(w, r)
		case http.MethodPost:
			if models.Login(w, r) == false {
				// Do work to tell user their password failed
				fmt.Println("login failed")
				loginFailed(w, r)
			} else {
				// do work to continue logging in user and set up session
				fmt.Println("login succeeded")
				models.CreateSession(w, r)
			}
			http.Redirect(w, r, "/", 302)
		default:
			fmt.Println("status not implemented")
			w.WriteHeader(http.StatusNotImplemented)
		}
	}
}

// pageLoad - Loads login page template - no data
func (lc *loginController) pageLoad(w http.ResponseWriter, r *http.Request) {
	//l := models.login("yes")
	logintpl.Execute(w, nil)
}

// loginFailed - displays failed log in message
func loginFailed(w http.ResponseWriter, r *http.Request) {
	// do work to help you figure why their log in failed
	Msg := "Log in failed"
	logintpl.Execute(w, Msg)
}

// loginSuccess - Unused/does no work
func loginSuccess() {
	// do work to finish logging in user and setting up a session
}
