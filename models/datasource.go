package models

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

var DataSource string

func SetDatasource()  {
	//os.e
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	SQLSERVER := os.Getenv("SQLSERVER")
	SQLUSER := os.Getenv("SQLUSER")
	SQLPW := os.Getenv("SQLPW")
	SQLPORT := os.Getenv("SQLPORT")
	SQLDBNAME := os.Getenv("SQLDBNAME")
	DataSource = SQLUSER + ":" + SQLPW + "@tcp(" + SQLSERVER + ":" + SQLPORT + ")/"+ SQLDBNAME + "?parseTime=true"
}

//var DataSource = sqlconnectionstring
