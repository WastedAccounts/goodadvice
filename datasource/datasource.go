package datasource

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

var DataSource string
var db string

func SetDatasource() {
	// Get environment variables
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	SQLSERVER := os.Getenv("SQLSERVER")
	SQLUSER := os.Getenv("SQLUSER")
	SQLPW := os.Getenv("SQLPW")
	SQLPORT := os.Getenv("SQLPORT")
	SQLDBNAME := os.Getenv("SQLDBNAME")

	// Create Datasource string
	DataSource = SQLUSER + ":" + SQLPW + "@tcp(" + SQLSERVER + ":" + SQLPORT + ")/" + SQLDBNAME + "?parseTime=true"

	//// create datasource
	//db, err := sql.Open("mysql", DataSource)
	//if err != nil {
	//	panic(err.Error())
	//}
	//defer db.Close()
}
