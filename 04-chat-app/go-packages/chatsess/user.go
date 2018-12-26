package chatsess

import (
	"fmt"
	"html"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type User struct {
	Username string
	Password string
}

func NewUser(name, pass string) User {
	return User{
		Username: html.EscapeString(name),
		Password: NewPassword(pass),
	}
}

func (u User) Put(sess *session.Session) error {
	cdb := dynamodb.New(sess)
	_, err := cdb.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("ch_users"),
		Item: map[string]*dynamodb.AttributeValue{
			"username": {S: aws.String(u.Username)},
			"password": {S: aws.String(u.Password)},
		},
	})

	return err
}

func (u User) Validate(passw string) error {
	if len(u.Username) < 3 {
		return fmt.Errorf("Invalid username.")
	}
	if len(passw) < 6 {
		return fmt.Errorf("Invalid password too small.")
	}
	return nil
}

func GetDBUser(uname string, sess *session.Session) (User, error) {
	cdb := dynamodb.New(sess)
	dbu, err := cdb.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("ch_users"),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {S: aws.String(html.EscapeString(uname))},
		},
	})

	if err != nil {
		return User{}, err
	}

	if dbu.Item == nil {
		return User{}, fmt.Errorf("User doesn't exists.")
	}

	return User{Username: *(dbu.Item["username"].S), Password: *(dbu.Item["password"].S)}, nil
}

func GetDBUserPass(uname, pass string, sess *session.Session) (User, error) {
	u, err := GetDBUser(uname, sess)
	if err != nil {
		return u, err
	}

	if !CheckPassword(pass, u.Password) {
		return User{}, fmt.Errorf("Passwords do not match.")
	}
	return u, err
}
