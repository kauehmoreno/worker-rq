package main

import (
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.WarnLevel)
}

func main() {
	var img ImageBucket
	var message SMS
	go img.Consumer()
	message.Consumer()
}
