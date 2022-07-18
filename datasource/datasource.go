package datasource

import (
	"database/sql"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var (
	DBconn *sql.DB
    DataSource string
)

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
}
