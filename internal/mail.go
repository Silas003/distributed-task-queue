package internal

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"github.com/joho/godotenv"
)

type Mail struct {
	Subject  string   `json:"subject"`
	Message  []byte   `json:"message"`
	Receiver []string `json:"receiver"`
}



func SendMail(mail *Mail) error {

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

	smtpServer := fmt.Sprintf("%s:%s", emailHost, smtpPort)

	auth := smtp.PlainAuth(
		"",
		emailUser,
		password,
		emailHost,
	)

	// Send the email
	err := smtp.SendMail(
		smtpServer,
		auth,
		emailUser,
		mail.Receiver,
		mail.Message,
	)

	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}
	log.Println("Mail Sent!")

	return nil
}