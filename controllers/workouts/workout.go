package workouts

import (
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

type shareController struct {
	shareIDPattern *regexp.Regexp
}

// Used to control NEW vs EDIT templates
var Edit bool

// html templates
var adminaddwodtpl = template.Must(template.ParseFiles("htmlpages/workouts/adminaddworkout.html", "htmlpages/templates/header.html", "htmlpages/templates/footer.html"))
var admineditwodtpl = template.Must(template.ParseFiles("htmlpages/workouts/admineditworkout.html", "htmlpages/templates/header.html", "htmlpages/templates/footer.html"))
var userwodtpl = template.Must(template.ParseFiles("htmlpages/workouts/userwod.html", "htmlpages/templates/header.html", "htmlpages/templates/footer.html"))
var guestwodtpl = template.Must(template.ParseFiles("htmlpages/workouts/guestwod.html", "htmlpages/templates/headerguest.html", "htmlpages/templates/footerguest.html"))
var useraddworkouttpl = template.Must(template.ParseFiles("htmlpages/workouts/useraddworkout.html", "htmlpages/templates/header.html", "htmlpages/templates/footer.html"))
var usereditworkouttpl = template.Must(template.ParseFiles("htmlpages/workouts/usereditworkout.html", "htmlpages/templates/header.html", "htmlpages/templates/footer.html"))
var shareworkouttpl = template.Must(template.ParseFiles("htmlpages/workouts/shareworkout.html", "htmlpages/templates/header.html", "htmlpages/templates/footer.html"))

// M - a map for passing multiple structs to htmltmpl
type M map[string]interface{}

// entry point from front.go for /workouts
func NewWorkoutController() *workoutController {
	return &workoutController{
		workoutIDPattern: regexp.MustCompile(`^/workouts/wod/(\d+)/?`),
	}
}

// entry point from front.go for /canyoubeatme
func NewShareController() *shareController {
	return &shareController{
		shareIDPattern: regexp.MustCompile(`^/workouts/wod/(\d+)/?`),
	}
}

//// START GetWOD Functions
// ServeHTTP - Entry point for the /wod page - Comes in from front.go
func (woc workoutController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check for a cookie first
	userauth := models.ValidateSession(w, r)
	// If there is no cookie found redirect to guest view
	if userauth.Exists == false || userauth.IsActive == false {
		woc.getWOD(w, r, userauth, true)
	} else {
		// At this point the user should be validated within 48 hour session time out
		// and a new cookie issued with new session start time stamp
		if r.URL.Path == "/workouts/wod" {
			switch r.Method {
			case http.MethodGet:
				u, _ := url.Parse(r.RequestURI)
				woid, _ := url.ParseQuery(u.RawQuery)
				if r.FormValue("date") != "" {
					// If a date is selected load workout from that date
					woc.getWODbydate(w, r.FormValue("date"), userauth)
				} else if r.FormValue("random") == "Random" {
					// Get random workout from Random button click
					woc.randomWorkout(w, userauth)
				} else if woid.Get("woid") == "0" {
				//get daily wod
				woc.getWOD(w, r, userauth, true)
				} else {
					// if no date is selected load today's workout
					woc.getWOD(w, r, userauth, false)
				}
			case http.MethodPost:
				err := r.ParseForm()
				if err != nil {
					log.Fatalf("Failed to decode postFormByteSlice: %v", err)
				}
				woc.SaveWorkoutResults(w, r)
				woc.getWOD(w, r, userauth, false)
			case http.MethodPut:
			default:
				w.WriteHeader(http.StatusNotImplemented)
			}
		} else if r.URL.Path == "/workouts/addwod" {
			switch r.Method {
			case http.MethodGet:
				if r.FormValue("date") == "" {
					pageLoadAddWorkout(w, true)
				} else {
					loadWODEdit(w, r, true, "")
				}
			case http.MethodPost:
				err := r.ParseForm()
				if err != nil {
					log.Fatalf("Failed to decode postFormByteSlice: %v", err)
				}
				if Edit == true {
					editWOD(w, r, userauth.Uid, true)
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
				if woid.Get("woid") == "new" {
					// Ready to create a new user workout
					pageLoadAddWorkout(w, false)
				} else if woid.Get("woid") != "" {
					// Call for WOD by woid
					woc.customizeWOD(w, r, woid.Get("woid"), userauth, "")
				} else {
					// Get User WOD by date selected
					loadWODEdit(w, r, false, userauth.Uid)
				}
			case http.MethodPost:
				err := r.ParseForm()
				if err != nil {
					log.Fatalf("Failed to decode postFormByteSlice: %v", err)
				}
				if Edit == true {
					editWOD(w, r, userauth.Uid, false)
				} else {
					postWOD(w, r, userauth.Uid, false)
				}
			default:
				w.WriteHeader(http.StatusNotImplemented)
			}
		} else if r.URL.Path == "/workouts/share" {
			switch r.Method {
			case http.MethodGet:
				u, _ := url.Parse(r.RequestURI)
				woid, _ := url.ParseQuery(u.RawQuery)
				woc.shareWOD(w, r, woid.Get("woid"))
			default:
				w.WriteHeader(http.StatusNotImplemented)
			}
		}
	}
}

func (share shareController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check for a cookie first
	userauth := models.ValidateSession(w, r)
	u, _ := url.Parse(r.RequestURI)
	woid, _ := url.ParseQuery(u.RawQuery)
	// If there is no cookie found redirect to guest view
	if userauth.Exists == false || userauth.IsActive == false {
		if r.URL.Path == "/canyoubeatme" {
			switch r.Method {
			case http.MethodGet:
				canyoubeatme(w, r, woid.Get("woid"), userauth)
			default:
				w.WriteHeader(http.StatusNotImplemented)
			}
		}
	} else {
		if r.URL.Path == "/canyoubeatme" {
			switch r.Method {
			case http.MethodGet:
				canyoubeatme(w, r, woid.Get("woid"), userauth)
			default:
				w.WriteHeader(http.StatusNotImplemented)
			}
		}
	}
}

// canyoubeatme - loads shared workouts for logged in users and guests
func canyoubeatme(w http.ResponseWriter, r *http.Request, woid string, userauth models.UserAuth) {

	wo, won, usr := workouts.GetWODbyID(woid, userauth.Uid)
	data := M{
		"wo":  wo,
		"won": won,
		"usr": usr,
	}
	if userauth.Exists == true && userauth.IsActive == true {
		// If Admin
		userwodtpl.Execute(w, data)
	} else {
		// If Guest
		guestwodtpl.Execute(w, data)
	}
}

func (woc *workoutController) shareWOD(w http.ResponseWriter, r *http.Request, Woid string) {
	// Create share WOD page with url to share
	data := map[string]interface{}{
		"woid": Woid,
	}
	shareworkouttpl.Execute(w, data)
}

// getWOD - Gets WOD and User data if they're logged in and depending on Role loads the correct template
func (woc *workoutController) getWOD(w http.ResponseWriter, r *http.Request, userauth models.UserAuth, dailywod bool) {
	// vars
	var uid string

	// daily WOD if they want that specifically.
	if dailywod == true {
		uid = ""
	} else {
		uid = userauth.Uid
	}

	// Call GetWOD to get structs for page load
	wo, won, usr := workouts.GetWOD(uid, r)

	// Setup link to Daily WOD if we loaded a different workout
	if wo.WODworkout == "Y" {
		wo.Linkhidden = "hidden"
	}
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

	// Setup link to Daily WOD if we loaded a different workout
	if wo.WODworkout == "Y" {
		wo.Linkhidden = "hidden"
	}
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
	wo, won, usr := workouts.GetWODbyID(woid, userauth.Uid)
	data := M{
		"wo":  wo,
		"won": won,
		"usr": usr,
	}
	if wo.WODworkout == "N" && wo.CreatedBy == userauth.Uid {
		// if user loaded their own work we'll send them to edit
		Edit = true
		usereditworkouttpl.Execute(w, data)
	} else {
		Edit = false
		// If a user loaded a Daily WOD we'll copy it and send them to create a new one for themself
		useraddworkouttpl.Execute(w, data)
	}
}

//// END GetWOD Functions

//// START AddWOD Functions
// pageLoadAddWorkout - initial page load
func pageLoadAddWorkout(w http.ResponseWriter, admin bool) {
	// load a blank create workout page for creating daily wods
	Edit = false
	if admin == true {
		adminaddwodtpl.Execute(w, nil)
	} else {
		useraddworkouttpl.Execute(w, nil)
	}
}

// loadWOD - loads workout for selected date
func loadWODEdit(w http.ResponseWriter, r *http.Request, admin bool, uid string) {
	wo := workouts.GetAddWODbydate(r.FormValue("date"), uid)
	data := M{
		"wo":  wo,
	}
	Edit = true
	if admin == true {
		admineditwodtpl.Execute(w, data)
	} else {
		usereditworkouttpl.Execute(w, data)
	}
}

// postWOD - write workout to the database and reloads it to the page
func postWOD(w http.ResponseWriter, r *http.Request, uid string, admin bool) {
	wo := workouts.AddWOD(r, uid, false)
	data := M{
		"wo":  wo,
	}
	if wo.Message == "" {
		//fmt.Println("msg null")
		Edit = true
		if admin == true {
			admineditwodtpl.Execute(w, data)
		} else {
			usereditworkouttpl.Execute(w, data)
		}
	} else {
		//fmt.Println("msg not null", wo.Message)
		Edit = false
		if admin == true {
			adminaddwodtpl.Execute(w, data)
		} else {
			useraddworkouttpl.Execute(w, data)
		}
	}
}

// editWOD - saves changes made to the USER CREATED workout and reloads it to the page
func editWOD(w http.ResponseWriter, r *http.Request, uid string, admin bool) {
	workouts.EditAddWOD(r, uid, true)
	Edit = true
	wo := workouts.GetAddWODbyID(r.FormValue("id"))
	data := M{
		"wo":  wo,
	}
	if admin == true {
		admineditwodtpl.Execute(w, data)
	} else {
		usereditworkouttpl.Execute(w, data)
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
