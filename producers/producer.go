package producers



import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"strconv"
	"time"
	"distributed-task-queue/internal"

)

// Mail is a simple struct to represent an email message.
type Mail struct {
	Subject  string   `json:"subject"`
	Message  []byte   `json:"message"`
	Receiver []string `json:"receiver"`
}


// Producer sends a mail to Redis queue.
func Producer(client *redis.Client, mail *internal.Mail) error {
	ctx := context.Background()
	jsonmail, err := json.Marshal(mail)

	if err != nil {
		log.Println(err)
	}
	var task_id string

	// Generate a unique task ID.
	for {
		random := rand.Intn(10)
		task_id += strconv.Itoa(random)
		if len(task_id) == 5 {
			break
		}
	}

	// Add task to Redis hash and list.
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
	
	// Add task ID to the task queue.
	err = client.LPush(ctx, "tasks_queue", task_id).Err()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("Task added to queue...")

	return nil
}