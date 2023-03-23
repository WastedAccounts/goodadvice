package workouts

import (
	"database/sql"
	"fmt"
	"goodadvice/v1/datasource"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Workout struct {
	ID           int    //`json:"ID"`
	Name         string //`json:"Name"`
	Strength     string //`json:"Strength"`
	Pace         string //`json:"Pace"`
	Conditioning string //`json:"Conditioning"`
	Date         string //`json:"Date"`
	WODworkout   string
	Linkhidden   string
	Type         string
	Types        []string
	Template     string
	//DOW string `json:"DOW"`
}

type WorkoutNotes struct {
	ID       string // `json:"ID"`
	WoId     string //`json:"WoId"`
	UserName string //`json:"UserName"`
	UserId   string //`json:"UserId"`
	Notes    string //`json:"Notes"`
	Minutes  string
	Seconds  string
	Rounds   string
	Reps     string
	Loved    sql.NullString //string
	Hated    sql.NullString //string
}

type WodUser struct {
	ID           int
	UserName     string
	FirstName    string
	LastName     string
	EmailAddress string
	Greeting     string
}

type AddWorkout struct {
	ID           string //int    //`json:"ID"`
	Name         string //`json:"Name"`
	Strength     string //`json:"Strength"`
	Pace         string //`json:"Pace"`
	Conditioning string //`json:"Conditioning"`
	Date         string //`json:"Date"`
	Message      string
	CreatedBy    string
	WODworkout   string
	Linkhidden   string
	Type         string
	Types        []string
	Template     string
}

type EditWorkout struct {
	ID           string
	Name         string
	Strength     string
	Pace         string
	Conditioning string
	Date         string
	Message      string
	WODworkout   string
	Type         string
	Types        []string
}

// START - Get Workout of the Day functions

// GetWOD get today's workout and post it to /wod on first page load. This will get a user created WOD first, Then look for a Coach assigned, THEN grab the WOD if that's all that's left.
func GetWOD(uid string, r *http.Request) (Workout, WorkoutNotes, WodUser) {
	//Set up Vars
	var wo Workout
	var id, wotypeid int
	var name, strength, pace, conditioning, date, wodworkout string
	var results *sql.Rows
	var err error

	// Query DB for CURRENT_DATE WOD
	if uid == "" {
		// If we don't have an ID we'll assume they're a guest or new here and get the latest WOD
		results, err = datasource.DBconn.Query("SELECT ID,wo_name,wo_strength,wo_pace,wo_conditioning,wo_date,wo_workoutoftheday,wo_type FROM workout WHERE wo_date = CURRENT_DATE() AND wo_workoutoftheday = 'Y'")
		defer results.Close()
		if err != nil {
			panic(err.Error())
		}
		// Load Workout struct for page load
		for results.Next() {
			err = results.Scan(&id, &name, &strength, &pace, &conditioning, &date, &wodworkout, &wotypeid)
			if err != nil {
				panic(err.Error())
			}
			wo = Workout{
				ID:           id,
				Name:         name,
				Strength:     strength,
				Pace:         pace,
				Conditioning: conditioning,
				Date:         date,
				WODworkout:   wodworkout,
			}
		}
	} else {
		// need to check for and then pull the workout
		// First we'll check if they wrote themself a work out for today
		results, err = datasource.DBconn.Query("SELECT ID,wo_name,wo_strength,wo_pace,wo_conditioning,wo_date,wo_workoutoftheday,wo_type FROM workout WHERE wo_date = CURRENT_DATE() AND wo_createdby = ?", uid)
		defer results.Close()
		if err != nil {
			panic(err.Error())
		}
		// Load Workout struct for page load
		for results.Next() {
			err = results.Scan(&id, &name, &strength, &pace, &conditioning, &date, &wodworkout, &wotypeid)
			if err != nil {
				panic(err.Error())
			}
			wo = Workout{
				ID:           id,
				Name:         name,
				Strength:     strength,
				Pace:         pace,
				Conditioning: conditioning,
				Date:         date,
				WODworkout:   wodworkout,
			}
		}
		// If we don't get results from that call then the user doesn't have a a user create workout for that day and we'll move on to see if they have a coach assigned workout
		if wo.ID == 0 {
			// Next well check if they have a Coach assigned workout
			// Need some way to get the correct WORKOUT_ID from a table and in put that here
			results, err = datasource.DBconn.Query("SELECT ID,wo_name,wo_strength,wo_pace,wo_conditioning,wo_date,wo_workoutoftheday,wo_type FROM workout WHERE wo_date = CURRENT_DATE() AND ID = ?", 0)
			defer results.Close()
			if err != nil {
				panic(err.Error())
			}
			// Load Workout struct for page load
			for results.Next() {
				err = results.Scan(&id, &name, &strength, &pace, &conditioning, &date, &wodworkout, &wotypeid)
				if err != nil {
					panic(err.Error())
				}
				wo = Workout{
					ID:           id,
					Name:         name,
					Strength:     strength,
					Pace:         pace,
					Conditioning: conditioning,
					Date:         date,
					WODworkout:   wodworkout,
				}
			}
		}
		// If we don't get results from that call then we'll just load the daily WOD
		if wo.ID == 0 {
			// If we don't have an ID or a user created or coach assigned WOD we'll assume they're a guest or new here and get the latest WOD
			results, err := datasource.DBconn.Query("SELECT ID,wo_name,wo_strength,wo_pace,wo_conditioning,wo_date,wo_workoutoftheday,wo_type FROM workout WHERE wo_date = CURRENT_DATE() AND wo_workoutoftheday = 'Y'")
			defer results.Close()
			if err != nil {
				panic(err.Error())
			}
			// Load Workout struct for page load
			for results.Next() {
				err = results.Scan(&id, &name, &strength, &pace, &conditioning, &date, &wodworkout, &wotypeid)
				if err != nil {
					panic(err.Error())
				}
				wo = Workout{
					ID:           id,
					Name:         name,
					Strength:     strength,
					Pace:         pace,
					Conditioning: conditioning,
					Date:         date,
					WODworkout:   wodworkout,
				}
			}
		}
	}

	// Data Ops
	splitdate := strings.Split(wo.Date, "T")
	wo.Date = splitdate[0]
	// Get wotypeName by ID
	wotype := getWorkoutTypeName(wotypeid)
	wo.Type = wotype
	wo.Template = strings.ReplaceAll(wotype, " ", "")

	// Load additional struct for usr and wod notes
	usr := getUser(uid)
	won := getWODNotes(strconv.Itoa(wo.ID), uid)

	// Send to Controller
	return wo, won, usr
}

// GetWODbydate get a workout by date select and post it to /wod
func GetWODbydate(d string, uid string) (Workout, WorkoutNotes, WodUser) {
	var wo Workout
	var id, wotypeid int
	var name, strength, pace, conditioning, date, wodworkout string
	var results *sql.Rows
	var err error

	// Clean up date format
	splitdate := strings.Split(d, " ")
	d = splitdate[0]

	//// Get default WOD if user does not have their own workout.
	//results, err := datasource.DBconn.Query("SELECT ID,wo_name,wo_strength,wo_pace,wo_conditioning,wo_date,wo_workoutoftheday,wo_type FROM workout WHERE wo_date = ? AND wo_workoutoftheday = 'Y'", d)
	//defer results.Close()
	//if err != nil {
	//	panic(err.Error())
	//}
	//
	//// load results into Workout struct
	//for results.Next() {
	//	err = results.Scan(&id, &name, &strength, &pace, &conditioning, &date, &wodworkout, &wotypeid)
	//	if err != nil {
	//		panic(err.Error())
	//	}
	//	wo = Workout{
	//		ID:           id,
	//		Name:         name,
	//		Strength:     strength,
	//		Pace:         pace,
	//		Conditioning: conditioning,
	//		Date:         date,
	//		WODworkout:   wodworkout,
	//	}
	//}
	// Query DB for CURRENT_DATE WOD
	if uid == "" {
		// If we don't have an ID we'll assume they're a guest or new here and get the latest WOD
		results, err = datasource.DBconn.Query("SELECT ID,wo_name,wo_strength,wo_pace,wo_conditioning,wo_date,wo_workoutoftheday,wo_type FROM workout WHERE wo_date = ? AND wo_workoutoftheday = 'Y'", d)
		defer results.Close()
		if err != nil {
			panic(err.Error())
		}
		// Load Workout struct for page load?
		for results.Next() {
			err = results.Scan(&id, &name, &strength, &pace, &conditioning, &date, &wodworkout, &wotypeid)
			if err != nil {
				panic(err.Error())
			}
			wo = Workout{
				ID:           id,
				Name:         name,
				Strength:     strength,
				Pace:         pace,
				Conditioning: conditioning,
				Date:         date,
				WODworkout:   wodworkout,
			}
		}
	} else {
		// need to check for and then pull the workout
		// First we'll check if they wrote themself a work out for today
		results, err = datasource.DBconn.Query("SELECT ID,wo_name,wo_strength,wo_pace,wo_conditioning,wo_date,wo_workoutoftheday,wo_type FROM workout WHERE wo_date = ? AND wo_createdby = ?", d, uid)
		defer results.Close()
		if err != nil {
			panic(err.Error())
		}
		// Load Workout struct for page load
		for results.Next() {
			err = results.Scan(&id, &name, &strength, &pace, &conditioning, &date, &wodworkout, &wotypeid)
			if err != nil {
				panic(err.Error())
			}
			wo = Workout{
				ID:           id,
				Name:         name,
				Strength:     strength,
				Pace:         pace,
				Conditioning: conditioning,
				Date:         date,
				WODworkout:   wodworkout,
			}
		}
		// If we don't get results from that call then the user doesn't have a a user create workout for that day and we'll move on to see if they have a coach assigned workout
		if wo.ID == 0 {
			// Next well check if they have a Coach assigned workout
			// Need some way to get the correct WORKOUT_ID from a table and in put that here
			results, err = datasource.DBconn.Query("SELECT ID,wo_name,wo_strength,wo_pace,wo_conditioning,wo_date,wo_workoutoftheday,wo_type FROM workout WHERE wo_date = ? AND ID = ?", d, 0)
			defer results.Close()
			if err != nil {
				panic(err.Error())
			}
			// Load Workout struct for page load
			for results.Next() {
				err = results.Scan(&id, &name, &strength, &pace, &conditioning, &date, &wodworkout, &wotypeid)
				if err != nil {
					panic(err.Error())
				}
				wo = Workout{
					ID:           id,
					Name:         name,
					Strength:     strength,
					Pace:         pace,
					Conditioning: conditioning,
					Date:         date,
					WODworkout:   wodworkout,
				}
			}
		}
		// If we don't get results from that call then we'll just load the daily WOD
		if wo.ID == 0 {
			// If we don't have an ID or a user created or coach assigned WOD we'll assume they're a guest or new here and get the latest WOD
			results, err := datasource.DBconn.Query("SELECT ID,wo_name,wo_strength,wo_pace,wo_conditioning,wo_date,wo_workoutoftheday,wo_type FROM workout WHERE wo_date = ? AND wo_workoutoftheday = 'Y'", d)
			defer results.Close()
			if err != nil {
				panic(err.Error())
			}
			// Load Workout struct for page load
			for results.Next() {
				err = results.Scan(&id, &name, &strength, &pace, &conditioning, &date, &wodworkout, &wotypeid)
				if err != nil {
					panic(err.Error())
				}
				wo = Workout{
					ID:           id,
					Name:         name,
					Strength:     strength,
					Pace:         pace,
					Conditioning: conditioning,
					Date:         date,
					WODworkout:   wodworkout,
				}
			}
		}
	}

	// Data ops
	if wo.ID == 0 {
		// handle no workout found condition - still display date selected on web
		//splitdate := strings.Split(d, " ")
		wo.Date = d //splitdate[0]
	} else {
		// cut date to just date, no time stamp
		splitdate := strings.Split(wo.Date, "T")
		wo.Date = splitdate[0]
	}

	// Get wotypeName by ID
	wotype := getWorkoutTypeName(wotypeid)
	wo.Type = wotype
	wo.Template = strings.ReplaceAll(wotype, " ", "")

	// Load additional struct for usr and wod notes
	won := getWODNotes(strconv.Itoa(wo.ID), uid)
	usr := getUser(uid)
	return wo, won, usr
}

// GetWODbyID get a workout by the ID -- Only returns the workout, no user values
func GetWODbyID(woid string, uid string) (AddWorkout, WorkoutNotes, WodUser) {
	var wo AddWorkout
	var id, wotypeid int
	var name, strength, pace, conditioning, date, createdby, wodworkout string
	//var types []string

	// Get default WOD if user does not have their own workout.
	results, err := datasource.DBconn.Query("SELECT ID ,wo_name, wo_strength, wo_pace, wo_conditioning, wo_date, wo_createdby, wo_workoutoftheday,wo_type FROM workout WHERE ID = ?", woid)
	defer results.Close()
	if err != nil {
		panic(err.Error())
	}

	// load results into Workout struct
	for results.Next() {
		err = results.Scan(&id, &name, &strength, &pace, &conditioning, &date, &createdby, &wodworkout, &wotypeid)
		if err != nil {
			panic(err.Error())
		}
		// Get wotypeName by ID
		template := getWorkoutTypeName(wotypeid)
		wo = AddWorkout{
			ID:           "",
			Name:         name,
			Strength:     strength,
			Pace:         pace,
			Conditioning: conditioning,
			Date:         date,
			CreatedBy:    createdby,
			WODworkout:   wodworkout,
			Template:     strings.ReplaceAll(template, " ", ""),
		}
	}

	// Data Ops
	splitdate := strings.Split(wo.Date, "T")
	wo.Date = splitdate[0]

	// Load additional struct for usr and wod notes
	usr := getUser(uid)
	won := getWODNotes(wo.ID, uid)
	wo.Types = GetWorkoutType()
	// Send to Controller
	return wo, won, usr
}

// GetWODNotes gets comment posted by user on WOD by ID
func getWODNotes(woid string, userid string) WorkoutNotes {
	var won WorkoutNotes
	var id string
	var notes, time, min, sec, amrapscore, rnds, reps, loved, hated string
	uid := userid

	// Query for workout notes
	results, err := datasource.DBconn.Query("SELECT c.ID,c.user_id,c.workout_id,c.comment,c.time,c.amrapscore,uwr.loved,uwr.hated FROM (SELECT 1) dummy LEFT JOIN comments c ON c.user_id = ? LEFT JOIN user_workout_rating uwr ON uwr.workoutid = c.workout_id WHERE c.workout_id = ?", uid, woid)
	defer results.Close()
	if err != nil {
		panic(err.Error())
	}
	for results.Next() {
		err = results.Scan(&id, &userid, &woid, &notes, &time, &amrapscore, &sql.NullString{String: loved, Valid: true}, &sql.NullString{String: hated, Valid: true})
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
		if amrapscore != "" {
			as := strings.Split(amrapscore, ":")
			rnds = as[0]
			reps = as[1]
		} else {
			rnds = ""
			reps = ""
		}
		if loved == "1" {
			loved = "checked"
		} else {
			loved = ""
		}
		if hated == "1" {
			hated = "checked"
		} else {
			hated = ""
		}
		won = WorkoutNotes{
			ID:      id,
			UserId:  userid,
			WoId:    woid,
			Notes:   notes,
			Minutes: min,
			Seconds: sec,
			Rounds:  rnds,
			Reps:    reps,
			Loved:   sql.NullString{String: loved, Valid: true},
			Hated:   sql.NullString{String: hated, Valid: true},
		}
	}
	return won
}

// SaveWODResults -
func SaveWorkoutResults(r *http.Request) {
	// string, uid string, woid string){
	// setup values from page
	var time, amrap string
	min := r.PostFormValue("minutes")
	sec := r.PostFormValue("seconds")
	if min == "" && sec == "" {
		time = ""
	} else {
		time = min + ":" + fmt.Sprintf("%02s", sec) //fmt.Sprintf("%02s", min)+
	}
	rounds := r.PostFormValue("rounds")
	reps := r.PostFormValue("reps")
	if rounds == "" && reps == "" {
		amrap = ""
	} else {
		amrap = rounds + ":" + reps
	}
	woid := r.PostFormValue("woid")
	uid := r.PostFormValue("uid")
	n := r.PostFormValue("notes")
	loved := r.PostFormValue("loved")
	hated := r.PostFormValue("hated")
	uidint, err := strconv.Atoi(uid)
	woidint, err := strconv.Atoi(woid)
	n = strings.Replace(n, "'", "\\'", -1)

	//check if notes for this work out exist already
	var checkValue int
	checkID, err := datasource.DBconn.Query("SELECT ID FROM comments WHERE user_id = ? AND workout_id = ?", uid, woid)
	defer checkID.Close()
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
		_, err := datasource.DBconn.Exec("INSERT INTO  comments (user_id,workout_id,comment,time,amrapscore) VALUES (?,?,?,?,?)", uidint, woidint, n, time, amrap)
		if err != nil {
			panic(err.Error())
		}
	} else {
		// if notes do exist, update them with the current values
		// This shouldn't overwrite but make a new note I think
		_, err := datasource.DBconn.Exec("UPDATE comments SET comment = ?, time = ?, amrapscore = ? WHERE ID = ? AND user_id = ? AND workout_id = ?", n, time, amrap, checkValue, uidint, woidint)
		if err != nil {
			panic(err.Error())
		}
	}
	checkID.Close()
	// check if user selected a work out rating
	// save if they did
	// If both are check just ignore as it cancels itself
	if (loved == "on" || hated == "on") && !(loved == "on" && hated == "on") {
		// Set up var for function
		var rateID, lovedval, hatedval int
		// check if this has been rating previously
		ratechk, err := datasource.DBconn.Query("SELECT ID FROM user_workout_rating WHERE userid = ? AND workoutid = ?;", uid, woid)
		defer ratechk.Close()
		if err != nil {
			panic(err.Error())
		}
		for ratechk.Next() {
			err := ratechk.Scan(&rateID)
			if err != nil {
				log.Fatal(err)
			}
		}
		// check which value is checked
		if loved == "on" {
			lovedval = 1
		}
		if hated == "on" {
			hatedval = 1
		}
		if rateID == 0 {
			// if no previous rating then insert a new record
			_, err := datasource.DBconn.Exec("INSERT INTO user_workout_rating (userid,workoutid,loved,hated) VALUE (?,?,?,?)", uidint, woidint, lovedval, hatedval)
			if err != nil {
				panic(err.Error())
			}
			//insert.RowsAffected()
		} else {
			// if a rating already exists then increment the existing record
			if loved == "on" {
				// increment loved value by 1
				_, err := datasource.DBconn.Exec("UPDATE user_workout_rating SET loved = loved+1 WHERE ID = ?;", rateID)
				if err != nil {
					panic(err.Error())
				}
				//update.RowsAffected()
			} else if hated == "on" {
				// increment hated value by 1
				_, err := datasource.DBconn.Exec("UPDATE user_workout_rating SET hated = hated+1 WHERE ID = ?;", rateID)
				if err != nil {
					panic(err.Error())
				}
				//update.RowsAffected()
			}
		}
	}
}

