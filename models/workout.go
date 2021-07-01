package models

import (
	"database/sql"
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
	min := r.PostFormValue("minutes")
	sec := r.PostFormValue("seconds")
	time := min+":"+sec
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
	//check if notes exist
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
		insert, err := db.Exec("insert into comments (user_id,workout_id,comment,time) values (?,?,?,?)",uidint,woidint,n,time)
		if err != nil {
			panic(err.Error())
		}
		insert.RowsAffected()
	} else {
		update, err := db.Exec("update comments set comment = ?, time = ? where ID = ? and user_id = ? and workout_id = ?",n,time,checkValue,uidint,woidint)
		if err != nil {
			panic(err.Error())
		}
		update.RowsAffected()
	}
	checkID.Close()
	if loved == "on" || hated == "on" {
		var rating int
		if loved == "on" {
			rating = 1
		}
		if hated == "on" {
			rating = 2
		}
		insert, err := db.Exec("INSERT INTO user_workout_rating (userid,workoutid,userrating) VALUE (?,?,?)",uidint,woidint,rating)
		if err != nil {
			panic(err.Error())
		}
		insert.RowsAffected()
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