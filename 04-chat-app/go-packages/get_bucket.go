package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/sampalm/serverless/04-chat-app/go-packages/chatsess"
)

type Event struct {
	SessID string
}

type Response struct {
	Value       int
	Description string
	Body        []chatsess.Object
}

func handler(ev Event) (Response, error) {
	sess := session.Must(session.NewSession())
	lg, err := chatsess.GetLogin(ev.SessID, sess)
	if err != nil {
		return Response{Value: 403, Description: err.Error()}, nil
	}

	objs, err := chatsess.ListObjects(lg.Username, sess)
	if err != nil {
		return Response{
			Value:       500,
			Description: err.Error(),
		}, nil
	}

	return Response{
		Value:       200,
		Description: "Objects listed successfully.",
		Body:        objs,
	}, nil
}

func main() {
	lambda.Start(handler)
}