// GetUser -
func getUser(uid string) WodUser {
	// Set up vars
	var wu WodUser
	var id int
	var username, firstname, lastname, emailaddress string

	// Query DB for user by ID
	results, err := datasource.DBconn.Query("SELECT ID, username, firstname,lastlogindate,emailaddress FROM users WHERE ID = ?", uid)
	defer results.Close()
	if err != nil {
		panic(err.Error())
	}
	for results.Next() {
		err = results.Scan(&id, &username, &firstname, &lastname, &emailaddress)
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

	// Call for a random greeting
	wu.Greeting = getRandomGreeting()

	// Return user data
	return wu
}

// getRandomGreeting - Get totally random greeting from the greeting table
func getRandomGreeting() string {
	// Set up vars
	type GreetID struct {
		ID int
	}
	var gid []GreetID
	var greeting string
	var id int

	// Query DB for list of indexex
	ids, err := datasource.DBconn.Query("SELECT ID FROM greetings")
	defer ids.Close()
	if err != nil {
		panic(err.Error())
	}
	for ids.Next() {
		err = ids.Scan(&id)
		if err != nil {
			panic(err.Error())
		}
		gid = append(gid, GreetID{id})
	}

	// Get Random index for greeting value
	randomIndex := rand.Intn(len(gid))
	pick := gid[randomIndex]

	// Query DB for randomized index id
	result, err := datasource.DBconn.Query("SELECT greeting FROM greetings WHERE ID = ?", pick.ID)
	defer result.Close()
	if err != nil {
		panic(err.Error())
	}
	for result.Next() {
		err = result.Scan(&greeting)
		if err != nil {
			panic(err.Error())
		}
	}

	// Return value
	return greeting
}

// GetRandomWorkout - Get totally random workout from the workout table
func GetRandomWorkout() string {
	type WorkoutID struct {
		WOID int
	}
	var woid []WorkoutID
	var wodate string
	var id int

	// get list of ALL workout IDs from workouts table
	ids, err := datasource.DBconn.Query("SELECT ID FROM workout;")
	defer ids.Close()
	if err != nil {
		panic(err.Error())
	}
	for ids.Next() {
		err = ids.Scan(&id)
		if err != nil {
			panic(err.Error())
		}
		woid = append(woid, WorkoutID{id})
	}
	//// I'm bailing on this system for now but I'll leave it just in case I bring it back .
	// This is to add weight to loved or hated workouts when generating random workouts.
	//if uid != "" {
	//	// get list of ALL Loved workout IDs from user_workout_rating table
	//	lids, err := datasource.DBconn.Query("SELECT workoutid FROM user_workout_rating WHERE userid = ? AND userrating = 1;", uid)
	//  defer lids.Close()
	//	if err != nil {
	//		panic(err.Error())
	//	}
	//	// add Loved IDs to list
	//	for lids.Next() {
	//		err = lids.Scan(&id)
	//		if err != nil {
	//			panic(err.Error())
	//		}
	//		woid = append(woid, WorkoutID{id})
	//	}
	//	// add them a second time to
	//	for lids.Next() {
	//		err = lids.Scan(&id)
	//		if err != nil {
	//			panic(err.Error())
	//		}
	//		woid = append(woid, WorkoutID{id})
	//	}
	//	// get all Hated IDs
	//	hids, err := datasource.DBconn.Query("SELECT workoutid FROM user_workout_rating WHERE userid = ? AND userrating = 2;", uid)
	//	defer hids.Close()
	//	if err != nil {
	//		panic(err.Error())
	//	}
	//	// add Loved IDs to list
	//	for hids.Next() {
	//		err = hids.Scan(&id)
	//		if err != nil {
	//			panic(err.Error())
	//		}
	//		woid = append(woid, WorkoutID{id})
	//	}
	//}
	// Pick random workout ID
	randomIndex := rand.Intn(len(woid))
	pick := woid[randomIndex]

	// Now get workout date from database
	date, err := datasource.DBconn.Query("SELECT wo_date FROM workout where ID = ?;", pick.WOID)
	defer date.Close()
	if err != nil {
		panic(err.Error())
	}
	for date.Next() {
		err = date.Scan(&wodate)
		if err != nil {
			panic(err.Error())
		}
	}
	// return date back to controller to call getWODbydate
	return wodate
}

// function to get workout types from DB in models.workout.go
func GetWorkoutType() []string {
	// vars
	var wotypes []string //AddWorkout
	var wotype string

	// get list of ALL workout types
	ids, err := datasource.DBconn.Query("SELECT wt_type FROM workout_type;")
	defer ids.Close()
	if err != nil {
		panic(err.Error())
	}
	for ids.Next() {
		err = ids.Scan(&wotype)
		if err != nil {
			panic(err.Error())
		}
		//wotypes.Types = append(wotypes.Types, wotype)
		wotypes = append(wotypes, wotype)
	}
	return wotypes
}

func getWorkoutTypeID(wotype string) int {
	// vars
	var wotypeid int
	// get list of ALL workout types
	wotypeidresults, err := datasource.DBconn.Query("SELECT ID FROM workout_type WHERE wt_type = ?;", wotype)
	defer wotypeidresults.Close()
	if err != nil {
		panic(err.Error())
	}
	for wotypeidresults.Next() {
		err = wotypeidresults.Scan(&wotypeid)
		if err != nil {
			panic(err.Error())
		}
	}
	return wotypeid
}

func getWorkoutTypeName(woid int) string {
	// vars
	var wotypename string
	// get list of ALL workout types
	wotypenameresults, err := datasource.DBconn.Query("SELECT wt_type FROM workout_type WHERE ID = ?;", woid)
	defer wotypenameresults.Close()
	if err != nil {
		panic(err.Error())
	}
	for wotypenameresults.Next() {
		err = wotypenameresults.Scan(&wotypename)
		if err != nil {
			panic(err.Error())
		}
	}
	return wotypename
}

// END - Get Workout of the Day functions

// START - Daily WOD functions - Add, Change, Select for working with the Workout of the Day

// AdminAddWOD - Add new WOD to workout table
func AdminAddWOD(r *http.Request, uid string) AddWorkout {
	// Vars
	awo := AddWorkout{
		ID:           "",
		Name:         r.FormValue("name"),
		Strength:     r.FormValue("strength"),
		Pace:         r.FormValue("pace"),
		Conditioning: r.FormValue("conditioning"),
		Date:         r.FormValue("date"),
		Message:      "",
		WODworkout:   r.PostFormValue("wodcb"),
	}

	// Check if this is a WOD workout or a just a usr/random workout
	if awo.WODworkout == "on" {
		// Write WOD workout to DB if it's a designated WOD workout
		checkWOD, err := datasource.DBconn.Query("SELECT * FROM workout WHERE wo_date = ? AND wo_workoutoftheday = 'Y'", awo.Date)
		defer checkWOD.Close()
		if err != nil {
			panic(err.Error())
		}
		if checkWOD.Next() != false {
			awo.Message = "A WOD already exists for " + awo.Date
		} else {
			insert, err := datasource.DBconn.Exec("INSERT INTO workout (wo_name, wo_strength, wo_pace, wo_conditioning, wo_date, wo_createdby,wo_workoutoftheday) VALUES (?,?,?,?,?,?,'Y')", awo.Name, awo.Strength, awo.Pace, awo.Conditioning, awo.Date, uid)
			if err != nil {
				panic(err.Error())
			}
			insertid, _ := insert.LastInsertId()
			awo.ID = strconv.FormatInt(insertid, 10)
		}
	} else {
		insert, err := datasource.DBconn.Exec("INSERT INTO workout (wo_name, wo_strength, wo_pace, wo_conditioning, wo_date, wo_createdby,wo_workoutoftheday) VALUES (?,?,?,?,?,?,'N')", awo.Name, awo.Strength, awo.Pace, awo.Conditioning, awo.Date, uid)
		if err != nil {
			panic(err.Error())
		}
		insertid, _ := insert.LastInsertId()
		awo.ID = strconv.FormatInt(insertid, 10)
	}
	// Date Ops
	if awo.WODworkout == "on" { // set value so we can set the checked state of the check box
		awo.WODworkout = "checked"
	} else { // Clear value if not so we don't break the page load
		awo.WODworkout = ""
	}
	// Set up date for display
	splitdate := strings.Split(awo.Date, "T")
	awo.Date = splitdate[0]

	// Return data
	return awo
}

// UserAddWOD - Add new WOD to workout table
func AddWOD(r *http.Request, uid string, edit bool) AddWorkout {
	// vars
	var wodworkout, wodworkoutchkbx string

	// Set wodworkout based on checkbox
	if r.PostFormValue("wodcb") == "on" {
		wodworkout = "Y"
		wodworkoutchkbx = "checked"
	} else {
		wodworkout = "N"
		wodworkoutchkbx = ""
	}
	awo := AddWorkout{
		ID:           "",
		Name:         r.FormValue("name"),
		Strength:     r.FormValue("strength"),
		Pace:         r.FormValue("pace"),
		Conditioning: r.FormValue("conditioning"),
		Date:         r.FormValue("date"),
		Message:      "",
		WODworkout:   wodworkoutchkbx,
		Types:        GetWorkoutType(),
		Type:         r.FormValue("typeddl"),
	}

	// Check adn make sure WOD is valid and not a dub
	awo.Message = checkWODValues(awo.Date, awo.WODworkout, uid, edit)

	if awo.Message != "" { // If not valid retun to user
		// Data Ops
		// Set up date for display
		splitdate := strings.Split(awo.Date, "T")
		awo.Date = splitdate[0]

		// Return struct to report issues to users
		return awo
	} else { // If value write to DB
		// Get ID of workouttype for writing to workouts table
		wotypeid := getWorkoutTypeID(awo.Type)

		// Write new workout to DB
		insert, err := datasource.DBconn.Exec("INSERT INTO workout (wo_name, wo_strength, wo_pace, wo_conditioning, wo_date, wo_createdby,wo_workoutoftheday, wo_type) VALUES (?,?,?,?,?,?,?,?)", awo.Name, awo.Strength, awo.Pace, awo.Conditioning, awo.Date, uid, wodworkout, wotypeid)
		if err != nil {
			panic(err.Error())
		}
		insertid, _ := insert.LastInsertId()
		awo.ID = strconv.FormatInt(insertid, 10)

		// Data Ops
		// Set up date for display
		splitdate := strings.Split(awo.Date, "T")
		awo.Date = splitdate[0]

		// Return data
		return awo
	}
}

// GetAddWODbydate - Get's WOD by Date in Datepicker
func GetAddWODbydate(d string, uid string) Workout {
	// Setup Vars
	var wo Workout
	var id, wotype int
	var name, strength, pace, conditioning, wodworkout string
	var date string

	if uid == "" {
		// today's WOD workout from db
		results, err := datasource.DBconn.Query("SELECT ID ,wo_name, wo_strength, wo_pace, wo_conditioning, wo_date,wo_workoutoftheday,wo_type FROM workout WHERE wo_date = ? AND wo_workoutoftheday = 'Y'", d)
		defer results.Close()
		if err != nil {
			panic(err.Error())
		}
		// Load results of query into struct
		for results.Next() {
			err = results.Scan(&id, &name, &strength, &pace, &conditioning, &date, &wodworkout, &wotype)
			if err != nil {
				panic(err.Error())
			}
			if wodworkout == "Y" {
				// set value so we can set the checked state of the check box
				wodworkout = "checked"
			} else {
				wodworkout = ""
			}
			wo = Workout{
				ID:           id,
				Name:         name,
				Strength:     strength,
				Pace:         pace,
				Conditioning: conditioning,
				Date:         date,
				WODworkout:   wodworkout,
			}
		}
	} else {
		// today's WOD workout from db
		results, err := datasource.DBconn.Query("SELECT ID ,wo_name, wo_strength, wo_pace, wo_conditioning, wo_date,wo_workoutoftheday, wo_type FROM workout WHERE wo_date = ? AND wo_createdby = ? AND wo_workoutoftheday = 'N'", d, uid)
		defer results.Close()
		if err != nil {
			panic(err.Error())
		}
		// Load results of query into struct
		for results.Next() {
			err = results.Scan(&id, &name, &strength, &pace, &conditioning, &date, &wodworkout, &wotype)
			if err != nil {
				panic(err.Error())
			}
			if wodworkout == "Y" {
				// set value so we can set the checked state of the check box
				wodworkout = "checked"
			} else {
				wodworkout = ""
			}
			wo = Workout{
				ID:           id,
				Name:         name,
				Strength:     strength,
				Pace:         pace,
				Conditioning: conditioning,
				Date:         date,
				WODworkout:   wodworkout,
			}
		}
	}

	// Data ops
	splitdate := strings.Split(wo.Date, "T")
	wo.Date = splitdate[0]

	// Get workout types for pageload into DDL
	wo.Types = GetWorkoutType()

	// Get wotypeName by ID
	wo.Type = getWorkoutTypeName(wotype)

	// Return the value
	return wo
}

// GetAddWODbyID - Get's WOD by ID tag hidden on page
func GetAddWODbyID(woid string) AddWorkout {
	// Setup Vars
	var wo AddWorkout
	var wotype int
	var name, strength, pace, conditioning, wodworkout, id string
	var date string

	// today's WOD workout from db
	// Removed  "AND wo_workoutoftheday = 'Y'"
	// I don't care about that here, right, I have the ID so it's hella specific anyway, right?
	results, err := datasource.DBconn.Query("SELECT ID ,wo_name, wo_strength, wo_pace, wo_conditioning, wo_date,wo_workoutoftheday, wo_type FROM workout WHERE ID = ?", woid)
	defer results.Close()
	if err != nil {
		panic(err.Error())
	}

	// Load results of query into struct
	for results.Next() {
		err = results.Scan(&id, &name, &strength, &pace, &conditioning, &date, &wodworkout, &wotype)
		if err != nil {
			panic(err.Error())
		}
		wo = AddWorkout{
			ID:           id,
			Name:         name,
			Strength:     strength,
			Pace:         pace,
			Conditioning: conditioning,
			Date:         date,
			WODworkout:   wodworkout,
		} //u = append(results)   //u, results)
	}

	// Date Ops
	if wo.WODworkout == "Y" { // set value so we can set the checked state of the check box
		wo.WODworkout = "checked"
	} else { // Clear value if not so we don't break the page load
		wo.WODworkout = ""
	}
	// Set up date for display
	splitdate := strings.Split(wo.Date, "T")
	wo.Date = splitdate[0]

	// Get workout types for pageload into DDL
	wo.Types = GetWorkoutType()

	// Get wotypeName by ID
	wo.Type = getWorkoutTypeName(wotype)

	// Return the values
	return wo
}

// EditAddWOD - Saves changes to WOD by ID tag hidden on page
func EditAddWOD(r *http.Request, uid string, edit bool) string {
	// vars
	var wodworkout string

	// Set wodworkout based on checkbox
	if r.PostFormValue("wodcb") == "on" {
		wodworkout = "Y"
	} else {
		wodworkout = "N"
	}

	// Load struct from request
	ew := EditWorkout{
		ID:           r.FormValue("id"),
		Name:         r.FormValue("name"),
		Strength:     r.FormValue("strength"),
		Pace:         r.FormValue("pace"),
		Conditioning: r.FormValue("conditioning"),
		Date:         r.FormValue("date"),
		Message:      "",
		WODworkout:   wodworkout, // r.PostFormValue("wodcb"),
		Type:         r.FormValue("typeddl"),
	}

	// Check workout values
	msg := checkWODValues(ew.Date, ew.WODworkout, uid, edit)

	if msg == "" { // If we're good write to db
		// Get ID of workouttype for writing to workouts table
		wotypeid := getWorkoutTypeID(ew.Type)

		// Write to DB
		_, err := datasource.DBconn.Exec("UPDATE workout SET wo_name = ?, wo_strength= ?, wo_pace = ?, wo_conditioning = ?, wo_workoutoftheday = ?, wo_type = ? WHERE ID = ?", ew.Name, ew.Strength, ew.Pace, ew.Conditioning, wodworkout, wotypeid, ew.ID)
		if err != nil {
			panic(err.Error())
		}

		// Return all clear in an empty string
		return msg
	} else { // If we're not ok report issues back to user to correct
		// Return issue in msg string
		return msg
	}
}

// checkWODValues - Check values being passed it to make sure it's not a duplicate user or Daily WOD or too for in the past(7 days) or future (6 months)
func checkWODValues(date string, wodworkout string, uid string, edit bool) string {
	// vars
	var msg string

	// Check if this is a WOD workout or a just a usr/random workout
	workoutdate, _ := time.Parse("2006-01-02", date)
	if workoutdate.Before(time.Now().AddDate(0, 0, -7)) {
		msg += "Can't schedule workouts that far in the past\n"
	}
	// check workout date make sure it's not over 6 months in the future
	if workoutdate.After(time.Now().AddDate(0, 6, 0)) {
		msg += "Can't schedule workouts that far in the future\n"
	}

	if edit == false { // If we're editing then we don't care if the dates or createdby or wodworkouts are the same since we're updating.
		if wodworkout == "" {
			// Check if it's a duplicate for a user wod
			checkUserWOD, err := datasource.DBconn.Query("SELECT ID FROM workout WHERE wo_date = ? AND wo_createdby = ? AND wo_workoutoftheday = 'N'", date, uid)
			defer checkUserWOD.Close()
			if err != nil {
				panic(err.Error())
			}
			if checkUserWOD.Next() != false {
				msg += "You already have a workout for " + date + "\n"
			}
		} else if wodworkout == "on" {
			// Check if it's a duplicate for an admin WOD workout
			checkAdminWOD, err := datasource.DBconn.Query("SELECT ID FROM workout WHERE wo_date = ? AND wo_workoutoftheday = 'Y'", date)
			defer checkAdminWOD.Close()
			if err != nil {
				panic(err.Error())
			}
			if checkAdminWOD.Next() != false {
				msg += "A Daily WOD already exists for " + date + "\n"
			}
		}
	}

	// Return values
	return msg
}

// END - Daily WOD functions
