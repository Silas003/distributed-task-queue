package mechanism

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"github.com/redis/go-redis/v9"
	"distributed-task-queue/internal"
)




func MarkCompleted(task_id string,client *redis.Client)error{
	ctx:=context.Background()
	err:=client.HSet(
		ctx,
		"task:"+ task_id,
		"status","completed",
		"completed_at",time.Now().Unix(),
	).Err()
	if err!=nil{
		log.Println(err)
		return err
	}else{
		client.LRem(ctx,"processing_tasks",1,task_id)
	}

	return nil
}

func MarkFailed(task_id string,client *redis.Client)error{
	ctx:=context.Background()
	err:=client.HSet(
		ctx,
		"task:"+ task_id,
		"status","failed",
		"failed_at",time.Now().Unix(),
	).Err()
	if err!=nil{
		log.Println(err)
		return err
	}else{
		client.LPush(ctx,"dead_letter",task_id)
		client.LRem(ctx,"processing_tasks",1,task_id)
	}
	return nil
}


func ViewDeadLetter(client *redis.Client) ([]internal.Mail,error){
	ctx:=context.Background()
	var list []Mail
	for{
		task_id,err:=client.RPop(ctx,"dead_letter").Result()
		if err !=nil{
			log.Println(err)
		}
		taskData,err:=client.HGetAll(ctx,"task:"+task_id).Result()
		if err !=nil{
			log.Print(err)
		}

		var mail Mail
		err=json.Unmarshal([]byte(taskData["payload"]), &mail)
		if err !=nil{
            log.Print(err)
			return nil, err
        }
		list = append(list, mail)
	}
	return list,nil

}


// 
func ProcessRetry(task_id string,current_retries int,client *redis.Client) error{
	ctx:=context.Background()
	_, err := client.HIncrBy(ctx, "task:"+task_id, "retries", 1).Result()
    if err != nil {
        return fmt.Errorf("failed to increment retries: %w", err)
    }
	err=client.HSet(
		ctx,
		"task:"+task_id,
		"status","retrying",
		"last_attempt",time.Now().Unix(),
	).Err()
	if err !=nil{
        log.Println(err)
    }

	err = client.LPush(ctx, "task_queue", task_id).Err()
    if err != nil {
        return fmt.Errorf("failed to requeue task: %w", err)
    }

    return nil
}


