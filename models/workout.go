package models

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

type Workout struct {
	ID int //`json:"ID"`
	Name string //`json:"Name"`
	Strength string //`json:"Strength"`
	Pace string //`json:"Pace"`
	Conditioning string //`json:"Conditioning"`
	Date string //`json:"Date"`
	//DOW string `json:"DOW"`
}

type WorkoutNotes struct {
	ID    int   // `json:"ID"`
	WoId  int    //`json:"WoId"`
	UserName string //`json:"UserName"`
	UserId string 	 //`json:"UserId"`
	Notes string //`json:"Notes"`
	Minutes string
	Seconds string
}

type WodUser struct {
	ID    int
	UserName string
	FirstName string
	LastName string
	EmailAddress string
	Greeting string
}

type GreetID struct {
	ID int
}

type WorkoutID struct {
	WOID int
}



func (w Workout) Write(bytes []byte) (int, error) {
	panic("implement me")
}

func (w Workout) WriteHeader(statusCode int) {
	panic("implement me")
}

var (
	workout []*Workout
)

var (
	greetid []*GreetID
)

// GetWOD get today's workout and post it to /wod on first page load
func GetWOD(uid string) (Workout, WorkoutNotes, WodUser)  {
	var wo Workout
	var id int
	var name, strength, pace, conditioning string
	var date string
	db, err := sql.Open("mysql", DataSource)
	if err != nil {
		panic(err.Error())
	}
	//qs := "select ID ,wo_name, wo_strength, wo_pace, wo_conditioning, wo_date from workout where wo_date = CURDATE()"
	results, err := db.Query("select ID ,wo_name, wo_strength, wo_pace, wo_conditioning, wo_date from workout where wo_date = CURDATE()")
	if err != nil {
		panic(err.Error())
	}
	for results.Next() {
		err = results.Scan(&id,&name,&strength,&pace,&conditioning,&date)
		if err != nil {
			panic(err.Error())
		}
		wo = Workout{ID: id, Name: name, Strength: strength, Pace: pace, Conditioning: conditioning, Date: date}  //u = append(results)   //u, results)
	}
	usr := GetUser(uid)
	won := GetWODNotes(wo.ID, uid)
	defer db.Close()
	return wo, won, usr
}

// GetWOD get today's workout and post it to /wod on first page load
func GetWODGuest() Workout {
	var wo Workout
	var id int
	var name, strength, pace, conditioning string
	var date string
	db, err := sql.Open("mysql", DataSource)
	if err != nil {
		panic(err.Error())
	}
	results, err := db.Query("select ID ,wo_name, wo_strength, wo_pace, wo_conditioning, wo_date from workout where wo_date = CURDATE()")
	if err != nil {
		panic(err.Error())
	}
	for results.Next() {
		err = results.Scan(&id,&name,&strength,&pace,&conditioning,&date)
		if err != nil {
			panic(err.Error())
		}
		wo = Workout{ID: id, Name: name, Strength: strength, Pace: pace, Conditioning: conditioning, Date: date}  //u = append(results)   //u, results)
	}
	defer db.Close()
	return wo //, won, usr
}

// GetWODbydate get a workout by date slect and post it to /wod
func GetWODbydate(d string, uid string) (Workout, WorkoutNotes, WodUser) {
	var wo Workout
	var id int
	var name, strength, pace, conditioning, date string
	db, err := sql.Open("mysql", DataSource)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	qs := "select ID ,wo_name, wo_strength, wo_pace, wo_conditioning, wo_date from workout where wo_date = '" + d + "'"
	results, err := db.Query(qs)
	if err != nil {
		panic(err.Error())
	}
	for results.Next() {
		err = results.Scan(&id,&name,&strength,&pace,&conditioning,&date)
		if err != nil {
			panic(err.Error())
		}
		wo = Workout{ID: id, Name: name, Strength: strength, Pace: pace, Conditioning: conditioning, Date: date}  //u = append(results)   //u, results)
	}
	won := GetWODNotes(wo.ID, uid)
	usr := GetUser(uid)
	return wo, won, usr
}

