package appinit

import (
	"goodadvice/v1/datasource"
	"goodadvice/v1/models/old"
	"os"
)

//var Global_ID = ""

func Init() {
	os.Setenv("TZ", "America/New_York")
	old.SetDatasource()
	datasource.SetDatasource()
}
