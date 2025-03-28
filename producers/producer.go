package producers



import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	// "net/smtp"
	// "os"
	// "bufio"
	// "github.com/google/uuid"
	// "github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"strconv"
	"time"
	"distributed-task-queue/internal"
	// "strings"
)

type Mail struct {
	Subject  string   `json:"subject"`
	Message  []byte   `json:"message"`
	Receiver []string `json:"receiver"`
}


func Producer(client *redis.Client, mail *internal.Mail) error {
	ctx := context.Background()
	jsonmail, err := json.Marshal(mail)

	if err != nil {
		log.Println(err)
	}
	var task_id string
	for {
		random := rand.Intn(10)
		task_id += strconv.Itoa(random)
		if len(task_id) == 5 {
			break
		}
	}
	err = client.HSet(
		ctx,
		"task:"+task_id,
		"payload", jsonmail,
		"retries", 0,
		"status", "pending",
		"created_at", time.Now().Unix(),
	).Err()

	if err != nil {
		log.Println(err)
	}
	err = client.LPush(ctx, "tasks_queue", task_id).Err()
	// err = client.LPush(ctx, "tasks", jsonmail).Err()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("Task added to queue...")

	return nil
}