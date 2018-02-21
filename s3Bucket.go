package main

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"os"
	"time"

	"github.com/minio/minio-go/pkg/policy"

	minio "github.com/minio/minio-go"
)

// S3_BUCKET_NAME = 'jacortei-login'
// S3_BUCKET_URL='https://ja-cortei.nyc3.digitaloceanspaces.com'
// S3_BUCKET_CLIENT_ID = 'BYHFFUOGOWHK3NG5JVTV'
// S3_BUCKET_SECRET = 'S1sVHID4eqNiGEPbtAvrRdmjwx/fdC6SOV0IphifLHw'

type ImageBucket struct {
	Image      string `json:"image"`
	BucketName string `json:"bucketName"`
	FileName   string `json:"fileName"`
	Extension  string `json:"extension"`
	Kind       string `json:"kind"`
}

func (img *ImageBucket) buildPath() string {
	t := time.Now()
	year := fmt.Sprintf("%d", t.Year())
	month := fmt.Sprintf("%02d", t.Month())

	shaYear := sha1.New()
	shaMonth := sha1.New()

	shaYear.Write([]byte(year))
	shaMonth.Write([]byte(month))

	newYear := fmt.Sprintf("%x\n", shaYear.Sum(nil))
	newMonth := fmt.Sprintf("%x\n", shaMonth.Sum(nil))

	return fmt.Sprintf("/static/%s/%x/%x/%s.%s", img.Kind, newYear, newMonth, img.FileName, img.Extension)
}

func (img *ImageBucket) SendBucket() string {
	accesKey := os.Getenv("S3_BUCKET_CLIENT_ID")
	secretKey := os.Getenv("S3_BUCKET_SECRET")

	client, err := minio.New(os.Getenv("S3_BUCKET_URL"), accesKey, secretKey, true)
	if err != nil {
		fmt.Println(err)
	}
	uri := img.buildPath()

	erro := client.SetBucketPolicy(
		os.Getenv("S3_BUCKET_NAME"),
		uri,
		policy.BucketPolicyReadWrite,
	)
	if erro != nil {
		fmt.Println(erro.Error())
	}

	buf := bytes.NewBufferString(img.Image)
	resp, s3Error := client.PutObject(
		os.Getenv("S3_BUCKET_NAME"),
		uri,
		buf,
		-1,
		minio.PutObjectOptions{
			ContentType: fmt.Sprintf("image/%s", img.Extension)})

	if s3Error != nil {
		fmt.Println(s3Error.Error())
	}

	return fmt.Sprintf("Upload %s quantity: %d Sucessfully", img.FileName, resp)
}
