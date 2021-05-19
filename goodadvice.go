package main

import (
	"fmt"
	"goodadvice/v1/appinit"
	"goodadvice/v1/controllers"
	"net/http"
)



func main()  {
	fmt.Println("Server is starting")
	//models.SetDatasource()
	appinit.Init()
	controllers.RegisterControllers()
	http.ListenAndServe(":3000", nil)
	fmt.Println("Server has started")
}


