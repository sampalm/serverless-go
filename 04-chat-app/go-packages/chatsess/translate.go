package chatsess

import (
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/translate"
)

const DATE_TL = "01-2006"

func TranslateText(src, tg, text string, sess *session.Session) (string, error) {
	ct := translate.New(sess)
	out, err := ct.Text(&translate.TextInput{
		SourceLanguageCode: aws.String(src),
		TargetLanguageCode: aws.String(tg),
		Text:               aws.String(text),
	})

	return *(out.TranslatedText), err
}

func CountCharacters(count int, sess *session.Session) error {
	cdb := dynamodb.New(sess)
	m := int(time.Now().Month())
	// DEBUG log.Println("Counting: " + strconv.Itoa(count))
	_, err := cdb.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String("ch_translation"),
		Key: map[string]*dynamodb.AttributeValue{
			"dateid":    {S: aws.String(time.Now().Format(DATE_TL))},
			"datemonth": {N: aws.String(strconv.Itoa(m))},
		},
		UpdateExpression:    aws.String("set datecount = datecount + :dcount"),
		ConditionExpression: aws.String("datecount < :max"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":dcount": {N: aws.String(strconv.Itoa(count))},
			":max":    {N: aws.String("1600000")},
		},
	})

	return err
}
