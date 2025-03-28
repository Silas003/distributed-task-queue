package workers


import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"github.com/redis/go-redis/v9"
	"distributed-task-queue/internal"
	"distributed-task-queue/mechanism"
	"time"
	"strconv"
)

func Worker(client *redis.Client,maxRetries int) error {
	ctx := context.Background()

	log.Println("Worker has started...")

	for {

		task_id, err := client.BRPopLPush(
			ctx,
			"tasks_queue",
			"processing_tasks",
			30*time.Second,
		).Result()
		if err == redis.Nil {
			fmt.Errorf("queue is empty %v", err)
			break
		}
		if err != nil {
			fmt.Println(err.Error())

		}
		taskData, err := client.HGetAll(ctx, "task:"+task_id).Result()

		retries, _ := strconv.Atoi(taskData["retries"])

		if retries > maxRetries{

			mechanism.MarkFailed(task_id,client)
		}
		var mail internal.Mail

		if err := json.Unmarshal([]byte(taskData["payload"]), &mail); err != nil {
			log.Printf("Error unmarshaling task: %v\n", err)
			continue
		}

		if err := internal.SendMail(&mail); err != nil {
			mechanism.ProcessRetry(task_id,retries,client)
			log.Printf("Failed to send email: %v\n", err)
			break
		} else {
			mechanism.MarkCompleted(task_id,client)
			log.Printf("Mail sent to %v", mail.Receiver)
		}

	}

	return nil
}