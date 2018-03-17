package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	minio "github.com/minio/minio-go"
	log "github.com/sirupsen/logrus"
)

type ImageBucket struct {
	Image      string `json:"image"`
	BucketName string `json:"bucketName"`
	FileName   string `json:"fileName"`
	Extension  string `json:"extension"`
	Kind       string `json:"kind"`
}

// SendBucket is method responsable to put object into s3 bucket
func (img ImageBucket) SendBucket() {

	client, err := minio.New(os.Getenv("S3_BUCKET_URI"), os.Getenv("S3_BUCKET_CLIENT_ID"), os.Getenv("S3_BUCKET_SECRET"), true)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
			"time":  time.Now(),
		}).Error("Error on setup s3 digitalocean client")
		return
	} else {

		decoded, decodeErr := base64.StdEncoding.DecodeString(img.Image)
		if decodeErr != nil {
			log.WithFields(log.Fields{
				"error": decodeErr.Error(),
				"time":  time.Now(),
			}).Error("Error decode string base64 to bytes")
		}

		buf := bytes.NewReader(decoded)

		err := client.MakeBucket("ja-cortei-user", "us-east-1")
		if err != nil {
			// Check to see if we already own this bucket (which happens if you run this twice)
			exists, err := client.BucketExists("ja-cortei-user")
			if err == nil && exists {
				log.WithFields(log.Fields{
					"time": time.Now(),
				}).Info("We already own %s\n", "ja-cortei-user")
			} else {
				log.WithFields(log.Fields{
					"time":  time.Now(),
					"error": err.Error(),
				}).Error("Error on makebucket on s3 digitalocean")
			}
		}

		log.WithFields(log.Fields{
			"time": time.Now(),
		}).Info("Successfully created %s\n", "ja-cortei-user")

		rules := make(map[string]string)
		rules["x-amz-acl"] = "public-read"

		resp, s3Error := client.PutObject(
			"ja-cortei-user",
			img.FileName,
			buf,
			-1,
			minio.PutObjectOptions{
				ContentType:  fmt.Sprintf("image/%s", img.Extension),
				UserMetadata: rules,
			})

		if s3Error != nil {
			log.WithFields(log.Fields{
				"time":     time.Now(),
				"error":    s3Error.Error(),
				"fileName": img.FileName,
			}).Error("Error on PutObject on s3 digitaocean bucket")
			img.errorOnSendImg()
		} else {
			log.WithFields(log.Fields{
				"time":     time.Now(),
				"quantity": resp,
				"fileName": img.FileName,
			}).Info("Image were success upload")
		}
	}
}
