package old

//type Records struct {
//	Record string
//	Movements []string
//	Date string
//}
//
//type Userprofile struct {
//	Name string
//	Birthday string //time.Time
//	Weight string
//	Sex string
//	About string
//	Age int
//}

//type Addpr struct {
//	Uid string
//	MovementName string
//	PRvalue string
//	Date string
//}

//func PageLoadUserProfile(uid string) Userprofile {
//	//var r Records
//	var up Userprofile
//
//	// get user
//	up = Userprofile(profile.UserProfileLoad(uid))
//
//	//// THis code will fill the Records struct
//	//movementresults, err := db.Query("select m.movementname, u.prvalue, u.prdate From user_pr u join movements m ON m.ID = u.movementid where u.userid = ?", uid)
//	//if err != nil {
//	//	panic(err.Error())
//	//}
//	//for movementresults.Next() {
//	//	err = movementresults.Scan(&movement,&pr,&date)
//	//	if err != nil {
//	//		panic(err.Error())
//	//	}
//	//	d := strings.Split(date.String(), " ")
//	//	display += movement + ": " + pr + " set on: " + d[0] + "\r"
//	//}
//	//r.Record = display
//	//movements, err := db.Query("SELECT movementname FROM mjs.movements;")
//	//if err != nil {
//	//	panic(err.Error())
//	//}
//	//for movements.Next() {
//	//	err = movements.Scan(&movementname)
//	//	if err != nil {
//	//		panic(err.Error())
//	//	}
//	//	r.Movements = append(r.Movements ,movementname)
//	//}
//	//currentTime := time.Now()
//	//r.Date = currentTime.Format("01/02/2006")
//	//return r,up
//	return up
//}

//func AddRecord (addpr Addpr) {
//	var movementid string
//	db, err := sql.Open("mysql", datasource.DataSource)
//	if err != nil {
//		panic(err.Error())
//	}
//	defer db.Close()
//	mid, err := db.Query("SELECT ID FROM movements WHERE movementname = ?;",addpr.MovementName)
//	if err != nil {
//		panic(err.Error())
//	}
//	for mid.Next() {
//		err = mid.Scan(&movementid)
//		if err != nil {
//			panic(err.Error())
//		}
//	}
//	insert, err := db.Exec("INSERT INTO user_pr (userid,movementid,prvalue,prdate) VALUES (?, ?, ?, ?)",addpr.Uid,movementid,addpr.PRvalue,addpr.Date)
//	if err != nil {
//		panic(err.Error())
//	}
//	insert.RowsAffected()
//}
