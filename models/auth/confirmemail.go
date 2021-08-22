package auth

import (
	"database/sql"
	"goodadvice/v1/datasource"
	"goodadvice/v1/models"
	"log"
	"net/http"
	"time"
)

// ConfirmEmail - validate code against DB value and timestamp and then remove from database
func ConfirmEmail(w http.ResponseWriter, r *http.Request) (bool,string) {
	//vars
	var expires time.Time
	var dbcode string
	var msg string
	c := models.ValidateCookie(w, r)
	webcode := r.FormValue("code")
	expiredby := time.Now().UTC()

	// open db conn
	db, err := sql.Open("mysql", datasource.DataSource)
	defer db.Close()
	// pulls code and expire time from db to validate
	checkCode, err := db.Query("SELECT verification_code,expires FROM email_verification WHERE userid = ?;", c.Uid)//(checkSessionAgeqs)
	if err != nil {
		panic(err.Error())
	}
	for checkCode.Next() {
		err := checkCode.Scan(&dbcode,&expires)
		if err != nil {
			log.Fatal(err)
		}
	}

	// now validate code and expiration
	if dbcode != webcode {
			msg = "Incorrect Confirmation Code"
			return false,msg
	}
	if expiredby.After(expires) {
		msg = "Expired Confirmation Code"
		return false,msg
	}
	if dbcode == webcode && expiredby.Before(expires) {
		msg = "Success"
		//activate user UPDATE users SET isactive = 1 WHERE ID = ?;
		// Get the existing entry present in the database for the given username
		db, err := sql.Open("mysql", datasource.DataSource)
		activateUser, err := db.Exec("UPDATE users SET isactive = 1 WHERE ID = ?;",c.Uid)
		if err != nil {
			panic(err.Error())
		}
		// get new user id value so we can store it in a cookie
		_, err = activateUser.LastInsertId()
		if err != nil {
			// If there is any issue with inserting into the database, return a 500 error
			panic(err.Error())
		}

		return true,msg
	} else {
		msg = "Something else is wrong"
		return false,msg
	}

}