package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"mjs/v1/models"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type workoutController struct {
	workoutIDPattern *regexp.Regexp
}

type WorkoutPageLoad struct {
	WoID int `json:"woID"`
	WoName string `json:"woName"`
	WoStrength string `json:"woStrength"`
	WoPace string `json:"woPace"`
	WoConditioning string `json:"woConditioning"`
	WoDate string `json:"woDate"`
	//WoDOW string `json:"woDOW"`
	UsrID string `json:"usrID"`
	UsrNoteID int`json:"usrNoteID"`
	UsrName string `json:"usrName""`
	UsrNotes string `json:"usrNotes""`
}

type Cookie struct {
	Exists bool
	Uid string
	Sessionkey string
}

// html templates
var	wodtpl = template.Must(template.ParseFiles("htmlpages/wod.html"))
var	wodguesttpl = template.Must(template.ParseFiles("htmlpages/guestwod.html"))

// entry point from front.go
func newWorkoutController() *workoutController {
	return &workoutController{
		workoutIDPattern: regexp.MustCompile(`^/wod/(\d+)/?`),
	}
}

//ServeHTTP
// Entry point for the /wod page
// Comes in from front.go
func (woc workoutController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check for a cookie first
	//cookie, err := r.Cookie("WODerator")
	//// No cookie then get guest WOD page
	//if err != nil {
	//	// if not exist redirect to login page
	//	fmt.Println(" cookie is doesn't exist ")
	//	woc.GetWODGuest(w, r)
	//	//http.Redirect(w,r, "/login", http.StatusUnauthorized)
	//	//return nil, err
	//} else if err == nil {
	//	cookievalue := cookie.Value
	//	splitcookie := strings.Split(cookievalue, "/")
	//	//splitcookie[0] == userid
	//	//splitcookie[1] == sessionkey
	c := models.ValidateSession(w, r)
	if c.Exists == false {
		//http.Redirect(w,r, "/login", 401)
		//fmt.Println("workout controller cookievalue:",cookievalue)
		//fmt.Println("workout controller splitcookie[0]:",splitcookie[0])
		//fmt.Println("workout controller splitcookie[1]:",splitcookie[1])
		woc.GetWODGuest(w, r)
	} else {
		// At this point the user should be validated within two hour session time out
		// and a new cookie issued with new start time stamp
		uid := c.Uid//splitcookie[0]
		if r.URL.Path == "/wod" {
			switch r.Method {
			case http.MethodGet:
				if r.FormValue("date") != "" {
					woc.getWODbydate(w, r, r.FormValue("date"), uid)
				} else {
					woc.getWOD(w, r, uid)
				}
			case http.MethodPost:
				err := r.ParseForm()
				if err != nil {
					log.Fatalf("Failed to decode postFormByteSlice: %v", err)
				}
				//woc.postWODnotes(w, r, r.PostFormValue("notes"), r.PostFormValue("woid"), r.PostFormValue("uid"))
				woc.postWODnotes(r.PostFormValue("notes"), r.PostFormValue("woid"), r.PostFormValue("uid"))
			//case http.MethodPut:
			default:
				w.WriteHeader(http.StatusNotImplemented)
			}
		}
	}
}


// getWOD - displays WOD for the current date if there is one
func (woc *workoutController) getWOD(w http.ResponseWriter, r *http.Request, uid string) {
	wo, won, usr := models.GetWOD(uid)
	wpl := WorkoutPageLoad{
		 wo.ID,
		wo.Name,
		wo.Strength,
		wo.Pace,
		wo.Conditioning,
		wo.Date,
		//wo.DOW,
		strconv.Itoa(usr.ID), //<-- I should be getting this from somewhere else NOT from the notes
		won.ID,
		usr.UserName,
		 won.Notes,
	}
	fmt.Println(wpl)
	wodtpl.Execute(w, wpl)
}

// GetWODGuest for non auth'd users
func (woc *workoutController) GetWODGuest(w http.ResponseWriter, r *http.Request) {
	wo := models.GetWODGuest()
	wpl := WorkoutPageLoad{
		wo.ID,
		wo.Name,
		wo.Strength,
		wo.Pace,
		wo.Conditioning,
		wo.Date,
		//wo.DOW,
		"", //strconv.Itoa(usr.ID), //<-- I should be getting this from somewhere else NOT from the notes
		0, //won.ID,
		"",//usr.UserName,
		"",//won.Notes,
	}
	fmt.Println(wpl)
	wodguesttpl.Execute(w, wpl)
}

