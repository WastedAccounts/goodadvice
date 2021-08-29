package old

import (
	_ "github.com/go-sql-driver/mysql"
)

//type AddWorkout struct {
//	ID           int    //`json:"ID"`
//	Name         string //`json:"Name"`
//	Strength     string //`json:"Strength"`
//	Pace         string //`json:"Pace"`
//	Conditioning string //`json:"Conditioning"`
//	Date         string //`json:"Date"`
//	Message      string
//	WODworkout   string
//}
//
//type EditWorkout struct {
//	ID 			 string
//	Name         string
//	Strength     string
//	Pace         string
//	Conditioning string
//	Date         string
//	WODworkout   string
//}
//
//// Daily WOD functions - Add, Change, Select for working with the Workout of the Day
//// AddWOD - Add new WOD to workout table
//func AddWOD(r *http.Request, uid string) AddWorkout {
//	// Open DB
//	awo := AddWorkout{
//		ID:           0,
//		Name:         r.FormValue("name"),
//		Strength:     r.FormValue("strength"),
//		Pace:         r.FormValue("pace"),
//		Conditioning: r.FormValue("conditioning"),
//		Date:         r.FormValue("date"),
//		Message:      "",
//		WODworkout:   r.PostFormValue("wodcb"),
//	}
//	// Open DB Conn
//	db, err := sql.Open("mysql", datasource.DataSource)
//	if err != nil {
//		panic(err.Error())
//	}
//	defer db.Close()
//
//	// Check if this is a WOD workout or a just a usr/random workout
//	if awo.WODworkout == "on" {
//		// Write WOD workout to DB if it's a designated WOD workout
//		checkWOD, err := db.Query("select * from workout where wo_date = ? and wo_workoutoftheday = 'Y'", awo.Date)
//		if err != nil {
//			panic(err.Error())
//		}
//		if checkWOD.Next() != false {
//			awo.Message = "A WOD already exists for " + awo.Date
//		} else {
//			insert, err := db.Exec("insert into workout (wo_name, wo_strength, wo_pace, wo_conditioning, wo_date, wo_createdby,wo_workoutoftheday) values (?, ?, ?, ?, ?, 'Y')", awo.Name, awo.Strength, awo.Pace, awo.Conditioning, awo.Date, uid)
//			if err != nil {
//				panic(err.Error())
//			}
//			insert.RowsAffected()
//		}
//	} else {
//		insert, err := db.Exec("insert into workout (wo_name, wo_strength, wo_pace, wo_conditioning, wo_date, wo_createdby,wo_workoutoftheday) values (?, ?, ?, ?, ?, 'N')", awo.Name, awo.Strength, awo.Pace, awo.Conditioning, awo.Date, uid)
//		if err != nil {
//			panic(err.Error())
//		}
//		insert.RowsAffected()
//	}
//	return awo
//}
//
//// GetAddWODbydate - Get's WOD by Date in Datepicker
//func GetAddWODbydate(d string) Workout {
//	// Setup Vars
//	var wo Workout
//	var id int
//	var name, strength, pace, conditioning, wodworkout string
//	var date string
//
//	// Open DB Conn
//	db, err := sql.Open("mysql", datasource.DataSource)
//	if err != nil {
//		panic(err.Error())
//	}
//	defer db.Close()
//
//	// today's WOD workout from db
//	results, err := db.Query("select ID ,wo_name, wo_strength, wo_pace, wo_conditioning, wo_date,wo_workoutoftheday from workout where wo_date = ? AND wo_workoutoftheday = 'Y'", d)
//	if err != nil {
//		panic(err.Error())
//	}
//
//	// Load results of query into struct
//	for results.Next() {
//		err = results.Scan(&id,&name,&strength,&pace,&conditioning,&date,&wodworkout)
//		if err != nil {
//			panic(err.Error())
//		}
//		if wodworkout == "Y" {
//			// set value so we can set the checked state of the check box
//			wodworkout = "checked"
//		} else {
//			wodworkout = ""
//		}
//		wo = Workout{ID: id, Name: name, Strength: strength, Pace: pace, Conditioning: conditioning, Date: date, WODworkout: wodworkout} //u = append(results)   //u, results)
//	}
//
//	// Manage data so we display what we want
//	splitdate := strings.Split(wo.Date, "T")
//	wo.Date = splitdate[0]
//
//	// Return the value
//	return wo
//}
//
//// GetAddWODbyID - Get's WOD by ID tag hidden on page
//func GetAddWODbyID(woid string) Workout {
//	// Setup Vars
//	var wo Workout
//	var id int
//	var name, strength, pace, conditioning, wodworkout string
//	var date string
//
//	// Open DB Conn
//	db, err := sql.Open("mysql", datasource.DataSource)
//	if err != nil {
//		panic(err.Error())
//	}
//	defer db.Close()
//
//	// today's WOD workout from db
//	results, err := db.Query("select ID ,wo_name, wo_strength, wo_pace, wo_conditioning, wo_date,wo_workoutoftheday from workout where ID = ? AND wo_workoutoftheday = 'Y'", woid)
//	if err != nil {
//		panic(err.Error())
//	}
//
//	// Load results of query into struct
//	for results.Next() {
//		err = results.Scan(&id,&name,&strength,&pace,&conditioning,&date,&wodworkout)
//		if err != nil {
//			panic(err.Error())
//		}
//		if wodworkout == "Y" {
//			// set value so we can set the checked state of the check box
//			wodworkout = "checked"
//		} else {
//			wodworkout = ""
//		}
//		wo = Workout{ID: id, Name: name, Strength: strength, Pace: pace, Conditioning: conditioning, Date: date, WODworkout: wodworkout} //u = append(results)   //u, results)
//	}
//
//	// Manage data so we display what we want
//	splitdate := strings.Split(wo.Date, "T")
//	wo.Date = splitdate[0]
//
//	// Return the value
//	return wo
//}
//
//// EditAddWOD - Saves changes to WOD by ID tag hidden on page
//func EditAddWOD(r *http.Request) {
//	var wodworkout string
//	if r.PostFormValue("wodcb") == "on" {
//		wodworkout = "Y"
//	} else {
//		wodworkout = "N"
//	}
//
//	ew := EditWorkout{
//		ID: 		  r.FormValue("ID"),
//		Name:         r.FormValue("name"),
//		Strength:     r.FormValue("strength"),
//		Pace:         r.FormValue("pace"),
//		Conditioning: r.FormValue("conditioning"),
//		Date:         r.FormValue("date"),
//		WODworkout:   wodworkout,
//	}
//	db, err := sql.Open("mysql", datasource.DataSource)
//	if err != nil {
//		panic(err.Error())
//	}
//	update, err := db.Exec("update workout set wo_name = ?, wo_strength= ?, wo_pace = ?, wo_conditioning = ?, wo_workoutoftheday = ? where ID = ?", ew.Name, ew.Strength, ew.Pace, ew.Conditioning, ew.WODworkout, ew.ID)
//	if err != nil {
//		panic(err.Error())
//	}
//	fmt.Println(update.RowsAffected())
//	defer db.Close()
//}
