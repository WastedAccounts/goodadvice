package auth

import (
	"fmt"
	"goodadvice/v1/models/auth"
	"html/template"
	"net/http"
	"regexp"
)

type authController struct {
	authIDPattern *regexp.Regexp
}


type Webvals struct {
	Msg string
}


var	confirmtpl = template.Must(template.ParseFiles("htmlpages/auth/confirmemail.html"))


func NewAuthController() *authController {
	return &authController{
		authIDPattern: regexp.MustCompile(`^/auth/(\d+)/?`),
	}
}

// set cookies: https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/06.1.html
func (authc authController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/auth/confirmemail" {
		switch r.Method  {
		case http.MethodGet:
			authc.confirmEmailLoad(w, r)
		case http.MethodPost:
			i,msg := auth.ConfirmEmail(w, r)
			if i == true {
				//redirect to home page
				http.Redirect(w, r, "/profile/userprofile", 302)
			} else if i == false {
				authc.confirmEmailLoadFailed(w,r,msg)
			}
		default:
			fmt.Println("status not implemented")
			w.WriteHeader(http.StatusNotImplemented)
		}
	}
}

func (authc *authController) confirmEmailLoad(w http.ResponseWriter, r *http.Request) {
	confirmtpl.Execute(w, nil)
}

func (authc *authController) confirmEmailLoadFailed(w http.ResponseWriter, r *http.Request, m string) {
	var webv = Webvals{
		Msg: 		  m,
	}
	confirmtpl.Execute(w, webv)
}