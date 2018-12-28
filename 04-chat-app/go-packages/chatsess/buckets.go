package chatsess

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const BUCKET = "s3://bucketname"

func ListObjects(username string, sess *session.Session) ([]*s3.Object, error) {
	cs3 := s3.New(sess)

	res, err := cs3.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(BUCKET),
		Prefix: aws.String(fmt.Sprintf("/public-files/%s/", username)),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				log.Println(s3.ErrCodeNoSuchBucket, aerr.Error())
			default:
				log.Println(aerr.Error())
			}
		} else {
			log.Println(err.Error())
		}
		return []*s3.Object{}, err
	}

	return res.Contents, nil
}
