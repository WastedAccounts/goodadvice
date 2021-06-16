package models

import (
	"database/sql"
	"strings"
)

type Version struct {
	Version string
	VersionDate string
}

type MovementTypes struct {
	MovementType []string
}

type Users struct {
	ID string
	Username string
	Firstname string
	Emailaddress string
	Isactive string
	Isadmin string
}

func GetVersion() Version {
	var v Version
	var ver, date string
	db, err := sql.Open("mysql", DataSource)
	if err != nil {
		panic(err.Error())
	}
	checkversion, err := db.Query("select db_version, date_updated from goodadvice_db_version order by ID desc limit 1;")
	if err != nil {
		panic(err.Error())
	}
	for checkversion.Next() {
		err = checkversion.Scan(&ver,&date)
		if err != nil {
			panic(err.Error())
		}
	}
	splitdate := strings.Split(date, "T")
	date = splitdate[0]
	v.VersionDate = date
	v.Version = ver
	defer db.Close()
	return v
}

func GetMovementTypes() MovementTypes {
	var mt MovementTypes
	var movementtype string
	db, err := sql.Open("mysql", DataSource)
	if err != nil {
		panic(err.Error())
	}
	getmovements, err := db.Query("select movement_type from movement_types;")
	if err != nil {
		panic(err.Error())
	}
	for getmovements.Next() {
		err = getmovements.Scan(&movementtype)
		if err != nil {
			panic(err.Error())
		}
		mt.MovementType = append(mt.MovementType, movementtype)
	}

	defer db.Close()
	return mt
}

func SaveMovement(m string, mt string) {
	var id int
	db, err := sql.Open("mysql", DataSource)
	if err != nil {
		panic(err.Error())
	}
	getid, err := db.Query("select ID from movement_types where movement_type = ?;",mt)
	if err != nil {
		panic(err.Error())
	}
	for getid.Next() {
		err = getid.Scan(&id)
	}
	//updateQry := fmt.Sprintf("insert into movements (movementtype,movementname) values (?,?)",id, m)
	insert, err := db.Exec("insert into movements (movementtype,movementname) values (?,?)",id, m)
	if err != nil {
		panic(err.Error())
	}
	insert.RowsAffected()
	defer db.Close()
}

func GetUsers() []Users {
	var u []Users
	var id,username,firstname,emailaddress,isactive,isadmin string
	db, err := sql.Open("mysql", DataSource)
	if err != nil {
		panic(err.Error())
	}
	getusers, err := db.Query("SELECT ID,username, firstname, emailaddress, isactive, isadmin FROM users;")
	if err != nil {
		panic(err.Error())
	}
	for getusers.Next() {
		err = getusers.Scan(&id,&username,&firstname,&emailaddress,&isactive,&isadmin)
		if err != nil {
			panic(err.Error())
		}
		if isactive == "1" {
			isactive = "Yes"
		}else {
			isactive = "No"
		}
		if isadmin == "5" {
			isadmin = "Admin"
		} else if isadmin == "3" {
			isadmin = "Moderator"
		} else {
			isadmin = "User"
		}
		u = append(u, Users{ID: id,Username: username,Firstname: firstname,Emailaddress: emailaddress,Isactive: isactive,Isadmin: isadmin} )
		//u.ID = append(u.ID,id)
		//u.Username = append(u.Username,username)
		//u.Firstname = append(u.Firstname,firstname)
		//u.Emailaddress = append(u.Emailaddress,emailaddress)
		//u.Isactive = append(u.Isactive,isactive)
		//u.Isadmin = append(u.Isadmin,isadmin)
	}
	defer db.Close()
	return u
}