package producers

import (
	"context"
	"distributed-task-queue/queue"
	"encoding/json"
	"log"
	"strconv"
	"math/rand"
	"github.com/redis/go-redis/v9"
)


func PriorityProducer(client *redis.Client,values queue.QueueList)error{
	ctx:=context.Background()
	length:=len(values)

	for i:=0;i<length;i++{
		jsonmail,err:=json.Marshal(values[i].Payload)
		if err !=nil{
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

		err=client.HSet(
			ctx,
			"task:"+task_id,
			"payload",jsonmail,
			"priority",values[i].Priority,
			"retries",0,
			"status","pending",
			"date_created",values[i].DateCreated,

		).Err()

		if err != nil {
			log.Println(err)
		}
		err=client.LPush(ctx,"priority_queue",task_id).Err()
		if err !=nil{
			log.Println(err)
		}
	}
	return nil
}