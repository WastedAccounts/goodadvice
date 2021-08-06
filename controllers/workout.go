package controllers

import (
	"goodadvice/v1/models"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type workoutController struct {
	workoutIDPattern *regexp.Regexp
}

type WorkoutPageLoad struct {
	WoID int
	WoName string
	WoStrength string
	WoPace string
	WoConditioning string
	WoDate string
	UsrID string
	UsrNoteID int
	UsrName string
	UsrNotes string
	UsrMinutes string
	UsrSeconds string
	UsrFirstname string
	UsrGreeting string
}

type Cookie struct {
	Exists bool
	Uid string
	Sessionkey string
}

// html templates
var	wodtpl = template.Must(template.ParseFiles("htmlpages/wod.html"))
var	wodguesttpl = template.Must(template.ParseFiles("htmlpages/guestwod.html"))
var	wodguesttpl2 = template.Must(template.ParseFiles("htmlpages/guestwod2.html"))
var guestfwtpl = template.Must(template.ParseFiles("htmlpages/guestframework.html"))

// entry point from front.go
func newWorkoutController() *workoutController {
	return &workoutController{
		workoutIDPattern: regexp.MustCompile(`^/wod/(\d+)/?`),
	}
}

//ServeHTTP
// Entry point for the /wod page
// Comes in from front.go
func (woc workoutController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check for a cookie first
	c := models.ValidateSession(w, r)
	// If there is no cookie found redirect to guest view
	if c.Exists == false {
		woc.GetWODGuest(w, r)
	} else {
		// At this point the user should be validated within two hour session time out
		// and a new cookie issued with new start time stamp
		uid := c.Uid//splitcookie[0]
		if r.URL.Path == "/wod" {
			switch r.Method {
			case http.MethodGet:
				//submit := r.FormValue("submit")
				// If a date is selected load workout from that date
				if r.FormValue("submitrandom") == "Random"{
					woc.randomWorkout(w, uid)
				} else if r.FormValue("date") != "" {
					woc.getWODbydate(w, r.FormValue("date"), uid)
				} else {
					// if no date is selected load today's workout
					woc.getWOD(w, uid)
				}
			case http.MethodPost:
				err := r.ParseForm()
				if err != nil {
					log.Fatalf("Failed to decode postFormByteSlice: %v", err)
				}
				woc.saveWODResults(w,r)
				woc.getWOD(w, uid)
				//// If a date is selected load workout from that date
				//if r.FormValue("date") != "" {
				//	woc.getWODbydate(w, r, r.FormValue("date"), r.PostFormValue("uid"))
				//} else {
				//	// if no date is selected load today's workout
				//	woc.getWOD(w, r,  r.PostFormValue("uid"))
				//}
			//case http.MethodPut:
			default:
				w.WriteHeader(http.StatusNotImplemented)
			}
		}
	}
}

// getWOD - displays WOD for the current date if there is one
func (woc *workoutController) getWOD(w http.ResponseWriter, uid string) {
	wo, won, usr := models.GetWOD(uid)
	wpl := WorkoutPageLoad{
		 wo.ID,
		wo.Name,
		wo.Strength,
		wo.Pace,
		wo.Conditioning,
		wo.Date,
		strconv.Itoa(usr.ID), //<-- I should be getting this from somewhere else NOT from the notes
		won.ID,
		usr.UserName,
		 won.Notes,
		 won.Minutes,
		 won.Seconds,
		 usr.FirstName,
		 usr.Greeting,
	}
	splitdate := strings.Split(wpl.WoDate, "T")
	wpl.WoDate = splitdate[0]
	wodtpl.Execute(w, wpl)
}

// GetWODGuest for non auth'd users
func (woc *workoutController) GetWODGuest(w http.ResponseWriter, r *http.Request) {
	wo := models.GetWODGuest()
	wpl := WorkoutPageLoad{
		wo.ID,
		wo.Name,
		wo.Strength,
		wo.Pace,
		wo.Conditioning,
		wo.Date,
		//wo.DOW,
		"", //strconv.Itoa(usr.ID), //<-- I should be getting this from somewhere else NOT from the notes
		0, //won.ID,
		"",//usr.UserName,
		"",//won.Notes,
		"",
		"",
		"",
		"",
	}
	wodguesttpl.Execute(w, wpl)
}

// getWODbydate - displays WOD for the current date if there is one
func (woc *workoutController) getWODbydate(w http.ResponseWriter, d string, uid string) {
	wo, won, usr := models.GetWODbydate(d, uid)
	wpl := WorkoutPageLoad{
		wo.ID,
		wo.Name,
		wo.Strength,
		wo.Pace,
		wo.Conditioning,
		wo.Date,
		//wo.DOW,
		strconv.Itoa(usr.ID), //<-- I should be getting this from somewhere else NOT from the notes
		won.ID,
		usr.UserName,
		won.Notes,
		won.Minutes,
		won.Seconds,
		usr.FirstName,
		usr.Greeting,
	}
	splitdate := strings.Split(wpl.WoDate, "T")
	wpl.WoDate = splitdate[0]
	wodtpl.Execute(w, wpl)
}

// saveWODResults - get notes for user for WOD being loaded
func (woc *workoutController) saveWODResults(w http.ResponseWriter, r *http.Request) {
	models.SaveWODResults(r)
	http.Redirect(w, r, r.Header.Get("Referer"), 302)
}

// randomWorkout = gets random workout from all workouts with weighted values for loved(2x times loved)/hated(1x times hated) rating
func (woc *workoutController) randomWorkout(w http.ResponseWriter, uid string) {
	//var date string
	date := models.GetRandomWorkout(uid)
	woc.getWODbydate(w,date,uid)
}