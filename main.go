package main

import (
	"fmt"
	"log"
	"github.com/redis/go-redis/v9"
	"distributed-task-queue/mechanism"
	"distributed-task-queue/producers"
	"distributed-task-queue/workers"
	"distributed-task-queue/internal"
)



func main() {
	fmt.Println("Welcome to a Simple Email Sending Distributed Task System")

	subject,message,receiver,_:=internal.GetInput()
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	email := internal.Mail{
		Subject:  subject,
		Message:  []byte(message),
		Receiver: []string{receiver},
	}

	if err := producers.Producer(client, &email); err != nil {
		log.Println(err)
		return
	}

	if err := workers.Worker(client,5); err != nil {
		log.Println(err)
		return
	}

}
