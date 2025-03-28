package main

import (
	"fmt"
	"log"
	"github.com/redis/go-redis/v9"
	"distributed-task-queue/producers"
	"distributed-task-queue/workers"
	"distributed-task-queue/internal"
)



func main() {
	fmt.Println("Welcome to a Simple Email Sending Distributed Task System")

	// get inputs from user via cli
	subject,message,receiver,_:=internal.GetInput()

	// start a new redis client
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", 
		DB:       0,  
	})

	email := internal.Mail{
		Subject:  subject,
		Message:  []byte(message),
		Receiver: []string{receiver},
	}

	// producer function with error handling
	if err := producers.Producer(client, &email); err != nil {
		log.Println(err)
		return
	}

	// Wotker function with error handling
	if err := workers.Worker(client,5); err != nil {
		log.Println(err)
		return
	}

}
