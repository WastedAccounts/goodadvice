package controllers

import (
	"encoding/json"
	"goodadvice/v1/controllers/auth"
	"goodadvice/v1/controllers/messaging"
	"goodadvice/v1/controllers/profile"
	"goodadvice/v1/controllers/workouts"
	"goodadvice/v1/models"
	"html/template"
	"io"
	"net/http"
)

// html template
var guestindexttpl = template.Must(template.ParseFiles("htmlpages/guestindex.html", "htmlpages/templates/headerguest.html", "htmlpages/templates/footerguest.html"))
var adminindextpl = template.Must(template.ParseFiles("htmlpages/adminindex.html", "htmlpages/templates/header.html", "htmlpages/templates/footer.html"))
var userindextpl = template.Must(template.ParseFiles("htmlpages/userindex.html", "htmlpages/templates/header.html", "htmlpages/templates/footer.html"))
var userabouttpl = template.Must(template.ParseFiles("htmlpages/about.html", "htmlpages/templates/header.html", "htmlpages/templates/footer.html"))
var guestabouttpl = template.Must(template.ParseFiles("htmlpages/about.html", "htmlpages/templates/headerguest.html", "htmlpages/templates/footerguest.html"))

func RegisterControllers() {
	//uc := newUserController()

	lc := newLoginController()
	sc := newSignupController()
	loc := newLogOutController()
	adm := newAdminController()
	pc := profile.NewProfileController()
	authc := auth.NewAuthController()
	msgc := messaging.NewMsgController()
	woc := workouts.NewWorkoutController()
	share := workouts.NewShareController()
	http.HandleFunc("/", index)
	http.HandleFunc("/about", about)
	http.Handle("/admin", *adm)
	http.Handle("/admin/", *adm)
	http.Handle("/login", *lc)
	http.Handle("/login/", *lc)
	http.Handle("/signup", *sc)
	http.Handle("/signup/", *sc)
	http.Handle("/logout", *loc)
	http.Handle("/logout/", *loc)
	http.Handle("/profile/", *pc)
	http.Handle("/auth/", *authc)
	http.Handle("/messaging/", *msgc)
	http.Handle("/workouts/", *woc)
	http.Handle("/canyoubeatme", *share)

	//// Test page
	//http.HandleFunc("/test", test)

	//// old stuff
	//upc := profile.newUserProfileController()
	//awc := workouts.NewAddWorkoutController()
	//http.Handle("/userprofile", *upc)
	//http.Handle("/userprofile/", *upc)
	//http.Handle("/workouts/", *awc)
	//http.Handle("/assets/", http.StripPrefix("/assets/css/", http.FileServer(http.Dir("/assets/css/"))))
	//http.Handle("/assets", http.FileServer(http.Dir("assets/css/style.css")))
	//http.Handle("/users", *uc)
	//http.Handle("/users/", *uc)
	//http.Handle("/wod/", *woc)
	//http.Handle("/addwod", *awc)
	//http.Handle("/addwod/", *awc)
	//http.Handle("/wod", *woc)
	//http.Handle("/wod/", *woc)
	//fs := http.FileServer(http.Dir("/assets/css/"))
	//http.Handle("/assets/css/", http.StripPrefix("/assets/css/", fs))
	// Handle client's requests for CSS
	//http.Handle("./css/", http.FileServer(http.Dir("./css/")))
}

// index - loads index.html based on user roles
func index(w http.ResponseWriter, r *http.Request) {
	userauth := models.ValidateSession(w, r)
	if userauth.Exists == true && userauth.IsAdmin == true && userauth.IsActive == true {
		//fmt.Fprint(w, "<h1 class='header'>Welcome, this is good advice</h1>")
		adminindextpl.Execute(w, nil)
	} else if userauth.Exists == true && userauth.IsAdmin != true && userauth.IsActive == true {
		//fmt.Fprint(w, "<h1 class='header'>Welcome, this is good advice</h1>")
		userindextpl.Execute(w, nil)
	} else {
		//fmt.Fprint(w, "<h1 class='header'>Welcome, this is good advice</h1>")
		guestindexttpl.Execute(w, nil)
	}
}

// about - loads about.html based on user roles
func about(w http.ResponseWriter, r *http.Request) {
	userauth := models.ValidateSession(w, r)
	if userauth.Exists == true && userauth.IsAdmin == true && userauth.IsActive == true {
		userabouttpl.Execute(w, nil)
	} else {
		guestabouttpl.Execute(w, nil)
	}
}

// encodeResponseAsJSON - Not is use not sure what it did
func encodeResponseAsJSON(data interface{}, w io.Writer) {
	enc := json.NewEncoder(w)
	enc.Encode(data)
	//return data
}

//// Test page stuff - uncomment to enable - access at //url.com/test
//func RegisterControllersTest() {
//http.HandleFunc("/test", test)
//}
//
//var testtpl = template.Must(template.ParseFiles("htmlpages/templates/test.html")) //,"htmlpages/templates/header.html","htmlpages/templates/footer.html"))
//
//func test(w http.ResponseWriter, r *http.Request) {
//		var check Check
//	check.Checked = ""
//	test := r.PostFormValue("wodcb")
//	fmt.Println("wodcb",test)
//	testtpl.Execute(w, check)
//}
