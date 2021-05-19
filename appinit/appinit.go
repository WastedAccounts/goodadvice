package appinit

import (
	"mjs/v1/models"
	"os"
)
//var Global_ID = ""

func Init() {
	os.Setenv("TZ", "America/New_York")
	models.SetDatasource()
}



