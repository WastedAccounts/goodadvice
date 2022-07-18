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
	prIDPattern      *regexp.Regexp //not being used
}

// M - a map for passing multiple structs to htmltmpl
type M map[string]interface{}

// entry point from front.go
func NewProfileController() *profileController {
	return &profileController{
		profileIDPattern: regexp.MustCompile(`^/profile/(\d+)/?`),
		prIDPattern:      regexp.MustCompile(`^/profile/editpr(\d+)/?`),
	}
}

var userprofiletpl = template.Must(template.ParseFiles("htmlpages/profile/userprofile.html", "htmlpages/templates/header.html", "htmlpages/templates/footer.html"))
var aboutmetpl = template.Must(template.ParseFiles("htmlpages/profile/aboutme.html", "htmlpages/templates/header.html", "htmlpages/templates/footer.html"))
var goalstpl = template.Must(template.ParseFiles("htmlpages/profile/goals.html", "htmlpages/templates/header.html", "htmlpages/templates/footer.html"))
var personalrecordstpl = template.Must(template.ParseFiles("htmlpages/profile/personalrecords.html", "htmlpages/templates/header.html", "htmlpages/templates/footer.html"))
var editprtpl = template.Must(template.ParseFiles("htmlpages/profile/editpr.html", "htmlpages/templates/header.html", "htmlpages/templates/footer.html"))

// set cookies: https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/06.1.html
func (pc profileController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Validate session
	userauth := models.ValidateSession(w, r)
	if userauth.IsActive == false {
		http.Redirect(w, r, "/login", 401)
	}
	if userauth.Exists == false {
		http.Redirect(w, r, "/login", 401)
		//} else if c.Isadmin != false {
		//	http.Redirect(w, r, "/login", 401)
	} else {
		if r.RequestURI == "/profile/userprofile" {
			switch r.Method {
			case http.MethodGet:
				// initial page load of userprofile.html
				pc.pageLoaduserProfile(w, r, userauth.Uid)
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
				pc.pageLoadAboutMe(w, r, userauth.Uid)
			case http.MethodPost:
				pc.pageSaveAboutMe(w, r, userauth.Uid)
				pc.pageLoadAboutMe(w, r, userauth.Uid)
			default:
				fmt.Println("status not implemented")
				w.WriteHeader(http.StatusNotImplemented)
			}
		} else if r.RequestURI == "/profile/goals" {
			switch r.Method {
			case http.MethodGet:
				pc.pageLoadGoals(w, r, userauth.Uid)
			case http.MethodPost:
			default:
				fmt.Println("status not implemented")
				w.WriteHeader(http.StatusNotImplemented)
			}
		} else if r.URL.Path == "/profile/personalrecords" {
			switch r.Method {
			case http.MethodGet:
				pc.pageLoadPersonalRecords(w, r, userauth.Uid)
			case http.MethodPost:
				time := r.FormValue("minutes") + ":" + r.FormValue("seconds")
				add := profile.Addpr{
					Uid:          userauth.Uid,
					MovementName: r.FormValue("prddl"),
					Weight:       r.FormValue("weight"),
					Date:         r.FormValue("prdate"),
					Time:         time,
					Notes:        r.FormValue("notes"),
				}
				profile.SaveNewPR(add)
				pc.pageLoadPersonalRecords(w, r, userauth.Uid)
			default:
				fmt.Println("status not implemented")
				w.WriteHeader(http.StatusNotImplemented)
			}
		} else if r.URL.Path == "/profile/editpr" {
			switch r.Method {
			case http.MethodGet:
				u, _ := url.Parse(r.RequestURI)
				prid, _ := url.ParseQuery(u.RawQuery)
				pc.pageLoadEditpr(w, r, userauth.Uid, prid.Get("prid"))

			case http.MethodPost:
				pc.pageSaveEditpr(w, r, userauth.Uid)
				pc.pageLoadEditpr(w, r, userauth.Uid, r.PostFormValue("prid"))
			default:
				fmt.Println("status not implemented")
				w.WriteHeader(http.StatusNotImplemented)
			}
		}
	}
}

func (pc profileController) pageLoaduserProfile(w http.ResponseWriter, r *http.Request, id string) {
	// get values for pages load
	rec, up := profile.PageLoadUserProfile(id)
	userprofiletpl.Execute(w, M{
		//"mov": mov,
		"up":  up,
		"rec": rec,
	})
}

func (pc profileController) pageLoadAboutMe(w http.ResponseWriter, r *http.Request, id string) {
	// load PRs on page lade
	am := profile.LoadAboutMe(id)
	aboutmetpl.Execute(w, am) //(w, userprofile)
}

func (pc profileController) pageSaveAboutMe(w http.ResponseWriter, r *http.Request, id string) {
	profile.UpdateAboutMe(r, id)
}

func (pc profileController) pageLoadGoals(w http.ResponseWriter, r *http.Request, id string) {
	// get goals for page lade
	rec := profile.LoadPersonalRecords(id)
	goalstpl.Execute(w, rec) //(w, userprofile)
}

func (pc profileController) pageLoadPersonalRecords(w http.ResponseWriter, r *http.Request, id string) {
	// get PRs for page load
	rec := profile.LoadPersonalRecords(id)
	mov := profile.LoadMovements()
	personalrecordstpl.Execute(w, M{
		//"mov": mov,
		"mov": mov,
		"rec": rec,
	})
}

func (pc profileController) pageLoadEditpr(w http.ResponseWriter, r *http.Request, id string, prid string) {
	//load a single PR value for editing
	pr, prhist := profile.LoadSinglePR(id, prid)
	editprtpl.Execute(w, M{
		"pr":     pr,
		"prhist": prhist,
	})
}

func (pc profileController) pageSaveEditpr(w http.ResponseWriter, r *http.Request, id string) {
	//  a single PR value for editing
	profile.UpdateSinglePR(r, id)
}
