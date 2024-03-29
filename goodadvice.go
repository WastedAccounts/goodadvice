package main

import (
	"fmt"
	"goodadvice/v1/appinit"
	"goodadvice/v1/controllers"
	"goodadvice/v1/datasource"
	"net/http"
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
	http.ListenAndServe(":3000", nil)
	fmt.Println("Server has stopped")
	defer datasource.DBconn.Close()
}

//func emailVerification() {
//	// https://medium.com/@stoyanov.veseline/self-hosting-a-mail-server-in-2019-6d29542dadd4
//	// https://mailcow.email/
//	// https://letsencrypt.org/getting-started/
//	// https://www.siteground.com/kb/gmail-smtp-server/?gclid=Cj0KCQjwkZiFBhD9ARIsAGxFX8CCkYj5_pWvqk2r5TIOIiMQLuUbU2bT_WK-44BqYcq2oQ9f7-muZswaAsUhEALw_wcB
//	// Set up authentication information.
//	fmt.Println("Email about to send")
//	auth := smtp.PlainAuth("", "matthewjaysimpson@gmail.com", "Ninjals124!", "smtp.gmail.com")
//
//	// Connect to the server, authenticate, set the sender and recipient,
//	// and send the email all in one step.
//	to := []string{"matthewjaysimpson@gmail.com"}
//	msg := []byte("To: matthewjaysimpson@gmail.com\r\n" +
//		"Subject: discount Gophers!\r\n" +
//		"\r\n" +
//		"This is the email body.\r\n")
//	err := smtp.SendMail("mail.example.com:25", auth, "sender@example.org", to, msg)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("Email sent")
//}
