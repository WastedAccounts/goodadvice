package controllers

import (
	"goodadvice/v1/models"
	"html/template"
	"log"
	"net/http"
	"regexp"
)

type addWorkoutController struct {
	addWorkoutIDPattern *regexp.Regexp
}

//type AddWorkoutPageLoad struct {
//	WoID int
//	WoName string
//	WoStrength string
//	WoPace string
//	WoConditioning string
//	WoDate string
//	WoDOW string
//	UsrID string
//	UsrNoteID int
//	UsrName string
//	UsrNotes string
//}
//
//type Workout struct {
//	ID int
//	Name string
//	Strength string
//	Pace string
//	Conditioning string
//	Date string
//}

// Used to control NEW vs EDIT templates
var Edit bool

// html templates
var addwodtpl = template.Must(template.ParseFiles("htmlpages/addworkout.html"))
var editwodtpl = template.Must(template.ParseFiles("htmlpages/editworkout.html"))

// entry point from front.go
func newAddWorkoutController() *addWorkoutController {
	return &addWorkoutController{
		addWorkoutIDPattern: regexp.MustCompile(`^/wod/(\d+)/?`),
	}
}

//ServeHTTP
// Entry point for the /addwod page
// Comes in from front.go
func (awc addWorkoutController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Validate session
	c := models.ValidateSession(w, r)
	if c.Exists == false {
		http.Redirect(w, r, "/login", 401)
	} else if c.Isadmin == false {
		http.Redirect(w, r, "/login", 401)
	} else {
		if r.URL.Path == "/addwod" {
			switch r.Method {
			case http.MethodGet:
				if r.FormValue("date") == "" {
					pageLoadAddWorkout(w)
				} else {
					loadWOD(w,r)
				}
			case http.MethodPost:
				err := r.ParseForm()
				if err != nil {
					log.Fatalf("Failed to decode postFormByteSlice: %v", err)
				}
				if Edit == true {
					editWOD(w,r)
				} else {
					postWOD(w,r,c.Uid)
				}
			default:
				w.WriteHeader(http.StatusNotImplemented)
			}
		}
	}
}

// pageLoadAddWorkout - initial page load
func pageLoadAddWorkout(w http.ResponseWriter) {
	// default load todays wod if there is one for quick edits
	//wo := models.GetWODGuest()
	Edit = false
	addwodtpl.Execute(w, nil)
}

// postWOD - write workout to the database and reloads it to the page
func postWOD(w http.ResponseWriter, r *http.Request, uid string) {
	wo := models.AddWOD(r, uid)
	if wo.Message == "" {
		Edit = true
		editwodtpl.Execute(w, wo)
	} else {
		Edit = false
		addwodtpl.Execute(w, wo)
	}
}

//loadWod - loads workout for selected date
func loadWOD(w http.ResponseWriter, r *http.Request) {
	wo := models.GetAddWODbydate(r.FormValue("date"))
	Edit = true
	editwodtpl.Execute(w, wo)
}

// editWOD - saves changes made to the workout and reloads it to the page
func editWOD (w http.ResponseWriter, r *http.Request) {
	models.EditAddWOD(r)
	Edit = true
	wo := models.GetAddWODbydate(r.FormValue("date"))
	editwodtpl.Execute(w, wo)
}