package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/sampalm/serverless/04-chat-app/go-packages/chatsess"
)

type Event struct {
	SessID   string
	Filename string
}

type Response struct {
	Value       int
	Description string
	ContentType string
	Body        string
}

func handler(ev Event) (Response, error) {
	sess := session.Must(session.NewSession())
	lg, err := chatsess.GetLogin(ev.SessID, sess)
	if err != nil {
		return Response{Value: 403, Description: err.Error()}, nil
	}

	ct, enc, err := chatsess.DownloadObject(lg.Username, ev.Filename, sess)
	if err != nil {
		return Response{
			Value:       500,
			Description: err.Error(),
		}, nil
	}

	return Response{
		Value:       200,
		Description: "Retrieved file successfully.",
		ContentType: ct,
		Body:        enc,
	}, nil
}

func main() {
	lambda.Start(handler)
}
