package controllers

import (
	"encoding/json"
	"goodadvice/v1/controllers/auth"
	"goodadvice/v1/controllers/profile"
	"goodadvice/v1/models"
	"html/template"
	"io"
	"net/http"
)

// html template
var guestindexttpl = template.Must(template.ParseFiles("htmlpages/guestindex.html"))
var adminindextpl = template.Must(template.ParseFiles("htmlpages/adminindex.html"))
var loggedinindextpl = template.Must(template.ParseFiles("htmlpages/index.html"))

func RegisterControllers() {
	//uc := newUserController()
	woc := newWorkoutController()
	lc := newLoginController()
	sc := newSignupController()
	awc := newAddWorkoutController()
	loc := newLogOutController()
	adm := newAdminController()
	upc := newUserProfileController()
	pc := profile.NewProfileController()
	authc := auth.NewAuthController()
	http.HandleFunc("/", index)
	//http.Handle("/assets/", http.StripPrefix("/assets/css/", http.FileServer(http.Dir("/assets/css/"))))
	//http.Handle("/assets", http.FileServer(http.Dir("assets/css/style.css")))
	//http.Handle("/users", *uc)
	//http.Handle("/users/", *uc)
	http.Handle("/admin", *adm)
	http.Handle("/admin/", *adm)
	http.Handle("/wod", *woc)
	http.Handle("/wod/", *woc)
	http.Handle("/login", *lc)
	http.Handle("/login/", *lc)
	http.Handle("/signup", *sc)
	http.Handle("/signup/", *sc)
	http.Handle("/addwod", *awc)
	http.Handle("/addwod/", *awc)
	http.Handle("/logout", *loc)
	http.Handle("/logout/", *loc)
	http.Handle("/userprofile", *upc)
	http.Handle("/userprofile/", *upc)
	http.Handle("/profile/", *pc)
	http.Handle("/auth/", *authc)
	//fs := http.FileServer(http.Dir("/assets/css/"))
	//http.Handle("/assets/css/", http.StripPrefix("/assets/css/", fs))
	// Handle client's requests for CSS
	//http.Handle("./css/", http.FileServer(http.Dir("./css/")))
}



func encodeResponseAsJSON(data interface{}, w io.Writer) {
	enc := json.NewEncoder(w)
	enc.Encode(data)
	//return data
}


func index(w http.ResponseWriter, r *http.Request) {
	active,c := models.ValidateSession(w,r)
	if c.Exists == true && c.Isadmin == true && active == true {
		//fmt.Fprint(w, "<h1 class='header'>Welcome, this is good advice</h1>")
		adminindextpl.Execute(w, nil)
	} else if c.Exists == true && c.Isadmin != true && active == true  {
		//fmt.Fprint(w, "<h1 class='header'>Welcome, this is good advice</h1>")
		loggedinindextpl.Execute(w, nil)
	} else {
		//fmt.Fprint(w, "<h1 class='header'>Welcome, this is good advice</h1>")
		guestindexttpl.Execute(w, nil)
	}

}