// GetWODNotes gets comment posted by user on WOD by ID
func GetWODNotes(woid int, userid string) WorkoutNotes {
	var won WorkoutNotes
	var id int
	var notes,time,min,sec string
	uid := userid
	wid := strconv.Itoa(woid)
	db, err := sql.Open("mysql", DataSource)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	results, err := db.Query("SELECT ID,user_id,workout_id,comment,time FROM comments where user_id = ? and workout_id = ?",uid, wid)
	if err != nil {
		panic(err.Error())
	}
	for results.Next() {
		err = results.Scan(&id,&userid,&woid,&notes,&time)
		if err != nil {
			panic(err.Error())
		}
		if time != "" {
			t := strings.Split(time, ":")
			min = t[0]
			sec = t[1]
		} else {
			min = ""
			sec = ""
		}
		won = WorkoutNotes{ID: id, UserId: userid, WoId: woid, Notes: notes,Minutes: min,Seconds:sec }  //u = append(results)   //u, results)
	}
	return won
}

func SaveWODResults (r *http.Request) {
	// string, uid string, woid string){
	// setup values from page
	var time string
	min := r.PostFormValue("minutes")
	sec := r.PostFormValue("seconds")
	if min == "" && sec == ""{
		time = ""
	} else {
		time = min+":"+fmt.Sprintf("%02s", sec) //fmt.Sprintf("%02s", min)+
	}
	woid := r.PostFormValue("woid")
	uid := r.PostFormValue("uid")
	n := r.PostFormValue("notes")
	loved := r.PostFormValue("loved")
	hated := r.PostFormValue("hated")
	uidint, err := strconv.Atoi(uid)
	woidint, err := strconv.Atoi(woid)
	// Open DB connection
	n = strings.Replace(n, "'", "\\'", -1)
	db, err := sql.Open("mysql", DataSource)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	//check if notes for this work out exist already
	var checkValue int
	checkID, err := db.Query("select ID from comments where user_id = ? and workout_id = ?",uid,woid)
	if err != nil {
		panic(err.Error())
	}
	for checkID.Next() {
		err := checkID.Scan(&checkValue)
		if err != nil {
			log.Fatal(err)
		}
	}
	if checkValue == 0 {
		// if no notes exist the insert a new record
		insert, err := db.Exec("insert into comments (user_id,workout_id,comment,time) values (?,?,?,?)",uidint,woidint,n,time)
		if err != nil {
			panic(err.Error())
		}
		insert.RowsAffected()
	} else {
		// if notes do exist, update them with the current values
		// This shouldn't overwrite but make a new note I think
		update, err := db.Exec("update comments set comment = ?, time = ? where ID = ? and user_id = ? and workout_id = ?",n,time,checkValue,uidint,woidint)
		if err != nil {
			panic(err.Error())
		}
		update.RowsAffected()
	}
	checkID.Close()
	// check if user selected a work out rating
	// save if they did
	// If both are check just ignore as it cancels itself
	if (loved == "on" || hated == "on") && !(loved == "on" && hated == "on") {
		// Set up var for function
		var rateID,lovedval,hatedval int
		// check if this has been rating previously
		ratechk, err := db.Query("SELECT ID FROM user_workout_rating WHERE userid = ? AND workoutid = ?;",uid,woid)
		if err != nil {
			panic(err.Error())
		}
		for ratechk.Next() {
			err := ratechk.Scan(&rateID)
			if err != nil {
				log.Fatal(err)
			}
		}
		fmt.Println("ratechk:", ratechk)
		fmt.Println("rateID:", rateID)
		// check which value is checked
		if loved == "on" {
			lovedval = 1
		}
		if hated == "on" {
			hatedval = 1
		}
		if rateID == 0{
			// if no previous rating then insert a new record
			insert, err := db.Exec("INSERT INTO user_workout_rating (userid,workoutid,loved,hated) VALUE (?,?,?,?)",uidint,woidint,lovedval,hatedval)
			if err != nil {
				panic(err.Error())
			}
			insert.RowsAffected()
		} else {
			// if a rating already exists then increment the existing record
			if loved == "on" {
				// increment loved value by 1
				update, err := db.Exec("UPDATE user_workout_rating SET loved = loved+1 WHERE ID = ?;",rateID)
				if err != nil {
					panic(err.Error())
				}
				update.RowsAffected()
			} else if hated == "on" {
				// increment hated value by 1
				update, err := db.Exec("UPDATE user_workout_rating SET hated = hated+1 WHERE ID = ?;",rateID)
				if err != nil {
					panic(err.Error())
				}
				update.RowsAffected()
			}
		}
	}
}

