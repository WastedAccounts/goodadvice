package models

import (
	"database/sql"
	"fmt"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Credentials struct {
	Id string
	Username string
	Password string
}

type Cookie struct {
	Exists bool
	Uid string
	Sessionkey string
	Isadmin bool
}

type Authenticated struct {
	Authenticated bool
	Active bool
}

// https://www.sohamkamani.com/blog/2018/02/25/golang-password-authentication-and-storage/#implementing-user-login
func Login(w http.ResponseWriter, r *http.Request) bool {
	var success bool
	var isactive, password, id string
	// Parse and decode the request body into a new `Credentials` instance
	creds := Credentials{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}
	// Get the existing entry present in the database for the given username
	db, err := sql.Open("mysql", DataSource)
	if err != nil {
		panic(err.Error())
	}
	GetCreds, err := db.Query("select ID, password,isactive from users where username = ?", creds.Username) //("select password,isactive from users where username='?'", creds.Username)
	for GetCreds.Next() {
		err := GetCreds.Scan(&id,&password,&isactive)
		if err != nil {
			log.Fatal(err)
		}
	}
	// Check if user is active, if not fail login immediately
	if isactive == "0"{
		success = false
		return success
	}
	// Compare the stored hashed password, with the hashed version of the password that was received
	if err = bcrypt.CompareHashAndPassword([]byte(password), []byte(creds.Password)); err != nil {
		// If the two passwords don't match, return a 401 status
		w.WriteHeader(http.StatusUnauthorized)
		success = false
		return success
	}
	// Capture login date and IP to login_history table
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	insert, err := db.Exec("INSERT INTO login_history (user_id, login_date, user_ip) VALUES (?, ?, ?)", id,currentTime,r.RemoteAddr )
	if err != nil {
		panic(err.Error())
	}
	insert.RowsAffected()
	//close DB connection
	defer db.Close()
	// return success
	success = true
	return success
	// If we reach this point, that means the users password was correct, and that they are authorized
	// The default 200 status is sent
}

func CreateSession(w http.ResponseWriter, r *http.Request) {
	// create vars
	var uid int
	var sessionID int
	var cookieID string
	// Generate unique session value
	id := ksuid.New()
	// write to DB
	// First check if the user exists in the session table
	db, err := sql.Open("mysql", DataSource)
	//getIDqs := fmt.Sprintf("select ID from users where username = '%s'", r.FormValue("username"))
	//fmt.Println(getIDqs)
	checkID, err := db.Query("select ID from users where username = ?", r.FormValue("username"))//(getIDqs)
	if err != nil {
		panic(err.Error())
	}
	for checkID.Next() {
		err := checkID.Scan(&uid)
		if err != nil {
			log.Fatal(err)
		}
	}
	checkSessionqs := fmt.Sprintf("select ID from user_session where userid = '%d'", uid)
	checkSessionID, err := db.Query(checkSessionqs)
	if err != nil {
		panic(err.Error())
	}
	for checkSessionID.Next() {
		err := checkSessionID.Scan(&sessionID)
		if err != nil {
			log.Fatal(err)
		}
	}
	if sessionID == 0 {
		currentTime := time.Now().Format("2006-01-02 15:04:05")
		//insertQry := fmt.Sprintf("insert into user_session (userid,sessionstart,sessionkey) value('%d', '%s', '%s')", uid, currentTime, id)
		insert, err := db.Query("insert into user_session (userid,sessionstart,sessionkey) value(?, ?, ?)", uid, currentTime, id)//(insertQry)
		if err != nil {
			panic(err.Error())
		}
		insert.Close()
	} else {
		currentTime := time.Now().Format("2006-01-02 15:04:05")
		//updateQry := fmt.Sprintf("update user_session set sessionstart = '%s', sessionkey = '%s' where userid = '%d'", currentTime, id, uid)
		update, err := db.Query("update user_session set sessionstart = ?, sessionkey = ? where userid = ?", currentTime, id, uid)//(updateQry)
		if err != nil {
			panic(err.Error())
		}
		update.Close()
	}
	checkID.Close()
	//create cookie on client
	cookieID = strconv.Itoa(uid) + "/" + id.String()
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{Name: "goodadvice", Value: cookieID , Expires: expiration}
	http.SetCookie(w, &cookie)
}

func ValidateSession(w http.ResponseWriter, r *http.Request) Cookie {
	// create vars
	var cookieID string
	var sessionID int
	var sessionAge time.Time
	var isAdmin string
	//var c Cookie
	c := validateCookie(w, r)
	if !c.Exists {
		return c
	}
	// Generate unique session value
	suid := ksuid.New()
	// write to DB
	db, err := sql.Open("mysql", DataSource)
	// validate session is LESS then 48 hours old
	//checkSessionAgeqs := fmt.Sprintf("select ID,sessionstart from user_session where userid = '%s' and sessionkey = '%s'", c.Uid, c.Sessionkey)
	checkSessionAge, err := db.Query("select ID,sessionstart from user_session where userid = ? and sessionkey = ?", c.Uid, c.Sessionkey)//(checkSessionAgeqs)
	if err != nil {
		panic(err.Error())
	}
	for checkSessionAge.Next() {
		err := checkSessionAge.Scan(&sessionID,&sessionAge)
		if err != nil {
			log.Fatal(err)
		}
	}

	expires := time.Now().Local().Add(-48 * time.Hour)//.Unix
	sessionAge,	expires = sessionAge.UTC(), expires.UTC()
	//var exp bool
	if expires.After(sessionAge) {
		//user will be redirected to login
		c.Exists = false
		c.Uid = ""
		c.Sessionkey = ""
	} else {
		currentTime := time.Now().Format("2006-01-02 15:04:05")
		//updateQry := fmt.Sprintf("update user_session set sessionstart = '%s', sessionkey = '%s' where ID = '%d'", currentTime, suid, sessionID)
		update, err := db.Query("update user_session set sessionstart = ?, sessionkey = ? where ID = ?", currentTime, suid, sessionID)//(updateQry)
		if err != nil {
			panic(err.Error())
		}
		update.Close()
		// update cookie on client
		cookieID = c.Uid + "/" + suid.String()
		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{Name: "goodadvice", Value: cookieID , Expires: expiration}
		http.SetCookie(w, &cookie)
		c.Exists = true
	}
	// Check if user is Admin
	checkAdmin, err := db.Query("select isadmin from users where ID = ?", c.Uid)//(checkAdminqs)
	if err != nil {
		panic(err.Error())
	}
	for checkAdmin.Next() {
		err := checkAdmin.Scan(&isAdmin)
		if err != nil {
			log.Fatal(err)
		}
	}
	if isAdmin == "5" {
		c.Isadmin = true
	} else {
		c.Isadmin = false
	}
	return c
}


func validateCookie (w http.ResponseWriter, r *http.Request) Cookie {
	var c Cookie
	cookie, err := r.Cookie("goodadvice")
	// No cookie then get guest WOD page
	if err != nil {
		// if not exist redirect to login page
		c.Exists = false
	} else if err == nil {
		cookievalue := cookie.Value
		splitcookie := strings.Split(cookievalue, "/")
		c.Exists = true
		c.Uid = splitcookie[0]
		c.Sessionkey = splitcookie[1]
	}
	return c
}

