package internal

import (
	"fmt"
	"log"
	// "net/smtp"
	"os"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
	"strings"
)

type Mail struct {
	From 	string `json:"from"`
	Subject  string   `json:"subject"`
	Message  []byte   `json:"message"`
	Receiver []string `json:"receiver"`
}



func SendMail(mail *Mail,dialer *gomail.Dialer) error {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	password := os.Getenv("EMAIL_HOST_PASSWORD")
	emailUser := os.Getenv("EMAIL_HOST_USER")
	emailHost := os.Getenv("EMAIL_HOST")
	smtpPort := os.Getenv("EMAIL_PORT")

	if emailUser == "" || password == "" || emailHost == "" || smtpPort == "" {
		return fmt.Errorf("missing required email configuration")
	}


	m := gomail.NewMessage()
    m.SetHeader("From", mail.From)
    m.SetHeader("To", strings.Join(mail.Receiver,""))
    m.SetHeader("Subject", mail.Subject)
    m.SetBody("text/plain", string(mail.Message))
    // d := gomail.NewDialer(emailHost, 587, emailUser, password)
    if err := dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send mail: %v",err)
    }else{
		log.Printf("Mail Sent to %v",mail.Receiver)
	}

		return nil
}