func GetUser(uid string) WodUser{
	var wu WodUser
	var id int
	var username,firstname,lastname,emailaddress string

	db, err := sql.Open("mysql", DataSource)
	if err != nil {
		panic(err.Error())
	}
	results, err := db.Query("select ID, username, firstname,lastlogindate,emailaddress from users where ID = ?", uid)
	if err != nil {
		panic(err.Error())
	}
	for results.Next() {
		err = results.Scan(&id,&username,&firstname,&lastname,&emailaddress)
		if err != nil {
			panic(err.Error())
		}
		wu = WodUser{
			ID:           id,
			UserName:     username,
			FirstName:    firstname,
			LastName:     lastname,
			EmailAddress: emailaddress,
		}
	}
	wu.Greeting = getRandomGreeting()
	defer db.Close()
	return wu
}

func getRandomGreeting() string {
	var gid []GreetID
	var greeting string
	var id int
	db, err := sql.Open("mysql", DataSource)
	if err != nil {
		panic(err.Error())
	}
	ids, err := db.Query("select ID from greetings")
	//fmt.Println("ids:",ids)
	if err != nil {
		panic(err.Error())
	}
	for ids.Next() {
		err = ids.Scan(&id)
		if err != nil {
			panic(err.Error())
		}
		gid = append(gid, GreetID{id} )
	}

	randomIndex := rand.Intn(len(gid))
	pick := gid[randomIndex]

	result, err := db.Query("select greeting from greetings where ID = ?", pick.ID)
	if err != nil {
		panic(err.Error())
	}
	for result.Next() {
		err = result.Scan(&greeting)
		if err != nil {
			panic(err.Error())
		}
	}
	defer db.Close()
	return greeting
}

func GetRandomWorkout(uid string) string {
	var woid []WorkoutID
	var wodate string
	var id int
	db, err := sql.Open("mysql", DataSource)
	if err != nil {
		panic(err.Error())
	}
	// get list of ALL workout IDs from workouts table
	ids, err := db.Query("SELECT ID FROM mjs.workout;")
	if err != nil {
		panic(err.Error())
	}
	for ids.Next() {
		err = ids.Scan(&id)
		if err != nil {
			panic(err.Error())
		}
		woid = append(woid, WorkoutID{id} )
	}
	// Check if user is logged in
	if uid != "" {
		// get list of ALL Loved workout IDs from user_workout_rating table
		lids, err := db.Query("SELECT workoutid FROM user_workout_rating WHERE userid = ? AND userrating = 1;",uid)
		if err != nil {
			panic(err.Error())
		}
		// add Loved IDs to list
		for lids.Next() {
			err = lids.Scan(&id)
			if err != nil {
				panic(err.Error())
			}
			woid = append(woid, WorkoutID{id})
		}
		// add them a second time to
		for lids.Next() {
			err = lids.Scan(&id)
			if err != nil {
				panic(err.Error())
			}
			woid = append(woid, WorkoutID{id} )
		}
		// get all Hated IDs
		hids, err := db.Query("SELECT workoutid FROM user_workout_rating WHERE userid = ? AND userrating = 2;",uid)
		if err != nil {
			panic(err.Error())
		}
		// add Loved IDs to list
		for hids.Next() {
			err = hids.Scan(&id)
			if err != nil {
				panic(err.Error())
			}
			woid = append(woid, WorkoutID{id})
		}
	}
	// Pick random workout ID
	randomIndex := rand.Intn(len(woid))
	pick := woid[randomIndex]
	// Now get workout date from database
	date, err := db.Query("SELECT wo_date FROM workout where ID = ?;",pick.WOID)
	if err != nil {
		panic(err.Error())
	}
	for date.Next() {
		err = date.Scan(&wodate)
		if err != nil {
			panic(err.Error())
		}
	}
	defer db.Close()
	// return date back to controller to call getWODbydate
	return wodate
}