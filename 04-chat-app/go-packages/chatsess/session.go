package chatsess

import (
	"crypto/rand"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/aws/aws-sdk-go/aws/session"
)

type Login struct {
	Sessid   string `db:"sessid"`
	Username string `db:"username"`
}

func NewLogin(name string) Login {
	bs := make([]byte, 20)
	rand.Read(bs)
	return Login{
		Sessid:   fmt.Sprintf("%x", bs),
		Username: name,
	}
}

func GetLogin(id string, sess *session.Session) (Login, error) {
	cdb := dynamodb.New(sess)
	res, err := cdb.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("ch_sessions"),
		Key: map[string]*dynamodb.AttributeValue{
			"sessid": {S: aws.String(id)},
		},
	})
	if err != nil {
		return Login{}, err
	}
	if _, ok := res.Item["username"]; !ok {
		return Login{}, fmt.Errorf("No username found")
	}
	return Login{Sessid: res.Item["sessid"].GoString(), Username: res.Item["username"].GoString()}, nil
}

func (l Login) Put(sess *session.Session) error {
	cdb := dynamodb.New(sess)
	_, err := cdb.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("ch_sessions"),
		Item: map[string]*dynamodb.AttributeValue{
			"sessid":   {S: aws.String(l.Sessid)},
			"username": {S: aws.String(l.Username)},
		},
	})
	return err
}
