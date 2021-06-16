package controllers

import (
	"fmt"
	"goodadvice/v1/models"
	"html/template"
	"log"
	"net/http"
	"regexp"
)

type adminController struct {
	adminIDPattern *regexp.Regexp
}

type Movement struct {
	MovementType string
	Options string
}

type LoadMovements struct {
	something string
	// Movement type DDL values
	DDLoptions []string
}

// html templates
var admintpl = template.Must(template.ParseFiles("htmlpages/admin.html"))
var adminmovementstpl = template.Must(template.ParseFiles("htmlpages/adminmovements.html"))
var adminuserstpl = template.Must(template.ParseFiles("htmlpages/adminusers.html"))

// entry point from front.go
func newAdminController() *adminController {
	return &adminController{
		adminIDPattern: regexp.MustCompile(`^/wod/(\d+)/?`),
	}
}

// ServeHTTP - Entry point from front.go
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
				submit := r.FormValue("submit")
				if submit == "users" {
					pageLoadUsers(w)
				} else if submit == "movements" {
					pageLoadMovements(w,r)
				} else if submit == "workouts" {
					pageLoadWorkouts()
				} else {
					pageLoadAdmin(w)
				}
			case http.MethodPost:
				movements := r.FormValue("addmovement")
				err := r.ParseForm()
				if err != nil {
					log.Fatalf("Failed to decode postFormByteSlice: %v", err)
				}
				//movements := r.FormValue("movements")
				if movements == "addmovement" {
					saveMovement(w,r)
					pageLoadMovements(w,r)
				} else {
					fmt.Println("notmovements lol")
				}
			default:
				w.WriteHeader(http.StatusNotImplemented)
			}
		}
	}
}



// pageLoadAdmin
func pageLoadAdmin(w http.ResponseWriter) {
	// default load todays wod if there is one for quick edits
	v := models.GetVersion()
	admintpl.Execute(w, v)
}

// pageLoadWorkouts - switch to workouts template do work there
func pageLoadWorkouts() {
}

// pageLoadUsers - switch to users template do work there
func pageLoadUsers(w http.ResponseWriter) {
	u := models.GetUsers()
	fmt.Println(u)
	adminuserstpl.Execute(w, u)
}

// pageLoadMovements - switch to Movements template - This doesn't dynamically populate the DDL yet so I hard coded the page.
func pageLoadMovements(w http.ResponseWriter, r *http.Request) {
	//var lm LoadMovements
	mt := models.GetMovementTypes()
	adminmovementstpl.Execute(w, mt)
}

// saveMovement - save new movement type from adminmovements.html
func saveMovement(w http.ResponseWriter, r *http.Request) {
	m := r.FormValue("movement")
	mt := r.FormValue("movementtypes")
	models.SaveMovement(m,mt)
}


