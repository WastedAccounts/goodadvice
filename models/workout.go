package models

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"math/rand"
	"strconv"
	"strings"
)

type Workout struct {
	ID int `json:"ID"`
	Name string `json:"Name"`
	Strength string `json:"Strength"`
	Pace string `json:"Pace"`
	Conditioning string `json:"Conditioning"`
	Date string`json:"Date"`
	//DOW string `json:"DOW"`
}

type WorkoutNotes struct {
	ID    int    `json:"ID"`
	WoId  int    `json:"WoId"`
	UserName string `json:"UserName"`
	UserId string 	 `json:"UserId"`
	Notes string `json:"Notes"`
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
	//t,_ := time.Parse(time.UnixDate, wo.Date.String())
	//year, month, day := time.Now().Date()
	//wo.DOW = t.Weekday().String() + " - " + month.String() + " " + strconv.Itoa(day) + " " + strconv.Itoa(year)
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
	//t,_ := time.Parse(time.UnixDate, wo.Date.String())
	//year, month, day := time.Now().Date()
	//wo.DOW = t.Weekday().String() + " - " + month.String() + " " + strconv.Itoa(day) + " " + strconv.Itoa(year)
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
	var notes string
	uid := userid
	wid := strconv.Itoa(woid)
	db, err := sql.Open("mysql", DataSource)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	qs := "SELECT ID,user_id, workout_id, comment FROM comments where user_id = " + uid + " and workout_id = " + wid
	results, err := db.Query(qs)
	if err != nil {
		panic(err.Error())
	}
	for results.Next() {
		err = results.Scan(&id,&userid,&woid,&notes)
		if err != nil {
			panic(err.Error())
		}
		won = WorkoutNotes{ID: id, UserId: userid, WoId: woid, Notes: notes}  //u = append(results)   //u, results)
	}
	return won
}

func PostWODNotes (n string, uid string, woid string){
	// Open DB connection
	uidint, err := strconv.Atoi(uid)
	woidint, err := strconv.Atoi(woid)
	n = strings.Replace(n, "'", "\\'", -1)
	db, err := sql.Open("mysql", DataSource)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	//check if notes exist
	var checkValue int
	qs := "select id from comments where user_id = '" + uid + "' and workout_id = '" + woid + "'"
	checkID, err := db.Query(qs)
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
		insertQry := fmt.Sprintf("insert into comments (user_id,workout_id,comment) values ('%d', '%d', '%s')",uidint,woidint,n)
		insert, err := db.Query(insertQry)
		if err != nil {
			panic(err.Error())
		}
		insert.Close()
	} else {
		updateQry := fmt.Sprintf("update comments set comment = '%s' where ID = '%d' and user_id = '%d' and workout_id = '%d'",n,checkValue, uidint,woidint)
		update, err := db.Query(updateQry)
		if err != nil {
			panic(err.Error())
		}
		update.Close()
	}
	checkID.Close()

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
	fmt.Println("ids:",ids)
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
	fmt.Println("result:",result)
	if err != nil {
		panic(err.Error())
	}
	for result.Next() {
		err = result.Scan(&greeting)
		if err != nil {
			panic(err.Error())
		}
	}
	fmt.Println("greeting:",greeting)
	defer db.Close()
	return greeting
}





//// AddWOD - Write new daily WOD to database
//func AddWOD1(wo Workout,w http.ResponseWriter, r *http.Request) (Workout, error) {
//	pace := wo.Pace
//	conditioning := wo.Conditioning
//	date := wo.Date //need to convert this to a date time somewhere
//	fmt.Println(wo.Name, wo.Strength, pace, conditioning, date)
//	// Open DB connection
//	db, err := sql.Open("mysql", DataSource)
//	if err != nil {
//		panic(err.Error())
//	}
//	defer db.Close()
//	insertQry := fmt.Sprintf("insert into workout (wo_name, wo_strength, wo_pace, wo_conditioning, wo_date) values ('%s', '%s', '%s', '%s', %t)", wo.Name, wo.Strength, pace, conditioning, date)
//	fmt.Println(insertQry)
//	insert, err := db.Query(insertQry)
//	if err != nil {
//		panic(err.Error())
//	}
//	insert.Close()
//	return wo, nil
//}
//
//
//
//// CheckIfWorkoutExists check if uset exist by their email
//// if they don't we return false and then we can add them
//// If they do then we return true and a fresh pull
//// of their info from the database to do more work
//// ID is not known to the user and emails are unique
//func CheckIfWorkoutExists(u User) (User, bool) {
//	emailAddress := u.EmailAddress
//	var count int
//	var exists bool
//	fmt.Println(emailAddress)
//	// Open DB connection
//	db, err := sql.Open("mysql", DataSource)
//	if err != nil {
//		panic(err.Error())
//	}
//	defer db.Close()
//	checkQry := fmt.Sprintf("select count(EmailAddress) from visitors where EmailAddress = '%s'", emailAddress)
//	check, err := db.Query(checkQry)
//	for check.Next() {
//		err = check.Scan(&count)
//	}
//	fmt.Println(count)
//	if count != 0 {
//		getuser := fmt.Sprintf("select ID,FirstName,LastName,EmailAddress,VisitDate from visitors where EmailAddress = '%s'", emailAddress)
//		getuserResults, err := db.Query(getuser)
//		if err != nil {
//			panic(err.Error())
//		}
//		for getuserResults.Next() {
//			err = getuserResults.Scan(&u.ID,&u.FirstName,&u.LastName,&u.EmailAddress,&u.VisitDate)
//		}
//		exists = true
//	}
//	if count == 0 {
//		exists = false
//	}
//	return u, exists
//}

// AddWorkout Adds new user to visitors table
// first checks if user exists based on email
//func AddWorkout(u User) (User, error) {
//	// create varables from User struct
//	firstName := u.FirstName
//	lastName := u.LastName
//	emailAddress := u.EmailAddress
//	fmt.Println(firstName, lastName, emailAddress)
//	// Open DB connection
//	db, err := sql.Open("mysql", DataSource)
//	if err != nil {
//		panic(err.Error())
//	}
//	defer db.Close()
//	insertQry := fmt.Sprintf("insert into visitors (FirstName,LastName,EmailAddress,VisitDate) values ('%s', '%s', '%s',NOW())", firstName, lastName, emailAddress)
//	fmt.Println(insertQry)
//	insert, err := db.Query(insertQry)
//	if err != nil {
//		panic(err.Error())
//	}
//	insert.Close()
//	return u, nil
//}

/* GetWorkoutByDate Not in use yet
func GetWorkoutByDate (id int) (User, error) {
	date := time.no
	var count int
	var exists bool
	fmt.Println(emailAddress)
	// Open DB connection
	db, err := sql.Open("mysql", DataSource)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	checkQry := fmt.Sprintf("select count(EmailAddress) from visitors where EmailAddress = '%s'", emailAddress)
	check, err := db.Query(checkQry)
	for check.Next() {
		err = check.Scan(&count)
	}
	return wo
}
*/

// UpdateWorkout Not in use yet
func UpdateWorkout(u User)	(User, error) {
	for i, candidate := range users {
		if candidate.ID == u.ID	{
			users[i] = &u
			return u, nil
		}
	}
	return User{}, fmt.Errorf("User with ID '%v' not found", u.ID)
}

func RemoveWorkoutByID(id int) error {
	//delete from visitors where ID = 8
	for i, u := range users {
		if u.ID == id {
			users = append(users[:i], users[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("User with ID '%v' not  found", id)
}
