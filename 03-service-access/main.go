package main

import (
	"bytes"
	"context"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
)

type Event struct {
	Text string
}
type Response struct {
	Msg string
	Err string
}

func handler(c context.Context, e Event) (Response, error) {
	sess := session.Must(session.NewSession())

	cs3 := s3.New(sess)
	_, err := cs3.PutObjectWithContext(c, &s3.PutObjectInput{
		Bucket: aws.String("sampalm"),
		Key:    aws.String("file01.txt"),
		Body:   bytes.NewReader([]byte("Hello " + e.Text)),
	})

	res := Response{Msg: "Uploading " + e.Text}
	if err != nil {
		res.Err = err.Error()
	}

	return res, nil
}

func main() {
	lambda.Start(handler)
}
