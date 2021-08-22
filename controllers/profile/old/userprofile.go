package old

import (
	"fmt"
	"goodadvice/v1/models"
	"goodadvice/v1/models/profile"
	"log"
	"net/http"
	"regexp"
)


type userProfileController struct {
	userProfileIDPattern *regexp.Regexp
}
// entry point from front.go
func newUserProfileController() *userProfileController {
	return &userProfileController{
		userProfileIDPattern: regexp.MustCompile(`^/userprofile/(\d+)/?`),
	}
}

//var	userprofiletpl = template.Must(template.ParseFiles("htmlpages/profile/userprofile.html"))

// set cookies: https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/06.1.html
func (upc userProfileController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
		switch r.Method  {
		case http.MethodGet:
			submit := r.FormValue("submit")
			err := r.ParseForm()
			if err != nil {
				log.Fatalf("Failed to decode getFormByteSlice: %v", err)
			}
			if submit == "" {
				upc.pageLoad(w, r, c.Uid)
			} else if submit == "Add Record" {
				add := profile.Addpr{
					Uid:          c.Uid,
					MovementName: r.FormValue("prddl"),
					PRvalue:     r.FormValue("prnew"),
					Date:         r.FormValue("prdate"),
				}
				profile.AddRecord(add)
				upc.pageLoad(w, r, c.Uid)
			}
		case http.MethodPost:
			if models.Login(w, r) == false {
			} else {
			}
			http.Redirect(w,r, "/", 302)
		default:
			fmt.Println("status not implemented")
			w.WriteHeader(http.StatusNotImplemented)
		}
	}
}

func (upc userProfileController) pageLoad(w http.ResponseWriter, r *http.Request, id string) {
	// load PRs on page lade
	//userprofile := profile.LoadAboutMe(id) // add a new struct for all values from all structs you want to use on the page.
	//profile2.userprofiletpl.Execute(w, userprofile)
}
