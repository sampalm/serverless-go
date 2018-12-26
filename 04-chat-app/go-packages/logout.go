package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/sampalm/serverless/04-chat-app/go-packages/chatsess"
)

type Event struct {
	Username string
	Sessid   string
}
type Response struct {
	Value       int
	Description string
}

func handler(ctx context.Context, ev Event) (Response, error) {
	sess := session.Must(session.NewSession())
	lg := chatsess.Login{Username: ev.Username, Sessid: ev.Sessid}
	if err := lg.Delete(sess); err != nil {
		return Response{
			Value: 500, Description: err.Error(),
		}, nil
	}

	return Response{Value: 200, Description: "User logged out with success."}, nil
}

func main() {
	lambda.Start(handler)
}
