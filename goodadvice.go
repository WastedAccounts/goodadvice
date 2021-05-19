package main

import (
	"fmt"
	"mjs/v1/appinit"
	"mjs/v1/controllers"
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


