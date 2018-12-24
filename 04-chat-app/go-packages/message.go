package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/sampalm/serverless/04-chat-app/go-packages/chatsess"
)

type Event struct {
	SessID string
	Text   string
}
type Response struct {
	Value       int
	Description string
}

func handler(c context.Context, ev Event) (Response, error) {
	sess := session.Must(session.NewSession())
	lg, err := chatsess.GetLogin(ev.SessID, sess)

	if err != nil {
		return Response{Value: 403, Description: err.Error()}, nil
	}

	ch := chatsess.NewChat(lg.Username, ev.Text)
	err = ch.Put(sess)
	if err != nil {
		return Response{Value: 500, Description: err.Error()}, nil
	}

	return Response{Value: 200, Description: "Sent: " + ev.Text}, nil
}

func main() {
	lambda.Start(handler)
}
