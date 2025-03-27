package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"os"

	// "String"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
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
func Producer(client *redis.Client, mail *Mail) error {
	ctx := context.Background()
	jsonmail, err := json.Marshal(mail)

	if err != nil {
		log.Println(err)
	}

	err = client.LPush(ctx, "tasks", jsonmail).Err()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("Task added to queue...")

	return nil
}

func Worker(client *redis.Client) error {
	ctx := context.Background()

	log.Println("Worker has started...")

	for  {

		task, err := client.RPop(ctx, "tasks").Result()
		if err== redis.Nil{
			fmt.Errorf("queue is empty %v",err)
			break
		}
		if err != nil {
			fmt.Println(err.Error())

			
		}
		var mail Mail

		if err := json.Unmarshal([]byte(task), &mail); err != nil {
			fmt.Printf("Error unmarshaling task: %v\n", err)
			continue
		}

		if err := SendMail(&mail); err != nil {
            log.Printf("Failed to send email: %v\n", err)
        }else{
			log.Printf("Mail sent to %v", mail.Receiver)
		}
		
	}

	return nil
}

func main() {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})


	email := Mail{
		Subject:  "Hello World",
		Message:  []byte("First email sending via golang"),
		Receiver: []string{"silaskumi4@gmail.com"},
	}


	if err:=Producer(client,&email);err!=nil{
		log.Println(err)
		return 
	}


	if err:=Worker(client);err!=nil{
		log.Println(err)
		return 
	}

	
}
