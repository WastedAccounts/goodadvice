package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"goodadvice/v1/models"
	"net/http"
)

// Port forward call for k8s
// kubectl port-forward mysql-5c9bd6bbb5-vq9vg 3306:3306 --address 0.0.0.0

// html template
var guesttpl = template.Must(template.ParseFiles("htmlpages/guestindex.html"))
var tpl = template.Must(template.ParseFiles("htmlpages/index.html"))

func RegisterControllers() {
	//uc := newUserController()
	woc := newWorkoutController()
	lc := newLoginController()
	sc := newSignupController()
	awc := newAddWorkoutController()
	loc := newLogOutController()
	http.HandleFunc("/", index)
	//http.HandleFunc("/resume", resume)
	//http.HandleFunc("/check", check)
	//http.Handle("/users", *uc)
	//http.Handle("/users/", *uc)
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
}

func encodeResponseAsJSON(data interface{}, w io.Writer) {
	enc := json.NewEncoder(w)
	enc.Encode(data)
	//return data
}


func index(w http.ResponseWriter, r *http.Request) {
	c := models.ValidateSession(w,r)
	if 	c.Exists == true {
		fmt.Fprint(w, "<h1>Welcome, this is good advice</h1>")
		tpl.Execute(w, nil)
	} else {
		fmt.Fprint(w, "<h1>Welcome, this is good advice</h1>")
		guesttpl.Execute(w, nil)
	}

}


//// addwod broken - trying to figure out how to display data
//func addwod(w http.ResponseWriter, r *http.Request) {
//	//template.ParseFiles("addworkout.html")
//	//fmt.Fprint(w, "<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n    <meta charset=\"UTF-8\">\n    <title>Add WOD</title>\n</head>\n<body>\n  <form action=\"/login\" method=\"post\">\n    Name:<input type=\"text\" name=\"username\">\n    <br />\n    Strength work:<input type=\"textbox\" aria-multiline=\"true\" name=\"username\">\n    <br />\n    Conditioning work:<input type=\"textbox\" aria-multiline=\"true\" name=\"password\">\n    <br />\n    testbox: <textarea name=\"textarea\" style=\"width:250px;height:150px;\"></textarea>\n    <br />\n    Pace:<input type=\"text\" name=\"username\">\n    <br />\n    Date:<input type=\"date\" data-date-inline-picker=\"true\" />\n    <input type=\"submit\" value=\"Login\">\n  </form>\n</body>\n</html>")
//	addwodtpl.Execute(w, nil)
//	//*woc
//}
//
//func check(w http.ResponseWriter, r *http.Request) {
//	fmt.Fprint(w, "<h1>Health Check</h1>")
//}
//
//func resume(w http.ResponseWriter, r *http.Request) {
//	fmt.Fprint(w, "<h1>Resume</h1>")
	//restpl.Execute(w, nil)
//}