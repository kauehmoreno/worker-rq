package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
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
		log.WithFields(log.Fields{
			"error": err.Error(),
			"time":  time.Now(),
		}).Fatal("Error on subscribe channel - it will break worker run")
	}

	for {
		msg, err := pubsub.ReceiveMessage()
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
				"time":  time.Now(),
			}).Error("Error on subscribe channel")
			// TODO think about store this image bucket into a queue if does not success send to s3
		}
		if marshalError := json.Unmarshal([]byte(msg.Payload), &img); marshalError != nil {
			// TODO same concept over - if fails what should we do with buffer? Image must be created on bucket anyway
			log.WithFields(log.Fields{
				"error": marshalError.Error(),
				"time":  time.Now(),
			}).Error("Error on subscribe channel")
		} else {
			go img.SendBucket()
		}
	}
}

// Consumer to send SMS to verification account
func (sms SMS) Consumer() {
	redis := GetInstance()
	pubsub, err := redis.Subscribe("sms")
	fmt.Println("Vai comecar a execucao")

	defer pubsub.Close()

	if err != nil {
		fmt.Println("Deu erro", err.Error())
		log.WithFields(log.Fields{
			"error": err.Error(),
			"time":  time.Now(),
		}).Fatal("Error on subscribe to channel")
		return
	}

	fmt.Println("Vai comecar o for")
	for {
		msg, erro := pubsub.ReceiveMessage()
		if erro != nil {
			log.WithFields(log.Fields{
				"error": erro.Error(),
				"time":  time.Now(),
			}).Error("Error on reading message from channel")
		}

		if erroToUnmarshal := json.Unmarshal([]byte(msg.Payload), &sms); erroToUnmarshal != nil {
			log.WithFields(log.Fields{
				"error": erroToUnmarshal.Error(),
				"time":  time.Now(),
				"msg":   msg.Payload,
			}).Error("Error on unmmarshal message")
		} else {
			go sms.SendSMS()
		}
	}
}

func (img ImageBucket) errorOnSendImg() {
	redis := GetInstance()
	names := strings.Split(img.FileName, "/")
	name := names[len(names)-1]
	key := fmt.Sprintf("ImageError:%s", name)
	redis.RPush(key, img)
}

// Set on Redis information about confirmation token
// which key is his email
func (sms *SMS) Set(expire time.Duration) {
	redis := GetInstance()
	erro := redis.Set(sms.Email, sms.ConfirmToken, expire).Err()
	if erro != nil {
		log.WithFields(log.Fields{
			"erro": erro.Error(),
			"time": time.Now(),
		}).Error("Error on set key into redis")
	}
}
