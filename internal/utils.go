package internal


import (
	"fmt"
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


func ConnectMail(smptHost string, smptPort int, emailUser string, password string)(*gomail.Dialer,error){
	if emailUser == "" || password == "" || smptHost == "" || smptPort == 0 {
		return nil, fmt.Errorf("missing required email configuration")
	}
	dialer := gomail.NewDialer(smptHost, smptPort, emailUser, password)

	return dialer, nil
}