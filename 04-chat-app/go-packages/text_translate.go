package main

import (
	"context"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/sampalm/serverless/04-chat-app/go-packages/chatsess"
)

type Event struct {
	SessID string
	Source string
	Target string
	Text   string
}

type Response struct {
	Value       int
	Description string
	Body        string
}

func handler(ctx context.Context, ev Event) (Response, error) {
	// DEBUG log.Printf("Source: %s, Target: %s, Text: %s", ev.Source, ev.Target, ev.Text)
	if ev.Source == "" {
		return Response{
			Value:       403,
			Description: "Source language must be set.",
		}, nil
	}
	if ev.Target == "" {
		return Response{
			Value:       403,
			Description: "Target language must be set.",
		}, nil
	}
	if ev.Text == "" || len(ev.Text) > 150 {
		return Response{
			Value:       403,
			Description: "Text can not be empty or larger than 150 characters.",
		}, nil
	}
	sess := session.Must(session.NewSession())

	_, err := chatsess.GetLogin(ev.SessID, sess)
	if err != nil {
		return Response{Value: 403, Description: err.Error()}, nil
	}

	if err := chatsess.CountCharacters(len(ev.Text), sess); err != nil {
		if strings.Contains(err.Error(), "ConditionalCheckFailedException") {
			return Response{Value: 500, Description: "Translation limits exceeded. Wait a while before a new request."}, nil
		}
		return Response{Value: 500, Description: "CountCharacters: " + err.Error()}, nil
	}

	res, err := chatsess.TranslateText(ev.Source, ev.Target, ev.Text, sess)
	if err != nil {
		return Response{
			Value:       500,
			Description: err.Error(),
		}, nil
	}
	return Response{
		Value:       200,
		Description: "Your text was translated successfully.",
		Body:        res,
	}, nil
}

func main() {
	lambda.Start(handler)
}
