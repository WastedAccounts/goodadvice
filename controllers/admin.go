package controllers

import (
	"fmt"
	"goodadvice/v1/models"
	"html/template"
	"log"
	"net/http"
	"regexp"
)

type adminController struct {
	adminIDPattern *regexp.Regexp
}

type Movement struct {
	MovementType string
	Options string
}

type LoadMovements struct {
	something string
	// Movement type DDL values
	DDLoptions []string
}

// html templates
var admintpl = template.Must(template.ParseFiles("htmlpages/admin.html"))
var adminmovementstpl = template.Must(template.ParseFiles("htmlpages/adminmovements.html"))

// entry point from front.go
func newAdminController() *adminController {
	return &adminController{
		adminIDPattern: regexp.MustCompile(`^/wod/(\d+)/?`),
	}
}

// ServeHTTP - Entry point from front.go
func (ac adminController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Validate session
	c := models.ValidateSession(w, r)
	if c.Exists == false {
		http.Redirect(w, r, "/login", 401)
	} else if c.Isadmin == false {
		http.Redirect(w, r, "/login", 401)
	} else {
		if r.URL.Path == "/admin" {
			switch r.Method {
			case http.MethodGet:
				submit := r.FormValue("submit")
				if submit == "users" {
					pageLoadUsers()
				} else if submit == "movements" {
					pageLoadMovements(w,r)
				} else if submit == "workouts" {
					pageLoadWorkouts()
				} else {
					pageLoadAdmin(w)
				}
			case http.MethodPost:
				movements := r.FormValue("addmovement")
				err := r.ParseForm()
				if err != nil {
					log.Fatalf("Failed to decode postFormByteSlice: %v", err)
				}
				//movements := r.FormValue("movements")
				if movements == "addmovement" {
					saveMovement(w,r)
					pageLoadMovements(w,r)
				} else {
					fmt.Println("notmovements lol")
				}
			default:
				w.WriteHeader(http.StatusNotImplemented)
			}
		}
	}
}



// pageLoadAdmin
func pageLoadAdmin(w http.ResponseWriter) {
	// default load todays wod if there is one for quick edits
	//wo := models.GetWODGuest()
	//Edit = false -- used on addworkout to control edit vs new templates
	v := models.GetVersion()
	admintpl.Execute(w, v)
}

// pageLoadUsers - switch to users template do work there
func pageLoadUsers() {

}


// pageLoadMovements - switch to Movements template - This doesn't dynamically populate the DDL yet so I hard coded the page.
func pageLoadMovements(w http.ResponseWriter, r *http.Request) {
	var lm LoadMovements
	//All of this is broken, well, not broken, just doesn't work as I expected
	mt := models.GetMovementTypes()
	//var DDLoptions string
	//for _,x := range mt {
	//	////i := strconv.Itoa(x)
	//	//i := x.MovementType
	//	DDLoptions += fmt.Sprintf("<option value=\"" + strings.ToLower(x.MovementType) + "\">" + x.MovementType + "</option>\r")
	//	//fmt.Println(x)
	//	fmt.Println(DDLoptions)
	//}
	for _,x := range mt {
		lm.DDLoptions = append(lm.DDLoptions, x.MovementType)
	}
	//adminmovementstpl.Execute(w, DDLoptions)
	adminmovementstpl.Execute(w, nil)
}

// saveMovement - save new movement type from adminmovements.html
func saveMovement(w http.ResponseWriter, r *http.Request) {
	m := r.FormValue("movement")
	mt := r.FormValue("movementtypes")
	models.SaveMovement(m,mt)
}

// pageLoadWorkouts - switch to workouts template do work there
func pageLoadWorkouts() {

}
