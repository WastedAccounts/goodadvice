package messaging

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"goodadvice/v1/datasource"
	"io"
	"log"
	"net/mail"
	"net/smtp"
	"strings"
)

type Email struct {
	recipientName  string
	recipientEmail string
	title          string
	body           string
	adminName      string
	adminEmail     string
	smtpServer     string
	smtpPort       string
	smtpUser       string
	smtpPW         string
}

func SendEmail(e Email) {
	//https://gist.github.com/andelf/5004821
	// Set up authentication information.
	auth := smtp.PlainAuth(
		"",
		e.smtpUser,
		e.smtpPW,
		e.smtpServer,
	)

	sender := mail.Address{e.adminName, e.adminEmail}
	recipient := mail.Address{e.recipientName, e.recipientEmail}

	header := make(map[string]string)
	header["From"] = sender.String()
	header["To"] = recipient.String()
	header["Subject"] = e.title //encodeRFC2047(title)
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(e.body))

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err := smtp.SendMail(
		e.smtpServer+":"+e.smtpPort,
		auth,
		sender.Address,
		[]string{recipient.Address},
		[]byte(message), //[]byte("This is the email body."),
	)
	if err != nil {
		log.Fatal(err)
	}
}

func VerificationEmail(newuid int64) {
	// Vars
	var e Email            // Email struct
	var regVal string      // for each value is registry query result
	var regValues []string // array of registry results to feed into Email struct

	// Generate confirmation code
	code := GenerateConfCode(6)

	// write code to db
	_, err := datasource.DBconn.Exec("INSERT INTO email_verification (userid, verification_code, expires) VALUE (?,?,UTC_TIMESTAMP() + INTERVAL 10 MINUTE)", newuid, code)
	if err != nil {
		// If there is any issue with inserting into the database, return a 500 error
		panic(err.Error())
	}
	//// Get code ID from verification table to use for verifying users code
	//_, err = writeCode.LastInsertId()
	if err != nil {
		// If there is any issue with inserting into the database, return a 500 error
		panic(err.Error())
	}
	// Query user table for user info
	getUserInfo, err := datasource.DBconn.Query("SELECT username,emailaddress FROM users WHERE ID = ?", newuid)
	defer getUserInfo.Close()
	if err != nil {
		// If there is any issue with inserting into the database, return a 500 error
		panic(err.Error())
	}
	for getUserInfo.Next() {
		err := getUserInfo.Scan(&e.recipientName, &e.recipientEmail)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Query for Admin and SMTP information
	getAdminInfo, err := datasource.DBconn.Query("SELECT value FROM registry WHERE name IN ('SMTP_SERVER','SMTP_USER','SMTP_PW','SMTP_PORT','ADMIN_NAME','ADMIN_EMAIL') ORDER BY name DESC;")
	defer getAdminInfo.Close()
	if err != nil {
		// If there is any issue with inserting into the database, return a 500 error
		panic(err.Error())
	}
	for getAdminInfo.Next() {
		err := getAdminInfo.Scan(&regVal)
		if err != nil {
			log.Fatal(err)
		}
		regValues = append(regValues, regVal)
	}
	//getUserInfo.Close()

	// load registry results from array into struct
	e.smtpUser = regValues[0]
	e.smtpServer = regValues[1]
	e.smtpPW = regValues[2]
	e.smtpPort = regValues[3]
	e.adminName = regValues[4]
	e.adminEmail = regValues[5]

	// Query for Message format
	getEmailMessage, err := datasource.DBconn.Query("SELECT message_subject,message_body FROM message_templates WHERE message_name = 'VerificationEmail';")
	defer getEmailMessage.Close()
	if err != nil {
		// If there is any issue with inserting into the database, return a 500 error
		panic(err.Error())
	}
	for getEmailMessage.Next() {
		err := getEmailMessage.Scan(&e.title, &e.body)
		if err != nil {
			log.Fatal(err)
		}
	}
	//getEmailMessage.Close()

	// load confirmation code into email body
	e.body = fmt.Sprintf(e.body, code)

	// Pass infor to SendEmail function so we can send an email
	SendEmail(e)

}

// Create 6 digit code for verification
func GenerateConfCode(max int) string {
	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

//encoding function - I guess I don't need it?
func encodeRFC2047(String string) string {
	// use mail's rfc2047 to encode any string
	addr := mail.Address{String, ""}
	return strings.Trim(addr.String(), " <>")
}

//old way
//https://www.loginradius.com/blog/async/sending-emails-with-golang/
//func SendConfirmationEmail1() {
//	// Sender data.
//	from := "moonshotlimited@gmail.com"
//	password := "tzkipjqayeobhmyu"
//
//	// Receiver email address.
//	to := []string{
//		"matthewjaysimpson@gmail.com",
//	}
//
//	// smtp server configuration.
//	smtpHost := "smtp.gmail.com"
//	smtpPort := "587"
//
//	// Message.
//	message := []byte("This is a test email message.")
//
//	// Authentication.
//	auth := smtp.PlainAuth("", from, password, smtpHost)
//
//	// Sending email.
//	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Println("Email Sent Successfully!")
//}
