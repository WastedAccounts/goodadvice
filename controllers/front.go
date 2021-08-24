package controllers

import (
	"encoding/json"
	"goodadvice/v1/controllers/auth"
	"goodadvice/v1/controllers/messaging"
	"goodadvice/v1/controllers/profile"
	"goodadvice/v1/models"
	"html/template"
	"io"
	"net/http"
)

// html template
var guestindexttpl = template.Must(template.ParseFiles("htmlpages/guestindex.html","htmlpages/templates/headerguest.html","htmlpages/templates/footerguest.html"))
var adminindextpl = template.Must(template.ParseFiles("htmlpages/adminindex.html","htmlpages/templates/header.html","htmlpages/templates/footer.html"))
var loggedinindextpl = template.Must(template.ParseFiles("htmlpages/index.html","htmlpages/templates/header.html","htmlpages/templates/footer.html"))
var abouttpl = template.Must(template.ParseFiles("htmlpages/about.html","htmlpages/templates/header.html","htmlpages/templates/footer.html"))
var aboutguesttpl = template.Must(template.ParseFiles("htmlpages/about.html","htmlpages/templates/headerguest.html","htmlpages/templates/footerguest.html"))
//var testtpl = template.Must(template.ParseFiles("htmlpages/templates/body.html","htmlpages/templates/header.html","htmlpages/templates/footer.html"))

func RegisterControllers() {
	//uc := newUserController()
	woc := newWorkoutController()
	lc := newLoginController()
	sc := newSignupController()
	awc := newAddWorkoutController()
	loc := newLogOutController()
	adm := newAdminController()
	pc := profile.NewProfileController()
	authc := auth.NewAuthController()
	msgc := messaging.NewMsgController()
	http.HandleFunc("/", index)
	http.HandleFunc("/about", about)
	//http.HandleFunc("/test", test)
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
	http.Handle("/profile/", *pc)
	http.Handle("/auth/", *authc)
	http.Handle("/messaging/", *msgc)

	//upc := profile.newUserProfileController()
	//http.Handle("/userprofile", *upc)
	//http.Handle("/userprofile/", *upc)
	//http.Handle("/assets/", http.StripPrefix("/assets/css/", http.FileServer(http.Dir("/assets/css/"))))
	//http.Handle("/assets", http.FileServer(http.Dir("assets/css/style.css")))
	//http.Handle("/users", *uc)
	//http.Handle("/users/", *uc)
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

func about(w http.ResponseWriter, r *http.Request) {
	active,c := models.ValidateSession(w,r)
	if c.Exists == true && c.Isadmin == true && active == true {
		abouttpl.Execute(w, nil)
	} else {
		aboutguesttpl.Execute(w, nil)
	}
}

//func test(w http.ResponseWriter, r *http.Request) {
//	testtpl.Execute(w, nil)
//	//testtpl.Execute(os.Stdout, nil)
//}

