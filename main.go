package main

import (
	"fmt"
	"log"
	// "log"
	"time"

	// "github.com/redis/go-redis/v9"
	// "distributed-task-queue/producers"
	// "distributed-task-queue/workers"
	"distributed-task-queue/internal"
	"distributed-task-queue/queue"
	"sort"
)



func main() {
	fmt.Println("Welcome to a Simple Email Sending Distributed Task System")

	// get inputs from user via cli
	// subject,message,receiver,_:=internal.GetInput()

	// start a new redis client
	// client := redis.NewClient(&redis.Options{
	// 	Addr:     "localhost:6379",
	// 	Password: "", 
	// 	DB:       0,  
	// })

	// email := internal.Mail{
	// 	Subject:  subject,
	// 	Message:  []byte(message),
	// 	Receiver: []string{receiver},
	// }

	test_mail:=internal.Mail{
		Subject: "Test mail",
		Message: []byte("Test message"),
		Receiver: []string{"silaskumi4@gmail.com"},
	}

	// value:=queue.Queue{
	// 	Payload: test_mail,
	// 	Priority: 1,
	// 	DateCreated: time.Now(),
	// }

	queuelist:=queue.QueueList{}
	values:=[]int{1,2,3,4,5}
	for _,k:= range values{
		value:=queue.Queue{
			Payload: test_mail,
			Priority: k,
			DateCreated: time.Now(),
		}

		queuelist.Enqueue(value)
	}
	// queuelist.Enqueue(value)
	
	sort.Slice(queuelist,func(i,j int)bool{
		return queuelist[i].Priority > queuelist[j].Priority
	})
	fmt.Println(queuelist)

	
	last,err:=queuelist.Dequeue()
	if err !=nil {
		fmt.Println("Printing last value...")
		
	}
	fmt.Print(last)
	queuelist.Remove()
	last,err=queuelist.Dequeue()
	if err !=nil {
		log.Println(err)
		
	}
	fmt.Print(last)
	// producer function with error handling
	// if err := producers.Producer(client, &email); err != nil {
	// 	log.Println(err)
	// 	return
	// }

	// Wotker function with error handling
	// if err := workers.Worker(client,5); err != nil {
	// 	log.Println(err)
	// 	return
	// }

}
