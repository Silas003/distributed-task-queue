package internal

import (
	"bufio"
	"fmt"
	"strings"
	"os"
	"log"
	"github.com/asaskevich/govalidator"
	"strconv"

)

func Isvalidmail(addr string) bool {
	// add more domains and make it more robust using regex
	return govalidator.IsEmail(addr)
}

func GetInput()(string,string,string,int,error){
	reader := bufio.NewReader(os.Stdin)
	var message, subject, receiver string

	fmt.Println("Draft your mail Subject to recipient:	")
	subject, _ = reader.ReadString('\n')
	subject = strings.TrimSpace(subject)
	fmt.Println("Draft your message to recipient:	")
	for {
		line, _ := reader.ReadString('\n')
		if strings.TrimSpace(line) == "" {
			break
		}
		message += line
	}

	message = strings.TrimSpace(message)
	fmt.Println("Enter Priority of mail. Leave as 0 if it has no priority else \n 4:\"CRITICAL\" \n 3:\"HIGH\" \n 2:\"MEDIUM\" \n 1:\"LOW\"")
	priority,_:= reader.ReadString('\n')
	priority=strings.TrimSpace(priority)
	priorityInt,_:=strconv.Atoi(priority)
	fmt.Println("Enter email address of recipient:	")
	receiver, _ = reader.ReadString('\n')
	receiver = strings.TrimSpace(receiver)
	if !Isvalidmail(receiver){
		log.Printf("Invalid email address %s", receiver)
	} 
	return subject,message,receiver,priorityInt,nil
}