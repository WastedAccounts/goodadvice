package main

import (
	"fmt"
	"goodadvice/v1/appinit"
	"goodadvice/v1/controllers"
	"goodadvice/v1/datasource"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Server is starting")
	appinit.Init()
	fmt.Println("App initialization complete")
	controllers.RegisterControllers()
	//// For testing
	//controllers.RegisterControllersTest()
	fmt.Println("Controllers registered")
	http.Handle("/assets/css/", http.StripPrefix("/assets/css/", http.FileServer(http.Dir("assets/css"))))
	fmt.Println("Serving asset files complete")
	err := http.ListenAndServeTLS(":3443", os.Getenv("WEBCERT"), os.Getenv("WEBKEY"), nil)
	if err != nil {
		fmt.Println(err)
		// this is for troubleshooting my container
		http.ListenAndServe(":3000", nil)
	}
	fmt.Println("Server has stopped")
	defer datasource.DBconn.Close()
}
