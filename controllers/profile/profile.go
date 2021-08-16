package profile

import (
	"fmt"
	"goodadvice/v1/models"
	"goodadvice/v1/models/profile"
	"html/template"
	//"log"
	"net/http"
	"regexp"
)

type profileController struct {
	profileIDPattern *regexp.Regexp
}
// entry point from front.go
func NewProfileController() *profileController {
	return &profileController{
		profileIDPattern: regexp.MustCompile(`^/profile/(\d+)/?`),
	}
}

var	aboutmetpl = template.Must(template.ParseFiles("htmlpages/profile/aboutme.html"))
var	goalstpl = template.Must(template.ParseFiles("htmlpages/profile/goals.html"))
var	personalrecordstpl = template.Must(template.ParseFiles("htmlpages/profile/personalrecords.html"))

// set cookies: https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/06.1.html
func (pc profileController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Validate session
	c := models.ValidateSession(w, r)
	if c.Exists == false {
		http.Redirect(w, r, "/login", 401)
		//} else if c.Isadmin != false {
		//	http.Redirect(w, r, "/login", 401)
	} else {
		if r.RequestURI == "/profile/aboutme" {
			switch r.Method {
			case http.MethodGet:
				//submit := r.FormValue("submit")
				//err := r.ParseForm()
				//if err != nil {
				//	log.Fatalf("Failed to decode getFormByteSlice: %v", err)
				//}
				//if submit == "" {
				//
				//} else if submit == "Add Record" {
				//
				//}
				pc.pageLoadAboutMe(w, r, c.Uid)
			case http.MethodPost:
				if models.Login(w, r) == false {
				} else {
				}
				http.Redirect(w, r, "/", 302)
			default:
				fmt.Println("status not implemented")
				w.WriteHeader(http.StatusNotImplemented)
			}
		} else if r.RequestURI == "/profile/goals" {
			switch r.Method {
			case http.MethodGet:
				//submit := r.FormValue("submit")
				//err := r.ParseForm()
				//if err != nil {
				//	log.Fatalf("Failed to decode getFormByteSlice: %v", err)
				//}
				//if submit == "" {
				//
				//} else if submit == "Add Record" {
				//
				//}
				pc.pageLoadGoals(w, r, c.Uid)
			case http.MethodPost:
				if models.Login(w, r) == false {
				} else {
				}
				http.Redirect(w, r, "/", 302)
			default:
				fmt.Println("status not implemented")
				w.WriteHeader(http.StatusNotImplemented)
			}
		} else if r.RequestURI == "/profile/personalrecords" {
			switch r.Method {
			case http.MethodGet:
				pc.pageLoadPersonalRecords(w, r, c.Uid)
			case http.MethodPost:
				if models.Login(w, r) == false {
				} else {
				}
				http.Redirect(w, r, "/", 302)
			default:
				fmt.Println("status not implemented")
				w.WriteHeader(http.StatusNotImplemented)
			}
		}
	}
}


func (pc profileController) pageLoadAboutMe(w http.ResponseWriter, r *http.Request, id string) {
	// load PRs on page lade
	_,am := profile.PageLoadAboutMe(id)
	aboutmetpl.Execute(w,am) //(w, userprofile)
}

func (pc profileController) pageLoadGoals(w http.ResponseWriter, r *http.Request, id string) {
	// load PRs on page lade
	//_,am := profile.PageLoadGoals(id)
	rec,_ := profile.PageLoadAboutMe(id)
	goalstpl.Execute(w,rec) //(w, userprofile)
}

func (pc profileController) pageLoadPersonalRecords(w http.ResponseWriter, r *http.Request, id string) {
	// load PRs on page lade
	rec := profile.PageLoadPersonalRecords(id)
	personalrecordstpl.Execute(w,rec) //(w, userprofile)
}
