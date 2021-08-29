package models

import (
	"database/sql"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
	"goodadvice/v1/datasource"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Credentials - Used for logging in a useer
type Credentials struct {
	Id       string
	Username string
	Password string
}

// UserAuth - Stores values for authenticating users around the app
type UserAuth struct {
	Exists     bool
	IsActive   bool
	IsAdmin    bool
	IsCoach    bool
	Uid        string
	Path       string
	Sessionkey string
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
	db, err := sql.Open("mysql", datasource.DataSource)
	if err != nil {
		panic(err.Error())
	}
	GetCreds, err := db.Query("select ID, password,isactive from users where username = ?", creds.Username) //("select password,isactive from users where username='?'", creds.Username)
	for GetCreds.Next() {
		err := GetCreds.Scan(&id, &password, &isactive)
		if err != nil {
			log.Fatal(err)
		}
	}
	// Check if user is active, if not fail login immediately
	if isactive == "0" {
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
	insert, err := db.Exec("INSERT INTO login_history (user_id, login_date, user_ip) VALUES (?, ?, ?)", id, currentTime, r.RemoteAddr)
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

// CreateSession - Creates a new session at login
func CreateSession(w http.ResponseWriter, r *http.Request) {
	// create vars
	var uid, sessionID int
	var cookieID string
	// Generate unique session value
	id := ksuid.New()

	// Open DB connection
	db, err := sql.Open("mysql", datasource.DataSource)
	defer db.Close()

	// Get user ID from users table by searching for username
	checkID, err := db.Query("select ID from users where username = ?", r.FormValue("username")) //(getIDqs)
	if err != nil {
		panic(err.Error())
	}
	for checkID.Next() {
		err := checkID.Scan(&uid)
		if err != nil {
			log.Fatal(err)
		}
	}

	checkSessionID, err := db.Query("select ID from user_session where userid = ?", uid)
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
		insert, err := db.Query("insert into user_session (userid,sessionstart,sessionkey) value(?, ?, ?)", uid, currentTime, id) //(insertQry)
		if err != nil {
			panic(err.Error())
		}
		insert.Close()
	} else {
		currentTime := time.Now().Format("2006-01-02 15:04:05")
		update, err := db.Query("update user_session set sessionstart = ?, sessionkey = ? where userid = ?", currentTime, id, uid) //(updateQry)
		if err != nil {
			panic(err.Error())
		}
		update.Close()
	}
	checkID.Close()
	//create cookie on client
	cookieID = strconv.Itoa(uid) + "/" + id.String()
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{
		Name:    "goodadvice",
		Value:   cookieID,
		Path:    "/",
		Expires: expiration,
	}
	http.SetCookie(w, &cookie)
}

// ValidateSession -
func ValidateSession(w http.ResponseWriter, r *http.Request) UserAuth {
	// create vars
	//var userauth UserAuth
	var cookieID, isAdmin, isActive string
	var sessionID int
	var sessionAge time.Time
	//var c Cookie
	userauth := ValidateCookie(w, r)
	if !userauth.Exists {
		userauth.IsActive = false
		return userauth
	}
	// write to DB
	db, err := sql.Open("mysql", datasource.DataSource)
	// Check if user is Active and their role
	checkAdmin, err := db.Query("select isactive,isadmin from users where ID = ?", userauth.Uid) //(checkAdminqs)
	if err != nil {
		panic(err.Error())
	}
	for checkAdmin.Next() {
		err := checkAdmin.Scan(&isActive, &isAdmin)
		if err != nil {
			log.Fatal(err)
		}
	}
	if isActive == "0" {
		userauth.IsActive = false
		return userauth
	} else if isActive == "1" {
		userauth.IsActive = true
	}
	if isAdmin == "5" {
		userauth.IsAdmin = true
	} else {
		userauth.IsAdmin = false
	}
	// Generate unique session value
	suid := ksuid.New()

	// validate session is LESS then 48 hours old
	//checkSessionAgeqs := fmt.Sprintf("select ID,sessionstart from user_session where userid = '%s' and sessionkey = '%s'", c.Uid, c.Sessionkey)
	checkSessionAge, err := db.Query("select ID,sessionstart from user_session where userid = ? and sessionkey = ?", userauth.Uid, userauth.Sessionkey) //(checkSessionAgeqs)
	if err != nil {
		panic(err.Error())
	}
	for checkSessionAge.Next() {
		err := checkSessionAge.Scan(&sessionID, &sessionAge)
		if err != nil {
			log.Fatal(err)
		}
	}

	expires := time.Now().Local().Add(-48 * time.Hour) //.Unix
	sessionAge, expires = sessionAge.UTC(), expires.UTC()
	//var exp bool
	if expires.After(sessionAge) {
		//user will be redirected to login
		userauth.Exists = false
		userauth.Uid = ""
		userauth.Sessionkey = ""
	} else {
		currentTime := time.Now().Format("2006-01-02 15:04:05")
		//updateQry := fmt.Sprintf("update user_session set sessionstart = '%s', sessionkey = '%s' where ID = '%d'", currentTime, suid, sessionID)
		update, err := db.Query("update user_session set sessionstart = ?, sessionkey = ? where ID = ?", currentTime, suid, sessionID) //(updateQry)
		if err != nil {
			panic(err.Error())
		}
		update.Close()
		// update cookie on client
		cookieID = userauth.Uid + "/" + suid.String()
		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{
			Name:    "goodadvice",
			Value:   cookieID,
			Path:    "/",
			Expires: expiration,
		}
		http.SetCookie(w, &cookie)
		userauth.Exists = true
	}
	return userauth
}

// ValidateCookie -
func ValidateCookie(w http.ResponseWriter, r *http.Request) UserAuth {
	var userauth UserAuth
	cookie, err := r.Cookie("goodadvice")
	// No cookie then get guest WOD page
	if err != nil {
		// if not exist redirect to login page
		userauth.Exists = false
	} else if err == nil {
		cookievalue := cookie.Value
		splitcookie := strings.Split(cookievalue, "/")
		userauth.Exists = true
		userauth.Uid = splitcookie[0]
		userauth.Sessionkey = splitcookie[1]
	}
	return userauth
}
