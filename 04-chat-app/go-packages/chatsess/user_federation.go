package chatsess

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

const (
	POLICY = "POLICY_CONTENT"
	ROLE   = "ROLE_URL"
)

type STSRole struct {
	RoleID string
	Arn    string
}

type STSToken struct {
	AccessKey       string
	SecretAccessKey string
	SessionToken    string
	Expiration      time.Time
	STSRole
}

func STSAssumeRole(username, userid string, sess *session.Session) (STSToken, error) {
	svc := sts.New(sess)

	res, err := svc.AssumeRole(&sts.AssumeRoleInput{
		DurationSeconds: aws.Int64(3600),
		ExternalId:      aws.String(userid),
		Policy:          aws.String(POLICY),
		RoleArn:         aws.String(ROLE),
		RoleSessionName: aws.String(username),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case sts.ErrCodeMalformedPolicyDocumentException:
				log.Println(sts.ErrCodeMalformedPolicyDocumentException, aerr.Error())
			case sts.ErrCodePackedPolicyTooLargeException:
				log.Println(sts.ErrCodePackedPolicyTooLargeException, aerr.Error())
			case sts.ErrCodeRegionDisabledException:
				log.Println(sts.ErrCodeRegionDisabledException, aerr.Error())
			default:
				log.Println(aerr.Error())
			}
		} else {
			log.Println(err.Error())
		}
		return STSToken{}, err
	}

	stk := STSToken{
		AccessKey:       *(res.Credentials.AccessKeyId),
		SecretAccessKey: *(res.Credentials.SecretAccessKey),
		SessionToken:    *(res.Credentials.SessionToken),
		Expiration:      *(res.Credentials.Expiration),
		STSRole: STSRole{
			RoleID: *(res.AssumedRoleUser.AssumedRoleId),
			Arn:    *(res.AssumedRoleUser.Arn),
		},
	}
	return stk, nil
}
