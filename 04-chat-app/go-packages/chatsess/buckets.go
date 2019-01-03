package chatsess

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const BUCKET = "public.sampalm.com"

type Object struct {
	Key  string
	Size int64
	Body []byte
}

func ListObjects(username string, sess *session.Session) ([]Object, error) {
	cs3 := s3.New(sess)

	res, err := cs3.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(BUCKET),
		Prefix: aws.String(fmt.Sprintf("%s/", username)),
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
		return []Object{}, err
	}
	log.Println("Response Body:\n", res.Contents)
	keys := []Object{}
	for _, v := range res.Contents {
		s := strings.SplitN(*(v.Key), "/", 2)[1]
		if s == "" {
			continue
		}
		keys = append(keys, Object{Key: s, Size: *(v.Size)})
	}
	return keys, nil
}

func CreateUserFolder(username string, sess *session.Session) error {
	cs3 := s3.New(sess)
	_, err := cs3.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(BUCKET),
		Key:    aws.String(fmt.Sprintf("%s/welcome.txt", username)),
		Body:   bytes.NewReader([]byte(fmt.Sprintf("Welcome %s!\nThis is your bucket, all of your items will be stored here.", username))),
	})
	if err != nil {
		return err
	}
	return nil
}

func (obj Object) Put(sess *session.Session) error {
	cs3 := s3.New(sess)
	_, err := cs3.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(BUCKET),
		Key:    aws.String(obj.Key),
		Body:   bytes.NewReader(obj.Body),
	})
	if err != nil {
		return err
	}
	return nil
}

func DownloadObject(username, filename string, sess *session.Session) (ContentType, Body string, err error) {
	cs3 := s3.New(sess)
	res, err := cs3.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(BUCKET),
		Key:    aws.String(fmt.Sprintf("%s/%s", username, filename)),
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
		return "", "", err
	}
	log.Println("Response Metadata:\n", res.Metadata, res.VersionId, res.ContentType)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err.Error())
		return "", "", err
	}
	ef := base64.StdEncoding.EncodeToString(body)
	return *(res.ContentType), ef, nil
}
