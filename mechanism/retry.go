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




func MarkCompleted(taskId string,client *redis.Client)error{
	ctx:=context.Background()
	taskData,err:=client.HGetAll(ctx,"task:"+taskId).Result()
	if err != nil {
		log.Println(err)
	}
	priority:= taskData["priority"]
	err=client.HSet(
		ctx,
		"task:"+ taskId,
		"status","completed",
		"completed_at",time.Now().Unix(),
	).Err()
	if err!=nil{
		log.Println(err)
		return err
	}else{
		if priority != ""{
			client.LRem(ctx,"priority_processing",1,taskId)
		}else{
			client.LRem(ctx,"processing_tasks",1,taskId)
		}
	}

	return nil
}

func MarkFailed(taskId string,client *redis.Client)error{
	ctx:=context.Background()
	taskData,err:=client.HGetAll(ctx,"task:"+taskId).Result()
	if err != nil {
		log.Println(err)
	}
	priority:=taskData["priority"]
	err=client.HSet(
		ctx,
		"task:"+ taskId,
		"status","failed",
		"failed_at",time.Now().Unix(),
	).Err()
	if err!=nil{
		log.Println(err)
		return err
	}else{
		client.LPush(ctx,"dead_letter",taskId)
		if priority != ""{
			client.LRem(ctx,"processing_tasks",1,taskId)
		}else{
			client.LRem(ctx,"priority_processing",1,taskId)
		}
	}
	return nil
}


func ViewDeadLetter(client *redis.Client) ([]internal.Mail,error){
	ctx:=context.Background()
	var list []internal.Mail
	for{
		taskId,err:=client.RPop(ctx,"dead_letter").Result()
		if err !=nil{
			log.Println(err)
		}
		taskData,err:=client.HGetAll(ctx,"task:"+taskId).Result()
		if err !=nil{
			log.Print(err)
		}

		var mail internal.Mail
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
func ProcessRetry(taskId string,current_retries int,client *redis.Client) error{
	ctx:=context.Background()
	_, err := client.HIncrBy(ctx, "task:"+taskId, "retries", 1).Result()
    if err != nil {
        return fmt.Errorf("failed to increment retries: %w", err)
    }
	err=client.HSet(
		ctx,
		"task:"+taskId,
		"status","retrying",
		"last_attempt",time.Now().Unix(),
	).Err()
	if err !=nil{
        log.Println(err)
    }
	taskData,err:=client.HGetAll(ctx,"task:"+taskId).Result()
	if err != nil{
		log.Println(err)
	}

	priority:=taskData["priority"]
	if priority !="" {
		err = client.LPush(ctx, "priority_queue", taskId).Err()
		if err != nil {
			return fmt.Errorf("failed to requeue task: %w", err)
		}
	}else{
		err = client.LPush(ctx, "task_queue", taskId).Err()

    	if err != nil {
        return fmt.Errorf("failed to requeue task: %w", err)
    }
	}
	

    return nil
}


