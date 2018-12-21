package main

import (
	"context"

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
	Sessid      string
}

func handler(ctx context.Context, ev Event) (Response, error) {
	sess := session.Must(session.NewSession())
	u, err := chatsess.GetDBUserPass(ev.Username, ev.Password, sess)
	if err != nil {
		return Response{
			Value: 403, Description: err.Error(),
		}, nil
	}

	lg := chatsess.NewLogin(u.Username)
	if err := lg.Put(sess); err != nil {
		return Response{
			Value: 500, Description: err.Error(),
		}, nil
	}

	return Response{Value: 200, Description: "User logged in with success.", Sessid: lg.Sessid}, nil
}

func main() {
	lambda.Start(handler)
}
