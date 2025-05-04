package workers

import (
	"context"
	"distributed-task-queue/internal"
	"distributed-task-queue/mechanism"
	"encoding/json"
	"log"
	"time"
	"strconv"
	"github.com/redis/go-redis/v9"
)

func PriorityWorker(client *redis.Client, maxRetries int) error{
	ctx:=context.Background()


	for {
		taskId,err := client.BRPopLPush(
			ctx,
			"priority_queue",
			"priority_processing",
			time.Second * 30,
		).Result()

		if err != nil {
			log.Println(err)
		}

		taskData,err:=client.HGetAll(ctx,"task:"+taskId).Result()

		retries, _ := strconv.Atoi(taskData["retries"])

		if retries > maxRetries {

			mechanism.MarkFailed(taskId, client)
		}

		if err != nil {
			log.Println(err)
		}

		var mail internal.Mail
		if err=json.Unmarshal([]byte(taskData["payload"]),&mail); err != nil{
			log.Println(err)
		}

		if err = internal.SendMail(&mail); err != nil{
			mechanism.ProcessRetry(taskId,retries,client)
			log.Println(err)
		}
		mechanism.MarkCompleted(taskId,client)
	}
	return nil
}