package profile

import (
	"goodadvice/v1/datasource"
	"net/http"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Records struct {
	Movement string
	Weight   string
	Date     string
	ID       string
	Time     string
	Minutes  string
	Seconds  string
	Notes    string
}

type Movements struct {
	Movements   []string
	Currentdate string
}

type Userinfo struct {
	Name     string `json:"Name"`
	Birthday string `json:"Birthday"`
	Weight   string `json:"Weight"`
	Sex      string `json:"Sex"`
	About    string `json:"About"`
	Age      int    `json:"Age"`
}

type Addpr struct {
	Uid          string
	MovementName string
	Weight       string
	Date         string
	Time         string
	Notes        string
}

// PageLoadAboutMe - loads user info and PR info for main profile page
func PageLoadUserProfile(uid string) ([]Records, Userinfo) {
	var up Userinfo
	var r []Records
	// Load up for page load
	up = LoadAboutMe(uid)
	r = LoadPersonalRecords(uid)
	// Need to add Goals to this call when I build it
	//gls := LoadGoals()
	return r, up
}

// PageLoadPersonalRecords - loads PR values to the Records Struct
//currently loads all records for all time
// After adding functions to edit PRs I can update this as I need to change the way I store the data so I can sort better
func LoadPersonalRecords(uid string) []Records {
	var rec []Records
	// Records struct var
	var movement, pr, prtime, id string //display,
	var date time.Time

	// THis code will fill the Records struct
	prs, err := datasource.DBconn.Query("SELECT m.movementname,u.prvalue,u.prdate,u.prtime,u.ID FROM user_pr u JOIN movements m ON m.ID = u.movementid WHERE u.userid = ? ORDER BY m.movementname,prdate desc", uid)
	defer prs.Close()
	if err != nil {
		panic(err.Error())
	}
	for prs.Next() {
		err = prs.Scan(&movement, &pr, &date, &prtime, &id)
		if err != nil {
			panic(err.Error())
		}
		d := strings.Split(date.String(), " ")
		rec = append(rec, Records{
			Movement: movement,
			Weight:   pr,
			Date:     d[0],
			Time:     prtime,
			ID:       id},
		)
	}
	return rec
}

// loadMovements - get all movements to load in DDL
func LoadMovements() Movements {
	var mov Movements
	var movementname string

	// query db
	movements, err := datasource.DBconn.Query("SELECT movementname FROM movements ORDER BY movementtype,movementname;")
	defer movements.Close()
	if err != nil {
		panic(err.Error())
	}
	for movements.Next() {
		err = movements.Scan(&movementname)
		if err != nil {
			panic(err.Error())
		}
		mov.Movements = append(mov.Movements, movementname)
	}
	currentTime := time.Now()
	mov.Currentdate = currentTime.Format("01/02/2006")
	return mov
}

// LoadAboutMe - Load personal data to the Userprofile Struct (name,weight,sex,about,birthday,age)
func LoadAboutMe(uid string) Userinfo {
	var up Userinfo
	// Userprofile Struct vars
	var name, weight, sex, about string //, birthday string
	var birthday time.Time

	// Load the Userprofile struct from db
	userresults, err := datasource.DBconn.Query("SELECT u.firstname,up.userbirthday,up.usersex,up.userweight,up.userabout FROM user_profile up JOIN users u ON u.ID = up.userid where userid = ?;", uid)
	defer userresults.Close()
	if err != nil {
		panic(err.Error())
	}
	for userresults.Next() {
		err = userresults.Scan(&name, &birthday, &sex, &weight, &about)
		if err != nil {
			panic(err.Error())
		}
		age := age(birthday, time.Now())
		up = Userinfo{Name: name, Age: age, Birthday: birthday.Format("01/02/2006"), Weight: weight, Sex: sex, About: about}
	}
	return up
}

func UpdateAboutMe(r *http.Request, id string) {
	// Manage incoming date values
	bday, err := time.Parse("01/02/2006", r.PostFormValue("bday"))
	if err != nil {
		panic(err.Error())
	}
	// Update the DB
	_, err = datasource.DBconn.Exec("UPDATE user_profile up JOIN users u ON u.ID = up.userid SET u.firstname = ?,up.userbirthday = ?,up.usersex = ?,up.userweight = ?,up.userabout = ? WHERE userid = ?;", r.PostFormValue("name"), bday, r.PostFormValue("sex"), r.PostFormValue("wgt"), r.PostFormValue("abme"), id)
	if err != nil {
		panic(err.Error())
	}
}

// AddRecord CHANGE to - SaveSinglePR - Write new PR value to database
func SaveNewPR(addpr Addpr) {
	var movementid string

	// query db
	mid, err := datasource.DBconn.Query("SELECT ID FROM movements WHERE movementname = ?;", addpr.MovementName)
	defer mid.Close()
	if err != nil {
		panic(err.Error())
	}
	for mid.Next() {
		err = mid.Scan(&movementid)
		if err != nil {
			panic(err.Error())
		}
	}
	// insert record to DB
	_, err = datasource.DBconn.Exec("INSERT INTO user_pr (userid,movementid,prvalue,prdate,prtime,prnotes) VALUES (?,?,?,?,?,?)", addpr.Uid, movementid, addpr.Weight, addpr.Date, addpr.Time, addpr.Notes)
	if err != nil {
		panic(err.Error())
	}
}

// LoadSinglePR Load a pr to editpr page for editing
func LoadSinglePR(uid string, prid string) (Records, []Records) {
	var r Records
	var rhist []Records
	var movement, value, id, prtime, notes string
	var date time.Time

	// Get single pr value
	rec, err := datasource.DBconn.Query("SELECT m.movementname, u.prvalue, u.prdate,u.ID,u.prtime,u.prnotes FROM user_pr u JOIN movements m ON m.ID = u.movementid WHERE u.userid = ? AND u.ID = ?;", uid, prid)
	defer rec.Close()
	if err != nil {
		panic(err.Error())
	}
	for rec.Next() {
		err = rec.Scan(&movement, &value, &date, &id, &prtime, &notes)
		if err != nil {
			panic(err.Error())
		}
		d := strings.Split(date.String(), " ")
		if prtime == "" {
			prtime = ":"
		}
		prtimesplit := strings.Split(prtime, ":")
		r = Records{
			Movement: movement,
			Weight:   value,
			Date:     d[0],
			ID:       id,
			Time:     prtime,
			Minutes:  prtimesplit[0],
			Seconds:  prtimesplit[1],
			Notes:    notes,
		}
	}
	// Get PR history
	prhist, err := datasource.DBconn.Query("SELECT m.movementname,u.prvalue,u.prdate,u.prtime,u.ID FROM user_pr u JOIN movements m ON m.ID = u.movementid WHERE u.userid = ? and m.movementname = ? ORDER BY m.movementname,prdate desc", uid, r.Movement)
	defer prhist.Close()
	if err != nil {
		panic(err.Error())
	}
	for prhist.Next() {
		err = prhist.Scan(&movement, &value, &date, &prtime, &id)
		if err != nil {
			panic(err.Error())
		}
		d := strings.Split(date.String(), " ")
		rhist = append(rhist, Records{
			Movement: movement,
			Weight:   value,
			Date:     d[0],
			ID:       id,
			Time:     prtime,
		})
	}
	return r, rhist
}

// UpdateSinglePR - Update a pr value after editing
func UpdateSinglePR(r *http.Request, id string) {
	// Format web values for storing in DB
	time := r.PostFormValue("minutes") + ":" + r.PostFormValue("seconds")


	// Update DB
	_, err := datasource.DBconn.Exec("UPDATE user_pr SET prvalue = ?, prdate = ?, prtime = ?, prnotes = ? WHERE ID = ?", r.PostFormValue("weight"), r.PostFormValue("date"), time, r.PostFormValue("notes"), r.PostFormValue("prid"))
	if err != nil {
		panic(err.Error())
	}

}

// // // //
// Local functions for calculating things
// // // //

// age - function to calculate age from a date -- Not sure it's accurate though
func age(birthdate, today time.Time) int {
	//https://forum.golangbridge.org/t/how-to-calculate-the-exact-age-from-given-date-until-today/20530/3
	today = today.In(birthdate.Location())
	ty, tm, td := today.Date()
	today = time.Date(ty, tm, td, 0, 0, 0, 0, time.UTC)
	by, bm, bd := birthdate.Date()
	birthdate = time.Date(by, bm, bd, 0, 0, 0, 0, time.UTC)
	if today.Before(birthdate) {
		return 0
	}
	age := ty - by
	anniversary := birthdate.AddDate(age, 0, 0)
	if anniversary.After(today) {
		age--
	}
	return age
}

// // // //
// Older unused functions I might need later
// // // //

// LoadALLPersonalRecords - loads all history of PR values to the Records Struct
func LoadALLPersonalRecords(uid string) Records {
	// Records struct var
	var rec Records

	// vars
	var movement, pr, display, movementname string
	var date time.Time


	// This code will fill the Records struct
	movementresults, err := datasource.DBconn.Query("select m.movementname, u.prvalue, u.prdate From user_pr u join movements m ON m.ID = u.movementid where u.userid = ?", uid)
	defer movementresults.Close()
	if err != nil {
		panic(err.Error())
	}
	for movementresults.Next() {
		err = movementresults.Scan(&movement, &pr, &date)
		if err != nil {
			panic(err.Error())
		}
		d := strings.Split(date.String(), " ")
		display += movement + ": " + pr + " set on: " + d[0] + "\r"
	}

	movements, err := datasource.DBconn.Query("SELECT movementname FROM movements;")
	defer movements.Close()
	if err != nil {
		panic(err.Error())
	}
	for movements.Next() {
		err = movements.Scan(&movementname)
		if err != nil {
			panic(err.Error())
		}
	}

	return rec
}
