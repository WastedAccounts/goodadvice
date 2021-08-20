package models

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"goodadvice/v1/models/messaging"
	"log"
	"net/http"
	"time"
)

type NewUser struct{
	User string
	Password string
	Firstname string
	Email string
	Date time.Time
}

func Signup(w http.ResponseWriter, r *http.Request) {
	var nu NewUser
	nu.Firstname = r.FormValue("firstname")
	nu.Password = r.FormValue("password")
	nu.User = r.FormValue("username")
	nu.Email = r.FormValue("email")

	// Salt and hash the password using the bcrypt algorithm
	// The second argument is the cost of hashing, which we arbitrarily set as 8 (this value can be more or less, depending on the computing power you wish to utilize)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(nu.Password), 8)
	// Next, insert the user values and hashed password into the database
	db, err := sql.Open("mysql", DataSource)
	insert, err := db.Exec("insert into users (username, firstname,lastlogindate,emailaddress,password,createdate) values (?,?,CURDATE(),?,?,CURDATE())",nu.User,nu.Firstname,nu.Email,string(hashedPassword))
	if err != nil {
		panic(err.Error())
	}
	// get new user id value so we can store it in a cookie
	newuid, err := insert.LastInsertId()
	if err != nil {
		// If there is any issue with inserting into the database, return a 500 error
		panic(err.Error())
	}
	// We reach this point if the credentials we correctly stored in the database, and the default status of 200 is sent back

	// now we send off a confirmation email and redirect to the confirmation page
	messaging.VerificationEmail(newuid,DataSource)

}

func CheckEmail(r *http.Request) bool {
	ef := false
	db, err := sql.Open("mysql", DataSource)
	//esq := fmt.Sprintf("select ID from users where emailaddress = '%s'", r.FormValue("email"))
	chkemail, err := db.Query("select ID from users where emailaddress = ?", r.FormValue("email"))
	if err != nil {
		// If there is any issue with inserting into the database, return a 500 error
		panic(err.Error())
	}
	var checkEmail int
	for chkemail.Next() {
		err := chkemail.Scan(&checkEmail)
		if err != nil {
			log.Fatal(err)
		}
	}
	chkemail.Close()
	if checkEmail != 0 {
		//if email already exists (not 0) return true
		ef = true
	}
	return ef
}

func CheckUsername(r *http.Request) bool {
	uf := false
	db, err := sql.Open("mysql", DataSource)
	usq := fmt.Sprintf("select ID from users where username = '%s'", r.FormValue("username"))
	chkusername, err := db.Query(usq)
	if err != nil {
		// If there is any issue with inserting into the database, return a 500 error
		panic(err.Error())
	}
	var checkUsername int
	for chkusername.Next() {
		err := chkusername.Scan(&checkUsername)
		if err != nil {
			log.Fatal(err)
		}
	}
	chkusername.Close()
	if checkUsername != 0 {
		//if username already exists (not 0) return true
		uf = true
	}
	return uf
}