// getWODbydate - displays WOD for the current date if there is one
func (woc *workoutController) getWODbydate(w http.ResponseWriter, r *http.Request, d string, uid string) {
	wo, won, usr := models.GetWODbydate(d, uid)
	wpl := WorkoutPageLoad{
		wo.ID,
		wo.Name,
		wo.Strength,
		wo.Pace,
		wo.Conditioning,
		wo.Date,
		//wo.DOW,
		strconv.Itoa(usr.ID), //<-- I should be getting this from somewhere else NOT from the notes
		won.ID,
		usr.UserName,
		won.Notes,
	}
	splitdate := strings.Split(wpl.WoDate, "T")
	wpl.WoDate = splitdate[0]
	wodtpl.Execute(w, wpl)
}

// postWODnotes - get notes for user for WOD being loaded
func (woc *workoutController) postWODnotes(n string, woid string, uid string) {
	models.PostWODNotes(n, uid, woid)
	//http.Redirect(w, r, r.Header.Get("Referer"), 302)
}

//setCookie
// https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/06.1.html
func setCookie(w http.ResponseWriter, r *http.Request)  {
	//var iu IdentifiedUser
	//cun, _ := r.Cookie("username")
	//cid, _ := r.Cookie("ID")
	//fmt.Fprint(w, cun, cid)
	//
	//iu.ID = strconv.Itoa(cid)
	//iu.Name = strconv.Itoa(cun)
	//return iu
}


///////////////////////////////
// Below is code examples for the future
// createWOD - displays WOD for the current date if there is one
func (woc *workoutController) createWOD(w http.ResponseWriter, r *http.Request) {
	// load page without data
	addwodtpl.Execute(w, nil)
	//After editing update  WOD in database
}

// editWOD - displays WOD for the current date if there is one
func (woc *workoutController) editWOD(w http.ResponseWriter, r *http.Request) {
	//load current day's WOD for speed of editing
	//wo := models.GetWOD(uid string)
	//wo.Strength = strings.Replace(wo.Strength,"\n","<br>",-1)
	//wo.Conditioning = strings.Replace(wo.Conditioning,"\n","<br>",-1)
	//addwodtpl.Execute(w, wo)
	//After editing update  WOD in database
}

//works
func (woc *workoutController) post(w http.ResponseWriter, r *http.Request) {
	var exists bool
	u, err := woc.parseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not parse User object"))
		return
	}
	u, exists = models.CheckIfUserExists(u)
	if exists == false {
		u, err = models.AddUser(u)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Visitor already exists"))
		return
	}
	encodeResponseAsJSON(u, w)
}

//not in use
func (woc *workoutController) get(id int, w http.ResponseWriter) {
	u, err := models.GetUserByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	encodeResponseAsJSON(u, w)
}

//not in use
func (woc *workoutController) put(id int, w http.ResponseWriter, r *http.Request) {
	u, err := woc.parseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not parse User object"))
		return
	}
	if id != u.ID {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ID of submitted user must match ID in URL"))
		return
	}
	u, err = models.UpdateUser(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	encodeResponseAsJSON(u, w)
}

func (woc *workoutController) delete(id int, w http.ResponseWriter) {
	err := models.RemoveUserByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (woc *workoutController) parseRequest(r *http.Request) (models.User, error) {
	dec := json.NewDecoder(r.Body)
	var u models.User
	err := dec.Decode(&u)
	if err != nil {
		return models.User{}, err
	}
	return u, nil
}


//else {
//matches := woc.workoutIDPattern.FindStringSubmatch(r.URL.Path)
//if len(matches) == 0 {
//w.WriteHeader(http.StatusNotFound)
//}
//id, err := strconv.Atoi(matches[1])
//if err != nil {
//w.WriteHeader(http.StatusNotFound)
//}
//switch r.Method {
//case http.MethodGet:
//woc.get(id, w)
//case http.MethodPut:
//woc.put(id, w, r)
//case http.MethodDelete:
//woc.delete(id, w)
//default:
//w.WriteHeader(http.StatusNotImplemented)
//}
//}

