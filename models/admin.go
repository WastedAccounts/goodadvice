package models

import (
	"database/sql"
	"goodadvice/v1/datasource"
	"strings"
)

type Version struct {
	Version string
	VersionDate string
}

type Movements struct {
	Movement string
	MovementType string
}

//type MovementTypes struct {
//	Movementtype string
//}

type Users struct {
	ID string
	Username string
	Firstname string
	Emailaddress string
	Isactive string
	Isadmin string
}

type User struct {
	ID string
	Username string
	Firstname string
	Emailaddress string
	Isactive string
	Isadmin string
	Active string
	Role1 string
	Role2 string
}

func GetVersion() Version {
	var v Version
	var ver, date string
	db, err := sql.Open("mysql", datasource.DataSource)
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

func GetMovements() []Movements {
	var m []Movements
	var movement,movementtype string
	db, err := sql.Open("mysql", datasource.DataSource)
	if err != nil {
		panic(err.Error())
	}
	getmovements, err := db.Query("SELECT movementname,movement_type FROM movements\nINNER JOIN movement_types\nON movements.movementtype = movement_types.ID\nORDER BY movement_type,movementname;")
	if err != nil {
		panic(err.Error())
	}
	for getmovements.Next() {
		err = getmovements.Scan(&movement,&movementtype)
		if err != nil {
			panic(err.Error())
		}
		m = append(m, Movements{Movement: movement,MovementType: movementtype})
	}
	defer db.Close()
	return m
}

func GetMovementTypes() []string {
	var mt []string
	var movementtype string
	db, err := sql.Open("mysql", datasource.DataSource)
	if err != nil {
		panic(err.Error())
	}
	getmovementtypes, err := db.Query("SELECT movement_type FROM movement_types ORDER BY movement_type;")
	if err != nil {
		panic(err.Error())
	}
	for getmovementtypes.Next() {
		err = getmovementtypes.Scan(&movementtype)
		if err != nil {
			panic(err.Error())
		}
		mt = append(mt,movementtype)
	}
	defer db.Close()
	return mt
}

func SaveMovement(m string, mt string) {
	var id int
	db, err := sql.Open("mysql", datasource.DataSource)
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
	db, err := sql.Open("mysql", datasource.DataSource)
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
	}
	defer db.Close()
	return u
}

func AdminGetUser(id string) User {
	var u User
	var username,firstname,emailaddress,isactive,isadmin,active,role1,role2 string
	db, err := sql.Open("mysql", datasource.DataSource)
	if err != nil {
		panic(err.Error())
	}
	getusers, err := db.Query("SELECT ID,username, firstname, emailaddress, isactive, isadmin FROM users WHERE ID = ?;",id)
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
			active = "Deactivate"
		}else {
			isactive = "No"
			active = "Activate"
		}
		if isadmin == "5" {
			isadmin = "Admin"
			role1 = "User"
			role2 = "Moderator"
		} else if isadmin == "3" {
			isadmin = "Moderator"
			role1 = "User"
			role2 = "Admin"
		} else {
			isadmin = "User"
			role1 = "Moderator"
			role2 = "Admin"
		}
		u = User{
			ID:           id,
			Username:     username,
			Firstname:    firstname,
			Emailaddress: emailaddress,
			Isactive: 	  isactive,
			Isadmin:   	  isadmin,
			Active: 	  active,
			Role1: 		  role1,
			Role2: 		  role2,
		}
	}
	defer db.Close()
	return u
}

func UpdateUser(id string,v string) {
	var role string
	var active string
	db, err := sql.Open("mysql", datasource.DataSource)
	if err != nil {
		panic(err.Error())
	}
	activeValue := map[string]bool {
		"Activate": true,
		"Deactivate": true,
	}
	roleValue := map[string]bool {
		"User": true,
		"Moderator": true,
		"Admin": true,
	}
	if activeValue[v]{
		if v == "Activate"{
			active = "1"
		} else if v == "Deactivate"{
			active = "0"
		}
		update, err := db.Exec("UPDATE users SET isactive = ? WHERE ID = ?;",active,id)
		if err != nil {
			panic(err.Error())
		}
		update.RowsAffected()
	} else if roleValue[v] {
		if v == "User"{
			role = "0"
		} else if v == "Moderator"{
			role = "3"
		} else if v == "Admin"{
			role = "5"
		}
		update, err := db.Exec("UPDATE users SET isadmin = ? WHERE ID = ?;",role,id)
		if err != nil {
			panic(err.Error())
		}
		update.RowsAffected()
	}
	defer db.Close()

}