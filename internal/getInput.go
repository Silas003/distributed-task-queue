package internal

import (
	"bufio"
	"fmt"
	"strings"
	"os"
	"log"
	"github.com/asaskevich/govalidator"

)

func Isvalidmail(addr string) bool {
	// add more domains and make it more robust using regex
	return govalidator.IsEmail(addr)
}

func GetInput()(string,string,string,error){
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
	fmt.Println("Enter email address of recipient:	")

	receiver, _ = reader.ReadString('\n')
	receiver = strings.TrimSpace(receiver)
	if Isvalidmail(receiver) {
		log.Print()
	} else {
		log.Fatalf("Invalid email address %s", receiver)
		

	}
	return subject,message,receiver,nil
}