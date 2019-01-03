package chatsess

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
)

// TextToSpeech converts a text into a mp3 audiofile. Returns a []byte and a error.
func TextToSpeech(text string, sess *session.Session) ([]byte, error) {
	if strings.Count(text, "") < 2 {
		return []byte{}, fmt.Errorf("Text is too small.")
	}
	if strings.Count(text, "") > 250 {
		return []byte{}, fmt.Errorf("Text is too large.")
	}

	cpl := polly.New(sess)
	res, err := cpl.SynthesizeSpeech(&polly.SynthesizeSpeechInput{
		LanguageCode: aws.String("en-US"),
		VoiceId:      aws.String("Salli"),
		Text:         aws.String(text),
		OutputFormat: aws.String("mp3"),
	})
	if err != nil {
		return []byte{}, err
	}
	body, err := ioutil.ReadAll(res.AudioStream)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}
