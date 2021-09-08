package appinit

import (
	"database/sql"
	"goodadvice/v1/datasource"
	//"goodadvice/v1/models/old"
	"os"
)


func Init() {
	os.Setenv("TZ", "America/New_York")
	datasource.SetDatasource()
	//datasource.DBconn.SetConnMaxIdleTime(0)
	//datasource.DBconn.SetMaxOpenConns(100)
	var err error
	datasource.DBconn, err = sql.Open("mysql", datasource.DataSource)
	if err != nil {
		panic(err.Error())
	}
}
