package models

import (
	"database/sql"
	"fmt"
	"net/http"
)

func LogOut(w http.ResponseWriter, r *http.Request){
	c := validateCookie(w, r)
	db, err := sql.Open("mysql", DataSource)
	// validate session is LESS then 2 hours old
	deleteqs := fmt.Sprintf("update user_session set sessionstart = '1900-01-01 00:00:00' where userid = '%s'",c.Uid)
	delete, err := db.Query(deleteqs)
	fmt.Println(delete)
	if err != nil {
		panic(err.Error())
	}
}
