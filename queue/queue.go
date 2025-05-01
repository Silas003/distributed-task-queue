package queue

import (
	"distributed-task-queue/internal"
	"time"
	// "sort"
	"fmt"
)


type Queue struct{
	Payload internal.Mail
	Priority int
	DateCreated time.Time
}


type QueueList []Queue


func (value *QueueList) Enqueue(task Queue) (QueueList,error){
	list:=append(*value,task)

	*value= list
	return *value,nil
}

func (value *QueueList) Dequeue()(any,error){
	if value == nil || len(*value) == 0 {
        return nil, fmt.Errorf("queue is empty or nil")
    }
	item := (*value)[len(*value)-1]
    
    // Remove last element
    // *value = (*value)[:len(*value)-1]
    
    return item, nil
}

func (value *QueueList) Remove()(any,error){
	if value == nil || len(*value) == 0 {
        return nil, fmt.Errorf("queue is empty or nil")
    }
	// item := (*value)[len(*value)-1]
    
    // Remove last element
    *value = (*value)[:len(*value)-1]
    
    return value, nil
}




func (value *QueueList) Pop()(any,error){
	if value == nil || len(*value) == 0 {
		return nil,fmt.Errorf("Queue is empty")
	}

	item := (*value)[len(*value)-1]

	return item,nil
}


func (value *QueueList) IsEmpty()bool{
	if len(*value) > 0 {
		return true
	}
	return false
}