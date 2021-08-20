package controllers

import (
	"fmt"
	"goodadvice/v1/models"
	"html/template"
	"net/http"
	"regexp"
)


type signupController struct {
	signupIDPattern *regexp.Regexp
}


type Webvals struct {
	Userval string
	Firstnameval string
	Emailval string
	Msg string
}

// html templates
var	signuptpl = template.Must(template.ParseFiles("htmlpages/signup.html"))

// newSignupController
func newSignupController() *signupController {
	return &signupController{
		signupIDPattern: regexp.MustCompile(`^/signup/(\d+)/?`),
	}
}

// set cookies: https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/06.1.html
func (sc signupController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/signup" {
		switch r.Method  {
		case http.MethodGet:
			sc.pageLoad(w, r)
		case http.MethodPost:
			if r.FormValue("firstname")  == "" {
				e := "Missing first name"
				sc.pageReload(w, r, e)
			}else if r.FormValue("email")  == "" {
				e := "Missing email"
				sc.pageReload(w, r, e)
			}else if r.FormValue("username")  == "" {
				e := "Missing username"
				sc.pageReload(w, r, e)
			}else if r.FormValue("password") != r.FormValue("confirmpassword"){
				e := "Passwords do not match"
				sc.pageReload(w, r, e)
			} else {
				if models.CheckEmail(r) == true {
					e := "Email already exists"
					sc.pageReload(w, r, e)
				} else if  models.CheckUsername(r) == true {
					e := "Username not available"
					sc.pageReload(w, r, e)
				} else {
					models.Signup(w, r)
					// now we issue a cookie to the client and activate a session -- client account is still inactive here
					models.CreateSession(w, r)
					//http.Redirect(w, r, r.Header.Get("/login"), 302)
					//sc.confirmemailLoad(w, r)
					http.Redirect(w, r, r.Header.Get("/auth/confirmemail"), 302)
				}
			}
		default:
			fmt.Println("status not implemented")
			w.WriteHeader(http.StatusNotImplemented)
		}
	}
}

// pageLoad - load sign up template
func (sc *signupController) pageLoad(w http.ResponseWriter, r *http.Request) {
	signuptpl.Execute(w, nil)
}

func (sc *signupController) pageReload(w http.ResponseWriter, r *http.Request,e string) {
	var webv = Webvals{
		Userval:      r.FormValue("username"),
		Firstnameval: r.FormValue("firstname"),
		Emailval:     r.FormValue("email"),
		Msg: 		  e,
	}
	signuptpl.Execute(w,webv)
}