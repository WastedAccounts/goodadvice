package workouts

import (
	"fmt"
	"goodadvice/v1/models"
	"goodadvice/v1/models/workouts"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

type workoutController struct {
	workoutIDPattern *regexp.Regexp
}

type WorkoutPageLoad struct {
	WoID           int
	WoName         string
	WoStrength     string
	WoPace         string
	WoConditioning string
	WoDate         string
	UsrID          string
	UsrNoteID      int
	UsrName        string
	UsrNotes       string
	UsrMinutes     string
	UsrSeconds     string
	UsrFirstname   string
	UsrGreeting    string
}

//// UserAuth - Stores values for authenticating users around the app
//type UserAuth struct {
//	Exists bool
//	IsActive bool
//	Isadmin bool
//	IsCoach bool
//	Uid string
//	Path string
//	Sessionkey string
//}

// Used to control NEW vs EDIT templates
var Edit bool

// html templates
var adminaddwodtpl = template.Must(template.ParseFiles("htmlpages/workouts/adminaddworkout.html", "htmlpages/templates/header.html", "htmlpages/templates/footer.html"))
var admineditwodtpl = template.Must(template.ParseFiles("htmlpages/workouts/admineditworkout.html", "htmlpages/templates/header.html", "htmlpages/templates/footer.html"))
var userwodtpl = template.Must(template.ParseFiles("htmlpages/workouts/userwod.html", "htmlpages/templates/header.html", "htmlpages/templates/footer.html"))
var guestwodtpl = template.Must(template.ParseFiles("htmlpages/workouts/guestwod.html", "htmlpages/templates/headerguest.html", "htmlpages/templates/footerguest.html"))
var useraddworkouttpl = template.Must(template.ParseFiles("htmlpages/workouts/useraddworkout.html", "htmlpages/templates/header.html", "htmlpages/templates/footer.html"))
var usereditworkouttpl = template.Must(template.ParseFiles("htmlpages/workouts/usereditworkout.html", "htmlpages/templates/header.html", "htmlpages/templates/footer.html"))

// M - a map for passing multiple structs to htmltmpl
type M map[string]interface{}

// entry point from front.go
func NewWorkoutController() *workoutController {
	return &workoutController{
		workoutIDPattern: regexp.MustCompile(`^/workouts/wod/(\d+)/?`),
	}
}

//// START GetWOD Functions
// ServeHTTP - Entry point for the /wod page - Comes in from front.go
func (woc workoutController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check for a cookie first
	userauth := models.ValidateSession(w, r)
	// If there is no cookie found redirect to guest view
	if userauth.Exists == false || userauth.IsActive == false {
		woc.getWOD(w, r, userauth)
	} else {
		// At this point the user should be validated within 48 hour session time out
		// and a new cookie issued with new session start time stamp
		if r.URL.Path == "/workouts/wod" {
			switch r.Method {
			case http.MethodGet:
				//submit := r.FormValue("submit")
				// If a date is selected load workout from that date
				if r.FormValue("random") == "Random" {
					woc.randomWorkout(w, userauth)
				} else if r.FormValue("date") != "" {
					woc.getWODbydate(w, r.FormValue("date"), userauth)
				} else {
					// if no date is selected load today's workout
					woc.getWOD(w, r, userauth)
				}
			case http.MethodPost:
				err := r.ParseForm()
				if err != nil {
					log.Fatalf("Failed to decode postFormByteSlice: %v", err)
				}
				woc.SaveWorkoutResults(w, r)
				woc.getWOD(w, r, userauth)
			case http.MethodPut:
			default:
				w.WriteHeader(http.StatusNotImplemented)
			}
		} else if r.URL.Path == "/workouts/addwod" {
			// This is for working with the Daily WOD.
			// If an admin wants to add an additional workout
			// they should do it as a user
			switch r.Method {
			case http.MethodGet:
				if r.FormValue("date") == "" {
					pageLoadAddWorkout(w)
				} else {
					loadWODEdit(w, r)
				}
			case http.MethodPost:
				err := r.ParseForm()
				if err != nil {
					log.Fatalf("Failed to decode postFormByteSlice: %v", err)
				}
				if Edit == true {
					editWOD(w, r, userauth.Uid,true)
				} else {
					postWOD(w, r, userauth.Uid, true)
				}
			default:
				w.WriteHeader(http.StatusNotImplemented)
			}
		} else if r.URL.Path == "/workouts/createwod" {
			switch r.Method {
			case http.MethodGet:
				// Parse woid out of URL
				u, _ := url.Parse(r.RequestURI)
				woid, _ := url.ParseQuery(u.RawQuery)
				// Call for WOD by woid
				woc.customizeWOD(w, r, woid.Get("woid"), userauth, "")
			case http.MethodPost:
				err := r.ParseForm()
				if err != nil {
					log.Fatalf("Failed to decode postFormByteSlice: %v", err)
				}
				//postWOD(w, r, userauth.Uid, false)
				if Edit == true {
					fmt.Println("here2")
					editWOD(w, r, userauth.Uid, false)
					// EDIT DOESN"T WORK AS USER
					// Reason is the Uid and Date get caught
					// but no message is sent.
				} else {
					fmt.Println("here3")
					postWOD(w, r, userauth.Uid, false)
					// pOST FAILS BECUASE Edit var is stuck on True after first
					// round of posting
					//Should just redirect to and /workouts/editWOD url
				}
			default:
				w.WriteHeader(http.StatusNotImplemented)
			}
		} else if r.URL.Path == "/workouts/editwod" {
			//switch r.Method {
			//case http.MethodGet:
			//	w.WriteHeader(http.StatusNotImplemented)
			//case http.MethodPost:
			//	fmt.Println("here2")
			//	userEditWOD(w, r, userauth.Uid)
			//	// EDIT DOESN"T WORK AS USER
			//	// Reason is the Uid and Date get caught
			//	// but no message is sent.
			//default:
			//	w.WriteHeader(http.StatusNotImplemented)
		}
	}
}

// getWOD - Gets WOD and User data if they're logged in and depending on Role loads the correct template
func (woc *workoutController) getWOD(w http.ResponseWriter, r *http.Request, userauth models.UserAuth) {
	// Call GetWOD to get structs for page load
	wo, won, usr := workouts.GetWOD(userauth.Uid, r)

	// Map structs to a single var to load into templates
	data := M{
		"wo":  wo,
		"won": won,
		"usr": usr,
	}
	if userauth.Exists == true && userauth.IsAdmin == true && userauth.IsActive == true {
		// If Admin
		userwodtpl.Execute(w, data)
	} else if userauth.Exists == true && userauth.IsCoach == true && userauth.IsActive == true {
		// If Coach
		userwodtpl.Execute(w, data)
	} else if userauth.Exists == true && userauth.IsActive == true {
		// If User
		userwodtpl.Execute(w, data)
	} else {
		// If Guest
		guestwodtpl.Execute(w, data)
	}
}

//// GetWODGuest for non auth'd users
//func (woc *workoutController) GetWODGuest(w http.ResponseWriter, r *http.Request) {
//	wo := workouts.GetWODGuest()
//	wpl := WorkoutPageLoad{
//		wo.ID,
//		wo.Name,
//		wo.Strength,
//		wo.Pace,
//		wo.Conditioning,
//		wo.Date,
//		//wo.DOW,
//		"", //strconv.Itoa(usr.ID), //<-- I should be getting this from somewhere else NOT from the notes
//		0, //won.ID,
//		"",//usr.UserName,
//		"",//won.Notes,
//		"",
//		"",
//		"",
//		"",
//	}
//	wodguesttpl.Execute(w, wpl)
//}

// getWODbydate - displays WOD for the current date if there is one
func (woc *workoutController) getWODbydate(w http.ResponseWriter, d string, userauth models.UserAuth) {
	wo, won, usr := workouts.GetWODbydate(d, userauth.Uid)
	// Map structs to a single var to load into templates
	data := M{
		"wo":  wo,
		"won": won,
		"usr": usr,
	}
	if userauth.Exists == true && userauth.IsAdmin == true && userauth.IsActive == true {
		// If Admin
		userwodtpl.Execute(w, data)
	} else if userauth.Exists == true && userauth.IsCoach == true && userauth.IsActive == true {
		// If User
		userwodtpl.Execute(w, data)
	} else if userauth.Exists == true && userauth.IsActive == true {
		// If User
		userwodtpl.Execute(w, data)
	} else {
		// If Guest
		guestwodtpl.Execute(w, data)
	}
}

// SaveWorkoutResults - get notes for user for WOD being loaded
func (woc *workoutController) SaveWorkoutResults(w http.ResponseWriter, r *http.Request) {
	workouts.SaveWorkoutResults(r)
	http.Redirect(w, r, r.Header.Get("Referer"), 302)
}

// randomWorkout = gets random workout from all workouts with weighted values for loved(2x times loved)/hated(1x times hated) rating
func (woc *workoutController) randomWorkout(w http.ResponseWriter, userauth models.UserAuth) {
	//var date string
	date := workouts.GetRandomWorkout()
	woc.getWODbydate(w, date, userauth)
}

// customizeWOD - Click the customize button and we gather the needed details and send the user off to make their own workout.
func (woc *workoutController) customizeWOD(w http.ResponseWriter, r *http.Request, woid string, userauth models.UserAuth, msg string) {
	data := workouts.GetWODbyID(woid)
	//fmt.Println("data",data)
	useraddworkouttpl.Execute(w, data)
}

//// END GetWOD Functions

//// START AddWOD Functions
// pageLoadAddWorkout - initial page load
func pageLoadAddWorkout(w http.ResponseWriter) {
	// default load todays wod if there is one for quick edits
	//wo := models.GetWODGuest()
	Edit = false
	adminaddwodtpl.Execute(w, nil)
}

// loadWOD - loads workout for selected date
func loadWODEdit(w http.ResponseWriter, r *http.Request) {
	wo := workouts.GetAddWODbydate(r.FormValue("date"))
	Edit = true
	admineditwodtpl.Execute(w, wo)
}

// postWOD - write workout to the database and reloads it to the page
func postWOD(w http.ResponseWriter, r *http.Request, uid string, admin bool) {
	wo := workouts.AddWOD(r, uid, false)
	if wo.Message == "" {
		fmt.Println("msg null")
		Edit = true
		if admin == true {
			admineditwodtpl.Execute(w, wo)
		} else {
			usereditworkouttpl.Execute(w, wo)
		}
	} else {
		fmt.Println("msg not null", wo.Message)
		Edit = false
		if admin == true {
			adminaddwodtpl.Execute(w, wo)
		} else {
			useraddworkouttpl.Execute(w, wo)
		}
	}
}

// editWOD - saves changes made to the USER CREATED workout and reloads it to the page
func editWOD(w http.ResponseWriter, r *http.Request, uid string, admin bool) {
	workouts.EditAddWOD(r, uid, true)
	Edit = true
	wo := workouts.GetAddWODbyID(r.FormValue("id"))
	if admin == true {
		admineditwodtpl.Execute(w, wo)
	} else {
		usereditworkouttpl.Execute(w, wo)
	}

}

//// END AddWOD Functions

//Not currently used
// adminEditWOD - saves changes made to the ADMIN CREATED workout and reloads it to the page
func adminEditWOD(w http.ResponseWriter, r *http.Request, uid string) {
	workouts.EditAddWOD(r, uid, true)
	Edit = true
	wo := workouts.GetAddWODbyID(r.FormValue("id"))
	admineditwodtpl.Execute(w, wo)
}

// adminPostWOD - write workout to the database and reloads it to the page
func adminPostWOD(w http.ResponseWriter, r *http.Request, uid string) {
	wo := workouts.AdminAddWOD(r, uid)
	if wo.Message == "" {
		Edit = true
		admineditwodtpl.Execute(w, wo)
	} else {
		Edit = false
		adminaddwodtpl.Execute(w, wo)
	}
}