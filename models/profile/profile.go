package profile

import (
	"database/sql"
	"goodadvice/v1/models"
	"strings"
	"time"
)

type Records struct {
	Record string
	Movements []string
	Date string
}

type Userprofile struct {
	Name string
	Birthday string //time.Time
	Weight string
	Sex string
	About string
	Age int
}

func PageLoadAboutMe(uid string) (Records,Userprofile){
	var up Userprofile
	var r Records
	// Userprofile Struct vars
	var name, weight, sex, about string //, birthday string
	var birthday time.Time
	// Records struct var
	var movement, pr, display, movementname string
	var date time.Time
	db, err := sql.Open("mysql", models.DataSource)
	if err != nil {
		panic(err.Error())
	}

	// THis code will fill the Userprofile struct
	userresults, err := db.Query("SELECT username,userbirthday,usersex,userweight,userabout FROM mjs.user_profile where userid = ?", uid)
	if err != nil {
		panic(err.Error())
	}
	for userresults.Next() {
		err = userresults.Scan(&name,&birthday,&sex,&weight,&about)
		if err != nil {
			panic(err.Error())
		}
		age := age(birthday,time.Now())
		up = Userprofile{Name: name,Age: age,Birthday: birthday.Format("01/02/2006"),Weight: weight,Sex: sex,About: about}
	}

	// THis code will fill the Records struct
	movementresults, err := db.Query("select m.movementname, u.prvalue, u.prdate From user_pr u join movements m ON m.ID = u.movementid where u.userid = ?", uid)
	if err != nil {
		panic(err.Error())
	}
	for movementresults.Next() {
		err = movementresults.Scan(&movement,&pr,&date)
		if err != nil {
			panic(err.Error())
		}
		d := strings.Split(date.String(), " ")
		display += movement + ": " + pr + " set on: " + d[0] + "\r"
	}
	r.Record = display
	movements, err := db.Query("SELECT movementname FROM mjs.movements;")
	if err != nil {
		panic(err.Error())
	}
	for movements.Next() {
		err = movements.Scan(&movementname)
		if err != nil {
			panic(err.Error())
		}
		r.Movements = append(r.Movements ,movementname)
	}
	currentTime := time.Now()
	r.Date = currentTime.Format("01/02/2006")
	return r,up
}

func PageLoadPersonalRecords(uid string) (Records){
	var rec Records
	// Records struct var
	var movement, pr, display, movementname string
	var date time.Time
	db, err := sql.Open("mysql", models.DataSource)
	if err != nil {
		panic(err.Error())
	}

	// THis code will fill the Records struct
	movementresults, err := db.Query("select m.movementname, u.prvalue, u.prdate From user_pr u join movements m ON m.ID = u.movementid where u.userid = ?", uid)
	if err != nil {
		panic(err.Error())
	}
	for movementresults.Next() {
		err = movementresults.Scan(&movement,&pr,&date)
		if err != nil {
			panic(err.Error())
		}
		d := strings.Split(date.String(), " ")
		display += movement + ": " + pr + " set on: " + d[0] + "\r"
	}
	rec.Record = display
	movements, err := db.Query("SELECT movementname FROM mjs.movements;")
	if err != nil {
		panic(err.Error())
	}
	for movements.Next() {
		err = movements.Scan(&movementname)
		if err != nil {
			panic(err.Error())
		}
		rec.Movements = append(rec.Movements ,movementname)
	}
	currentTime := time.Now()
	rec.Date = currentTime.Format("01/02/2006")
	return rec
}

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
