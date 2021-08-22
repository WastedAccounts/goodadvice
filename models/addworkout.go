package models

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"goodadvice/v1/datasource"
	"net/http"
	//"os"
	"strings"
)

type AddWorkout struct {
	ID int //`json:"ID"`
	Name string //`json:"Name"`
	Strength string //`json:"Strength"`
	Pace string //`json:"Pace"`
	Conditioning string //`json:"Conditioning"`
	Date string //`json:"Date"`
	Message string
}

type EditWorkout struct {
	Name string
	Strength string
	Pace string
	Conditioning string
	Date string
}

func AddWOD(r *http.Request,uid string) AddWorkout {
	// Open DB
	awo := AddWorkout{
		ID:           0,
		Name:         r.FormValue("name"),
		Strength:     r.FormValue("strength"),
		Pace:         r.FormValue("pace"),
		Conditioning: r.FormValue("conditioning"),
		Date:         r.FormValue("date"),
		Message:      "",
	}
	db, err := sql.Open("mysql", datasource.DataSource)
	if err != nil {
		panic(err.Error())
	}
	checkdate, err := db.Query("select * from workout where wo_date = ?",awo.Date)
	if err != nil {
		panic(err.Error())
	}
	if checkdate.Next() != false  {
		awo.Message = "A workout already exists for " + awo.Date
	} else {
		insert, err := db.Exec("insert into workout (wo_name, wo_strength, wo_pace, wo_conditioning, wo_date, wo_createdby) values (?, ?, ?, ?, ?, ?)", awo.Name, awo.Strength, awo.Pace, awo.Conditioning, awo.Date, uid)
		if err != nil {
			panic(err.Error())
		}
		insert.RowsAffected()
	}
	defer db.Close()
	return awo
}

func GetAddWODbydate(d string) Workout {
	var wo Workout
	var id int
	var name, strength, pace, conditioning string
	var date string
	db, err := sql.Open("mysql", datasource.DataSource)
	if err != nil {
		panic(err.Error())
	}
	results, err := db.Query("select ID ,wo_name, wo_strength, wo_pace, wo_conditioning, wo_date from workout where wo_date = ?",d)
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
	splitdate := strings.Split(wo.Date, "T")
	wo.Date = splitdate[0]
	defer db.Close()
	return wo
}

func EditAddWOD (r *http.Request) {
	ew := EditWorkout{
		Name:         r.FormValue("name"),
		Strength:     r.FormValue("strength"),
		Pace:         r.FormValue("pace"),
		Conditioning: r.FormValue("conditioning"),
		Date:         r.FormValue("date"),
	}
	db, err := sql.Open("mysql", datasource.DataSource)
	if err != nil {
		panic(err.Error())
	}
	update, err := db.Exec("update workout set wo_name = ?, wo_strength= ?, wo_pace = ?, wo_conditioning = ? where wo_date = ?", ew.Name, ew.Strength, ew.Pace, ew.Conditioning, ew.Date)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(update.RowsAffected())
	defer db.Close()
}