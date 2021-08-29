package messaging

import (
	"fmt"
	"goodadvice/v1/models"
	"goodadvice/v1/models/messaging"
	"html/template"
	"net/http"
	"regexp"
)

type msgController struct {
	msgIDPattern *regexp.Regexp
}

var suggestionboxtpl = template.Must(template.ParseFiles("htmlpages/messaging/suggestionbox.html", "htmlpages/templates/header.html", "htmlpages/templates/footer.html"))

func NewMsgController() *msgController {
	return &msgController{
		msgIDPattern: regexp.MustCompile(`^/messaging/suggestionbox/(\d+)/?`),
	}
}

// set cookies: https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/06.1.html
func (msgc msgController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Validate session
	userauth := models.ValidateSession(w, r)
	if userauth.IsActive == false {
		http.Redirect(w, r, "/login", 401)
	}
	if userauth.Exists == false {
		http.Redirect(w, r, "/login", 401)
		//} else if c.Isadmin != false {
		//	http.Redirect(w, r, "/login", 401)
	} else if r.URL.Path == "/messaging/suggestionbox" {
		switch r.Method {
		case http.MethodGet:
			msgc.pageLoadSuggestionbox(w, r)
		case http.MethodPost:
			messaging.SaveSuggestion(w, r, userauth.Uid)
		default:
			fmt.Println("status not implemented")
			w.WriteHeader(http.StatusNotImplemented)

		}
		//Template
		if r.URL.Path == "pathtohtmlpage" {
			switch r.Method {
			case http.MethodGet:
			case http.MethodPost:
			default:
				fmt.Println("status not implemented")
				w.WriteHeader(http.StatusNotImplemented)
			}
		}
	}
}

func (msgc *msgController) pageLoadSuggestionbox(w http.ResponseWriter, r *http.Request) {
	suggestionboxtpl.Execute(w, nil)
}
