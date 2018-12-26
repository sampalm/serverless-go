package chatsess

import (
	"html"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Chat struct {
	DateID   string
	Time     time.Time
	Username string
	Text     string
}

func NewChat(Username, Text string) Chat {
	return Chat{
		DateID:   time.Now().Format(DATE_FMT),
		Time:     time.Now(),
		Username: Username,
		Text:     html.EscapeString(Text),
	}
}

func ItemToChat(item map[string]*dynamodb.AttributeValue) Chat {
	dti, ok := item["dateid"]
	if !ok {
		return Chat{}
	}
	tmi, ok := item["timestamp"]
	if !ok {
		return Chat{}
	}
	uni := item["username"]
	txi := item["text"]

	return Chat{
		DateID:   *(dti.S),
		Time:     DBtoTime(tmi.N),
		Username: *(uni.S),
		Text:     *(txi.S),
	}
}

func (c Chat) Put(sess *session.Session) error {
	cdb := dynamodb.New(sess)

	_, err := cdb.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("ch_chats"),
		Item: map[string]*dynamodb.AttributeValue{
			"dateid":    {S: aws.String(c.DateID)},
			"timestamp": {N: TimetoDB(c.Time)},
			"username":  {S: aws.String(c.Username)},
			"text":      {S: aws.String(c.Text)},
		},
	})

	return err
}

func GetChat(sess *session.Session) ([]Chat, error) {
	cdb := dynamodb.New(sess)

	res, err := cdb.Query(&dynamodb.QueryInput{
		TableName:              aws.String("ch_chats"),
		KeyConditionExpression: aws.String("dateid = :d"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":d": {S: aws.String(time.Now().Format(DATE_FMT))},
		},
	})
	if err != nil {
		return []Chat{}, err
	}

	chs := []Chat{}
	for _, v := range res.Items {
		chs = append(chs, ItemToChat(v))
	}
	return chs, nil
}

func GetChatAfter(DateID string, t time.Time, sess *session.Session) ([]Chat, error) {
	cdb := dynamodb.New(sess)

	res, err := cdb.Query(&dynamodb.QueryInput{
		TableName:              aws.String("ch_chats"),
		KeyConditionExpression: aws.String("dateid = :d"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":d": {S: aws.String(DateID)},
		},
		ExclusiveStartKey: map[string]*dynamodb.AttributeValue{
			"dateid":    {S: aws.String(DateID)},
			"timestamp": {N: TimetoDB(t)},
		},
	})
	if err != nil {
		return []Chat{}, err
	}

	chs := []Chat{}
	for _, v := range res.Items {
		chs = append(chs, ItemToChat(v))
	}

	return chs, nil
}
