package profile

import (
	"fmt"
	"goodadvice/v1/models"
	"goodadvice/v1/models/profile"
	"html/template"
	"net/http"
	"net/url"
	"regexp"
)

// profileController - Not sure I still need this.
type profileController struct {
	profileIDPattern *regexp.Regexp // Not sure if I need this now that I can parse URLs better.
	prIDPattern *regexp.Regexp //not being used
}

// Userprofile - not in use - a struct to combine structs from model.profile
type userprofile struct {
	//Userinfo struct data
	Name string
	Birthday string
	Weight string
	Sex string
	About string
	Age int
	// Records struct data
	Movement []string
	PrId []string
	Record []string
	Movements []string
	Date []string
	Currentdate string
}

// LoadUserProfileDate - not in use
type LoadUserProfileDate struct {
	Movements profile.Movements
	Record []profile.Records
	UserProfile profile.Userinfo
}

// M - a map for passing multiple structs to htmltmpl
type M map[string]interface{}

// entry point from front.go
func NewProfileController() *profileController {
	return &profileController{
		profileIDPattern: regexp.MustCompile(`^/profile/(\d+)/?`),
		prIDPattern: regexp.MustCompile(`^/profile/editpr(\d+)/?`),
	}
}

var	userprofiletpl = template.Must(template.ParseFiles("htmlpages/profile/userprofile.html"))
var	aboutmetpl = template.Must(template.ParseFiles("htmlpages/profile/aboutme.html"))
var	goalstpl = template.Must(template.ParseFiles("htmlpages/profile/goals.html"))
var	personalrecordstpl = template.Must(template.ParseFiles("htmlpages/profile/personalrecords.html"))
var	editprtpl = template.Must(template.ParseFiles("htmlpages/profile/editpr.html"))

// set cookies: https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/06.1.html
func (pc profileController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Validate session
	active,c := models.ValidateSession(w, r)
	if active == false {
		http.Redirect(w, r, "/login", 401)
	}
	if c.Exists == false {
		http.Redirect(w, r, "/login", 401)
		//} else if c.Isadmin != false {
		//	http.Redirect(w, r, "/login", 401)
	} else {
		if r.RequestURI == "/profile/userprofile" {
			switch r.Method {
			case http.MethodGet:
				// initial page load of userprofile.html
				pc.pageLoaduserProfile(w,r,c.Uid)
			case http.MethodPost:
				//if models.Login(w, r) == false {
				//} else {
				//}
				//http.Redirect(w, r, "/", 302)
			default:
				fmt.Println("status not implemented")
				w.WriteHeader(http.StatusNotImplemented)
			}
		} else if r.RequestURI == "/profile/aboutme" {
			switch r.Method {
			case http.MethodGet:
				pc.pageLoadAboutMe(w, r, c.Uid)
			case http.MethodPost:
			default:
				fmt.Println("status not implemented")
				w.WriteHeader(http.StatusNotImplemented)
			}
		} else if r.RequestURI == "/profile/goals" {
			switch r.Method {
			case http.MethodGet:
				pc.pageLoadGoals(w, r, c.Uid)
			case http.MethodPost:
			default:
				fmt.Println("status not implemented")
				w.WriteHeader(http.StatusNotImplemented)
			}
		} else if r.URL.Path == "/profile/personalrecords" {
			switch r.Method {
			case http.MethodGet:
				pc.pageLoadPersonalRecords(w, r, c.Uid)
			case http.MethodPost:
				time := r.FormValue("minutes") + ":" + r.FormValue("seconds")
				add := profile.Addpr{
					Uid:          c.Uid,
					MovementName: r.FormValue("prddl"),
					Weight:     r.FormValue("prnew"),
					Date:         r.FormValue("prdate"),
					Time: time,
					Notes: r.FormValue("prnotes"),
				}
				profile.SaveNewPR(add)
				pc.pageLoadPersonalRecords(w, r, c.Uid)
			default:
				fmt.Println("status not implemented")
				w.WriteHeader(http.StatusNotImplemented)
			}
		} else if r.URL.Path == "/profile/editpr" {
			switch r.Method {
			case http.MethodGet:
				u,_ := url.Parse(r.RequestURI)
				prid,_ :=  url.ParseQuery(u.RawQuery)
				pc.pageLoadEditpr(w, r, c.Uid, prid.Get("prid"))

			case http.MethodPost:
				pc.pageSaveEditpr(w, r, c.Uid, )
				pc.pageLoadEditpr(w, r, c.Uid, r.PostFormValue("prid"))
			default:
				fmt.Println("status not implemented")
				w.WriteHeader(http.StatusNotImplemented)
			}
		}
	}
}

func (pc profileController) pageLoaduserProfile(w http.ResponseWriter, r *http.Request, id string) {
	//var up userprofile

	// get values for pages load
	//https://stackoverflow.com/questions/50080640/how-to-pass-multiple-variables-to-go-html-template
	rec,up := profile.PageLoadUserProfile(id)

	userprofiletpl.Execute(w,M{
		//"mov": mov,
		"up": up,
		"rec": rec,
	})
}

func (pc profileController) pageLoadAboutMe(w http.ResponseWriter, r *http.Request, id string) {
	// load PRs on page lade
	am := profile.LoadAboutMe(id)
	aboutmetpl.Execute(w,am) //(w, userprofile)
}

func (pc profileController) pageLoadGoals(w http.ResponseWriter, r *http.Request, id string) {
	// get goals for page lade
	rec := profile.LoadPersonalRecords(id)
	goalstpl.Execute(w,rec) //(w, userprofile)
}

func (pc profileController) pageLoadPersonalRecords(w http.ResponseWriter, r *http.Request, id string) {
	// get PRs for page load
	rec := profile.LoadPersonalRecords(id)
	mov := profile.LoadMovements()
	personalrecordstpl.Execute(w,M{
		//"mov": mov,
		"mov": mov,
		"rec": rec,
	})
}

func (pc profileController) pageLoadEditpr(w http.ResponseWriter, r *http.Request, id string, prid string) {
	//load a single PR value for editing
	pr := profile.LoadSinglePR(id, prid)
	editprtpl.Execute(w,pr)
}

func (pc profileController) pageSaveEditpr(w http.ResponseWriter, r *http.Request, id string) {
	//  a single PR value for editing
	profile.UpdateSinglePR(r, id)
}
