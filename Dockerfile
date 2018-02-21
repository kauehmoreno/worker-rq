FROM golang:1.9

# COPY . /go/src/github.com/kauehmoreno/worker-rq
# WORKDIR /go/src/github.com/kauehmoreno/worker-rq

# RUN go get ./
# RUN go build -o main main.go redis.go s3Bucket.go

# CMD [ "go run main.go redis.go s3Bucket.go" ]


WORKDIR /go/src/github.com/kauehmoreno/worker-rq
COPY . .

ENV  S3_BUCKET_NAME jacortei-login
ENV S3_BUCKET_URL https://ja-cortei.nyc3.digitaloceanspaces.com

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["worker-rq"]