package messaging

import (
	"goodadvice/v1/datasource"
	"net/http"
)

func SaveSuggestion(w http.ResponseWriter, r *http.Request, uid string) {

	// write code to db
	_, err := datasource.DBconn.Exec("INSERT INTO suggestionbox (suggestion_subject,suggestion_msg,userid) VALUES (?,?,?)", r.PostFormValue("subject"), r.PostFormValue("suggestions"), uid)
	if err != nil {
		// If there is any issue with inserting into the database, return a 500 error
		panic(err.Error())
	}
}
