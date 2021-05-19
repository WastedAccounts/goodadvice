package models

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
)

type AddWorkout struct {
	ID int `json:"ID"`
	Name string `json:"Name"`
	Strength string `json:"Strength"`
	Pace string `json:"Pace"`
	Conditioning string `json:"Conditioning"`
	Date string `json:"Date"`
}

func AddWOD(w http.ResponseWriter, r *http.Request,uid string) AddWorkout {
	// Open DB
	awo := AddWorkout{
		ID:           0,
		Name:         r.FormValue("name"),
		Strength:     r.FormValue("strength"),
		Pace:         r.FormValue("pace"),
		Conditioning: r.FormValue("conditioning"),
		Date:         r.FormValue("date"),
	}
	db, err := sql.Open("mysql", DataSource)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	insertQry := fmt.Sprintf("insert into workout (wo_name, wo_strength, wo_pace, wo_conditioning, wo_date, wo_createdby) values ('%s', '%s', '%s', '%s', '%s', '%s')", awo.Name, awo.Strength, awo.Pace, awo.Conditioning, awo.Date, uid)
	fmt.Println(insertQry)
	insert, err := db.Query(insertQry)
	if err != nil {
		panic(err.Error())
	}
	insert.Close()
	return awo
}

func GetAddWODbydate(w http.ResponseWriter, r *http.Request) Workout {
	var wo Workout
	var id int
	var name, strength, pace, conditioning string
	var date string
	db, err := sql.Open("mysql", DataSource)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	qs := fmt.Sprintf("select ID ,wo_name, wo_strength, wo_pace, wo_conditioning, wo_date from workout where wo_date = '%s'",r.FormValue("date"))
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
	splitdate := strings.Split(wo.Date, "T")
	wo.Date = splitdate[0]
	return wo
}

func EditAddWOD (w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", DataSource)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	qs := fmt.Sprintf("update workout set wo_name = '%s', wo_strength= '%s', wo_pace = '%s', wo_conditioning = '%s' where wo_date = '%s'",r.FormValue("name"),r.FormValue("strength"),r.FormValue("pace"),r.FormValue("conditioning"),r.FormValue("date"))
	update, err := db.Query(qs)
	fmt.Println(update)
	if err != nil {
		panic(err.Error())
	}
}