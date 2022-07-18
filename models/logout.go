package models

import (
	"goodadvice/v1/datasource"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func LogOut(w http.ResponseWriter, r *http.Request) {
	c := ValidateCookie(w, r)

	// update session so it's invalid by setting start date to the turn of the 20th century
	_, err := datasource.DBconn.Exec("UPDATE user_session SET sessionstart = '1900-01-01 00:00:00' WHERE userid = ?", c.Uid)

	if err != nil {
		panic(err.Error())
	}
	//_, err = delete.LastInsertId()
	if err != nil {
		panic(err.Error())
	}
}
