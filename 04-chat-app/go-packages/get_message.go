package main

import (
	"context"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/sampalm/serverless/04-chat-app/go-packages/chatsess"
)

type Event struct {
	SessID   string
	LastID   string
	LastTime string
}

type Response struct {
	Value       int
	Description string
	Chats       []chatsess.Chat
}

func handler(ctx context.Context, ev Event) (Response, error) {
	sess := session.Must(session.NewSession())

	_, err := chatsess.GetLogin(ev.SessID, sess)
	if err != nil {
		return Response{
			Value:       403,
			Description: err.Error(),
		}, nil
	}

	if ev.LastID != "" {
		ltime, err := time.Parse(time.RFC3339, ev.LastTime)
		if err != nil {
			return Response{
				Value:       500,
				Description: err.Error(),
			}, nil
		}

		ch, err := chatsess.GetChatAfter(ev.LastID, ltime, sess)
		if err != nil {
			return Response{
				Value:       500,
				Description: err.Error(),
			}, nil
		}

		return Response{Value: 200, Description: "Messages from " + ltime.Format(chatsess.DATE_FMT) + " successfully returned", Chats: ch}, nil
	}

	ch, err := chatsess.GetChat(sess)
	if err != nil {
		return Response{
			Value:       500,
			Description: err.Error(),
		}, nil
	}

	return Response{Value: 200, Description: "Messages successfully returned", Chats: ch}, nil
}

func main() {
	lambda.Start(handler)
}
