package controllers

import (
	"html/template"
	"log"
	"goodadvice/v1/models"
	"net/http"
	"regexp"
	"time"
)

type addWorkoutController struct {
	addWorkoutIDPattern *regexp.Regexp
}

type AddWorkoutPageLoad struct {
	WoID int `json:"woID"`
	WoName string `json:"woName"`
	WoStrength string `json:"woStrength"`
	WoPace string `json:"woPace"`
	WoConditioning string `json:"woConditioning"`
	WoDate time.Time `json:"woDate"`
	WoDOW string `json:"woDOW"`
	UsrID string `json:"usrID"`
	UsrNoteID int`json:"usrNoteID"`
	UsrName string `json:"usrName""`
	UsrNotes string `json:"usrNotes""`
}

type Workout struct {
	ID int `json:"ID"`
	Name string `json:"Name"`
	Strength string `json:"Strength"`
	Pace string `json:"Pace"`
	Conditioning string `json:"Conditioning"`
	Date string `json:"Date"`
	//DOW string `json:"DOW"`
}

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
					pageLoad(w,r)
				} else {
					loadWod(w,r)
				}

			case http.MethodPost:
				err := r.ParseForm()
				if err != nil {
					log.Fatalf("Failed to decode postFormByteSlice: %v", err)
				}
				if Edit == true {
					editWOD(w, r)
				} else {
					postWOD(w, r, c.Uid)
				}
			default:
				w.WriteHeader(http.StatusNotImplemented)
			}
		}
	}
}

func pageLoad(w http.ResponseWriter, r *http.Request) {
	// default load todays wod if there is one for quick edits
	//wo := models.GetWODGuest()
	Edit = false
	addwodtpl.Execute(w, nil)
}

func postWOD(w http.ResponseWriter, r *http.Request, uid string) {
	wo := models.AddWOD(w, r, uid)
	Edit = true
	editwodtpl.Execute(w, wo)
}

func loadWod(w http.ResponseWriter, r *http.Request) {
	wo := models.GetAddWODbydate(w,r)
	Edit = true
	editwodtpl.Execute(w, wo)
}

func editWOD (w http.ResponseWriter, r *http.Request) {
	models.EditAddWOD(w,r)
	Edit = true
	wo := models.GetAddWODbydate(w,r)
	editwodtpl.Execute(w, wo)
}