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

//type Movement struct {
//	MovementType string
//	Options string
//}
//
//type Moves struct {
//	Movement string
//	MovementType string
//}
type LoadMovementsData struct {
	MovementList []models.Movements
	// Movement type DDL values
	Ddloptions []string
}

// html templates
var admintpl = template.Must(template.ParseFiles("htmlpages/admin.html"))
var adminmovementstpl = template.Must(template.ParseFiles("htmlpages/adminmovements.html"))
var adminuserstpl = template.Must(template.ParseFiles("htmlpages/adminusers.html"))
var adminusertpl = template.Must(template.ParseFiles("htmlpages/adminuser.html"))

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
				} else if submit == "user" {
					pageLoadUser(w,r)
				} else {
					pageLoadAdmin(w)
				}
			case http.MethodPost:
				submit := r.FormValue("submit")
				err := r.ParseForm()
				changeValue := map[string]bool {
					"Activate": true,
					"Deactivate": true,
					"User": true,
					"Moderator": true,
					"Admin": true,
				}
				if err != nil {
					log.Fatalf("Failed to decode postFormByteSlice: %v", err)
				}
				//movements := r.FormValue("movements")
				if submit == "addmovement" {
					saveMovement(w,r)
					pageLoadMovements(w,r)
				} else if changeValue[submit] {
					models.UpdateUser(r.FormValue("userID"),submit)
					pageLoadUser(w,r)
				} else {
					fmt.Println("notmovements lol")
				}
			default:
				w.WriteHeader(http.StatusNotImplemented)
			}
		}
	}
}

func submitTrue(ss ...string) bool {
	for _, s := range ss {
		if s == "" {
			return true
		}
	}
	return false
}

func pageLoadUser(w http.ResponseWriter, r *http.Request ) {
	u := models.AdminGetUser(r.FormValue("userID"))
	adminusertpl.Execute(w, u)
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
	adminuserstpl.Execute(w, u)
}

// pageLoadMovements - switch to Movements template - This doesn't dynamically populate the DDL yet so I hard coded the page.
func pageLoadMovements(w http.ResponseWriter, r *http.Request) {
	//var lm LoadMovements
	m := models.GetMovements()
	mt := models.GetMovementTypes()
	data := LoadMovementsData{
		MovementList: m,
		Ddloptions: mt,
	}
	fmt.Printf("%v",data)
	adminmovementstpl.Execute(w,data)
}

// saveMovement - save new movement type from adminmovements.html
func saveMovement(w http.ResponseWriter, r *http.Request) {
	m := r.FormValue("movement")
	mt := r.FormValue("movementtypes")
	models.SaveMovement(m,mt)
}


