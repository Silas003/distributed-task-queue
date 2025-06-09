package dtq

import (
	"fmt"

	"github.com/Silas003/distributed-task-queue/internal"
	"github.com/Silas003/distributed-task-queue/producers"
	"github.com/Silas003/distributed-task-queue/queue"
	"github.com/Silas003/distributed-task-queue/workers"
	"github.com/redis/go-redis/v9"
	"gopkg.in/gomail.v2"
)

func ConnectRD(host string, port string, password string, db int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       0,
	})

	return client, nil
}

func ConnectMail(smptHost string, smptPort int, emailUser string, password string) (*gomail.Dialer, error) {
	if emailUser == "" || password == "" || smptHost == "" || smptPort == 0 {
		return nil, fmt.Errorf("missing required email configuration")
	}
	dialer := gomail.NewDialer(smptHost, smptPort, emailUser, password)

	return dialer, nil
}

type Mail *internal.Mail
type QueueList *queue.QueueList

func NewMail(from string, subject string, message []byte, receiver []string) internal.Mail {
	mail := internal.Mail{
		From:     from,
		Subject:  subject,
		Message:  message,
		Receiver: receiver,
	}
	return mail
}

func Producer(client *redis.Client, mail *internal.Mail) error {
	return producers.Producer(client, mail)
}


func Worker(client *redis.Client, dialer *gomail.Dialer) error {
	return workers.Worker(client, 5, dialer)
}
