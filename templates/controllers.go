package templates

import (
	"fmt"
	"goodadvice/v1/models"
	"html/template"
	"net/http"
	"regexp"
)

type msgController struct {
	msgIDPattern *regexp.Regexp
}

// html template
var suggestionboxtpl = template.Must(template.ParseFiles("htmlpages/messaging/suggestionbox.html"))

// new controller
func NewMsgController() *msgController {
	return &msgController{
		msgIDPattern: regexp.MustCompile(`^/messaging/suggestionbox/(\d+)/?`),
	}
}

// set cookies: https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/06.1.html
func (msgc msgController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Validate session
	active, c := models.ValidateSession(w, r)
	if active == false {
		http.Redirect(w, r, "/login", 401)
	}
	if c.Exists == false {
		http.Redirect(w, r, "/login", 401)
	} else if r.URL.Path == "pathtohtmlpage" {
		switch r.Method {
		case http.MethodGet:
			msgc.pageLoadSuggestionbox(w, r)
		case http.MethodPost:
		default:
			fmt.Println("status not implemented")
			w.WriteHeader(http.StatusNotImplemented)
		}
	}
}

func (msgc *msgController) pageLoadSuggestionbox(w http.ResponseWriter, r *http.Request) {
	suggestionboxtpl.Execute(w, nil)
}
