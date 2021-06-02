package models

import (
	"database/sql"
	//"fmt"
	//"golang.org/x/crypto/openpgp/packet"
	//"net/http"
	"strings"
	"time"
)

type Records struct {
	Record string
	Movements []string
	Date string
}

type Addpr struct {
	Uid string
	MovementName string
	PRvalue string
	Date string
}

func PageLoadUserProfile(uid string) Records {
	var r Records
	var name, pr, display, movementname string
	var date time.Time
	db, err := sql.Open("mysql", DataSource)
	if err != nil {
		panic(err.Error())
	}
	results, err := db.Query("select m.movementname, u.prvalue, u.prdate From user_pr u join movements m ON m.ID = u.movementid where u.userid = ?", uid)
	if err != nil {
		panic(err.Error())
	}
	for results.Next() {
		err = results.Scan(&name,&pr,&date)
		if err != nil {
			panic(err.Error())
		}
		d := strings.Split(date.String(), " ")
		display += name + ": " + pr + " set on: " + d[0] + "\r"
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
	return r
}

func AddRecord (addpr Addpr) {
	var movementid string
	db, err := sql.Open("mysql", DataSource)
	if err != nil {
		panic(err.Error())
	}
	mid, err := db.Query("SELECT ID FROM movements WHERE movementname = ?;",addpr.MovementName)
	if err != nil {
		panic(err.Error())
	}
	for mid.Next() {
		err = mid.Scan(&movementid)
		if err != nil {
			panic(err.Error())
		}
	}
	insert, err := db.Exec("INSERT INTO user_pr (userid,movementid,prvalue,prdate) VALUES (?, ?, ?, ?)",addpr.Uid,movementid,addpr.PRvalue,addpr.Date)
	if err != nil {
		panic(err.Error())
	}
	insert.RowsAffected()
}
