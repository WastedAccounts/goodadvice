package models

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID int
	FirstName string
	LastName string
	EmailAddress string
	VisitDate string
}

var (
	users []*User
	nextID = 1
)


func GetUsers() []User {
	var u []User
	var userId int
	var firstName, lastName, emailAddress, visitDate string
	db, err := sql.Open("mysql", DataSource)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	results, err := db.Query("select ID,FirstName,LastName,EmailAddress,VisitDate from visitors")
	if err != nil {
		panic(err.Error())
	}
	for results.Next() {
		err = results.Scan(&userId,&firstName,&lastName,&emailAddress,&visitDate)
		fmt.Println(results)
		if err != nil {
			panic(err.Error())
		}
		u = append(u, User{ID: userId, FirstName: firstName, LastName: lastName, EmailAddress: emailAddress, VisitDate: visitDate}  ) //u = append(results)   //u, results)
	}
	return u
}

// CheckIfUserExists check if uset exist by their email
// if they don't we return false and then we can add them
// If they do then we return true and a fresh pull
// of their info from the database to do more work
// ID is not known to the user and emails are unique
func CheckIfUserExists(u User) (User, bool) {
	emailAddress := u.EmailAddress
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
	fmt.Println(count)
	if count != 0 {
		getuser := fmt.Sprintf("select ID,FirstName,LastName,EmailAddress,VisitDate from visitors where EmailAddress = '%s'", emailAddress)
		getuserResults, err := db.Query(getuser)
		if err != nil {
			panic(err.Error())
		}
		for getuserResults.Next() {
			err = getuserResults.Scan(&u.ID,&u.FirstName,&u.LastName,&u.EmailAddress,&u.VisitDate)
		}
		exists = true
	}
	if count == 0 {
		exists = false
	}
	return u, exists
}

// AddUser Adds new user to visitors table
// first checks if user exists based on email
func AddUser(u User) (User, error) {
	// create varables from User struct
	firstName := u.FirstName
	lastName := u.LastName
	emailAddress := u.EmailAddress
	fmt.Println(firstName, lastName, emailAddress)
	// Open DB connection
	db, err := sql.Open("mysql", DataSource)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	insertQry := fmt.Sprintf("insert into visitors (FirstName,LastName,EmailAddress,VisitDate) values ('%s', '%s', '%s',NOW())", firstName, lastName, emailAddress)
	fmt.Println(insertQry)
	insert, err := db.Query(insertQry)
	if err != nil {
		panic(err.Error())
	}
	insert.Close()
	return u, nil
}

// GetUserByID Not in use yet
func GetUserByID (id int) (User, error) {
	//db := dataconnection.dbconn()
	for _, u := range users {
		if u.ID == id {
			return *u, nil
		}
	}
	return User{}, fmt.Errorf("User with ID '%v' not  found", id)
}

// UpdateUser Not in use yet
func UpdateUser(u User)	(User, error) {
	for i, candidate := range users {
		if candidate.ID == u.ID	{
			users[i] = &u
			return u, nil
		}
	}
	return User{}, fmt.Errorf("User with ID '%v' not found", u.ID)
}

func RemoveUserByID(id int) error {
	//delete from visitors where ID = 8
	for i, u := range users {
		if u.ID == id {
			users = append(users[:i], users[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("User with ID '%v' not  found", id)
}

/*
//var u []User
	///var userId int
	//var firstName, lastName, emailAddress, visitDate string


	for results.Next() {
		err = results.Scan(&userId,&firstName,&lastName,&emailAddress,&visitDate)
		fmt.Println(results)
		if err != nil {
			panic(err.Error())
		}
		u = append(u, User{ID: userId, FirstName: firstName, LastName: lastName, EmailAddress: emailAddress, VisitDate: visitDate}  ) //u = append(results)   //u, results)
	}
	return u
	/*
 */