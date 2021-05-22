package controllers

import (
	"goodadvice/v1/models"
	"html/template"
	"log"
	"net/http"
	"regexp"
)

type adminController struct {
	adminIDPattern *regexp.Regexp
}


// html templates
var admintpl = template.Must(template.ParseFiles("htmlpages/admin.html"))
//var adminOTHERtpl = template.Must(template.ParseFiles("htmlpages/adminOTHER.html"))

// entry point from front.go
func newAdminController() *adminController {
	return &adminController{
		adminIDPattern: regexp.MustCompile(`^/wod/(\d+)/?`),
	}
}

//ServeHTTP
// Entry point for the /addwod page
// Comes in from front.go
func (ac adminController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Validate session
	c := models.ValidateSession(w, r)
	if c.Exists == false {
		http.Redirect(w, r, "/login", 401)
	} else if c.Isadmin == false {
		http.Redirect(w, r, "/login", 401)
	} else {
		if r.URL.Path == "/admin" {
			switch r.Method {
			case http.MethodGet:
				if r.FormValue("date") == "" {
					pageLoadAdmin(w,r)
				} else {
					//loadWod(w,r)
				}

			case http.MethodPost:
				err := r.ParseForm()
				if err != nil {
					log.Fatalf("Failed to decode postFormByteSlice: %v", err)
				}
				//if Edit == true {
				//	//editWOD(w, r)
				//} else {
				//	//postWOD(w, r, c.Uid)
				//}
			default:
				w.WriteHeader(http.StatusNotImplemented)
			}
		}
	}
}

func pageLoadAdmin(w http.ResponseWriter, r *http.Request) {
	// default load todays wod if there is one for quick edits
	//wo := models.GetWODGuest()
	//Edit = false -- used on addworkout to control edit vs new templates
	admintpl.Execute(w, nil)
}
