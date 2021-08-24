package models

import (
	"database/sql"
	"goodadvice/v1/datasource"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func LogOut(w http.ResponseWriter, r *http.Request){
	c := ValidateCookie(w, r)
	// open DB conn
	db, err := sql.Open("mysql", datasource.DataSource)
	defer db.Close()

	// update session so it's invalid by setting start date to the turn of the 20th century
	delete, err := db.Exec("update user_session set sessionstart = '1900-01-01 00:00:00' where userid = ?",c.Uid)
	//fmt.Println(delete)
	if err != nil {
		panic(err.Error())
	}
	_,err = delete.LastInsertId()
	if err != nil {
		panic(err.Error())
	}
}
