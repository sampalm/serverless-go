package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/sampalm/serverless/04-chat-app/go-packages/chatsess"
)

type Event struct {
	Username string
	Password string
}

type Response struct {
	Value       int
	Description string
}

func handler(ctx context.Context, ev Event) (Response, error) {
	sess := session.Must(session.NewSession())

	_, err := chatsess.GetDBUser(ev.Username, sess)
	if err == nil {
		return Response{Value: 403, Description: "User already exists"}, nil
	}

	u := chatsess.NewUser(ev.Username, ev.Password)
	if err := u.Validate(ev.Password); err != nil {
		return Response{Value: 403, Description: err.Error()}, nil
	}
	if err := u.Put(sess); err != nil {
		return Response{Value: 500, Description: "Could not add to database: " + err.Error()}, nil
	}
	if err := chatsess.CreateUserFolder(ev.Username, sess); err != nil {
		return Response{Value: 202, Description: "User created successfully, but something went wrong with your bucket"}, nil
	}
	fl, err := chatsess.TextToSpeech(fmt.Sprintf("Welcome %s!\nThis is your bucket, all of your items will be stored here.", ev.Username), sess)
	if err != nil {
		return Response{Value: 202, Description: "User created successfully, but something went wrong with Polly"}, nil
	}
	obj := chatsess.Object{
		Key:  fmt.Sprintf("%s/welcome.mp3", ev.Username),
		Body: fl,
	}
	if err := obj.Put(sess); err != nil {
		return Response{Value: 202, Description: "User created successfully, but wasn't possible to generate Polly audio"}, nil
	}
	return Response{Value: 200, Description: "User created successfully"}, nil
}

func main() {
	lambda.Start(handler)
}
