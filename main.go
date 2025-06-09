package main

import (
	"fmt"
	// "log"
	"log"
	"time"

	// "github.com/redis/go-redis/v9"
	"distributed-task-queue/producers"
	"distributed-task-queue/workers"
	"distributed-task-queue/internal"
	"distributed-task-queue/queue"
		"github.com/joho/godotenv"
		"strconv"
		"os"
	// "sort"
)



func main() {
	fmt.Println("Welcome to a Simple Email Sending Distributed Task System")

	// get inputs from user via cli
	from,subject,message,receiver,priority,_:=internal.GetInput()

	// start a new redis client
	// client := redis.NewClient(&redis.Options{
	// 	Addr:     "localhost:6379",
	// 	Password: "", 
	// 	DB:       0,  
	// })

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	password := os.Getenv("EMAIL_HOST_PASSWORD")
	emailUser := os.Getenv("EMAIL_HOST_USER")
	emailHost := os.Getenv("EMAIL_HOST")
	smtpPort := os.Getenv("EMAIL_PORT")
	smtpPortInt,err:=strconv.Atoi(smtpPort)
	client,err:=internal.ConnectRD("localhost","6379","",0)
	dialer,err:=internal.ConnectMail(emailHost,smtpPortInt,emailUser,password)
	if err != nil {
		log.Println(err)
		return
	}
	queuelist:=queue.QueueList{}
	email := internal.Mail{
		From: from,
		Subject:  subject,
		Message:  []byte(message),
		Receiver: []string{receiver},
	}

	if priority != 0 {
		value:=queue.Queue{
			Payload: email,
			Priority: priority,
			DateCreated: time.Now(),
		}

		queuelist.Enqueue(value)
		// producer function with error handling
	if err := producers.PriorityProducer(client, queuelist); err != nil {
		log.Println(err)
		return
	}

	// Worker function with error handling
	if err := workers.PriorityWorker(client,5,dialer); err != nil {
		log.Println(err)
		return
	}
	}else{
		if err := producers.Producer(client, &email); err != nil {
			log.Println(err)
			return
		}
	
		// Wotker function with error handling
		if err := workers.Worker(client,5,dialer); err != nil {
			log.Println(err)
			return
		}
	}


}
