package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"gopkg.in/redis.v5"
)

var once sync.Once
var redisInstance *redis.Client

//GetInstance return a unique redis instances troughtout whole application
func GetInstance() *redis.Client {
	once.Do(func() {
		redisInstance = redis.NewClient(&redis.Options{
			Addr:     "queue:6379",
			Password: "", // no password was set
			DB:       0,  // use default DB
		})
	})
	return redisInstance
}

// Consumer will get any message on channel and send it to S3 bucket - image store
func (img ImageBucket) Consumer() {
	redis := GetInstance()
	pubsub, err := redis.Subscribe("bucket")

	defer pubsub.Close()

	if err != nil {
		fmt.Println(err.Error())
	}

	for {
		msg, err := pubsub.ReceiveMessage()
		if err != nil {
			fmt.Println(err.Error())
			// TODO think about store this image bucket into a queue if does not success send to s3
		}

		if marshalError := json.Unmarshal([]byte(msg.Payload), &img); marshalError != nil {
			// TODO same concept over - if fails what should we do with buffer? Image must be created on bucket anyway
			fmt.Println(marshalError.Error())
		}

		go img.SendBucket()
		fmt.Println(img)
	}

}
