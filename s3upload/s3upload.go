package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var s3session *s3.S3
var sqsQueue *sqs.SQS
var region, bucket, queueURL string

// InputEvent is the input
type InputEvent struct {
	Link string `json:"link"`
	Key  string `json:"key"`
}

// OutputEvent is the output
type OutputEvent struct {
	Link   string `json:"link"`
	Key    string `json:"key"`
	S3Link string `json:"s3link"`
}

func init() {
	region = os.Getenv("REGION")
	bucket = os.Getenv("BUCKET")
	queueURL = os.Getenv("QUEUE")

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))

	s3session = s3.New(sess)
	sqsQueue = sqs.New(sess)
}

func main() {
	lambda.Start(Handler)
}

// Handler for Lambda runtime
func Handler(ctx context.Context, event InputEvent) (err error) {
	image, err := getImage(event.Link)
	if err != nil {
		return err
	}

	_, err = s3session.PutObject(&s3.PutObjectInput{
		Body:   bytes.NewReader(image),
		Bucket: aws.String(bucket),
		Key:    aws.String(event.Key),
	})

	if err != nil {
		return err
	}

	var output = OutputEvent{
		Key:  event.Key,
		Link: event.Link,
	}

	output.S3Link = fmt.Sprintf("https://%s.%s.amazonaws.com/%s", bucket, region, output.Key)

	_, err = sqsQueue.SendMessage(&sqs.SendMessageInput{
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"Link": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(output.Link),
			},
			"Key": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(output.Key),
			},
			"S3Link": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(output.S3Link),
			},
		},
		MessageBody: aws.String("New event!"),
		QueueUrl:    aws.String(queueURL),
	})

	return err
}

func getImage(url string) (bytes []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return bytes, err
	}

	defer resp.Body.Close()

	bytes, err = ioutil.ReadAll(resp.Body)
	return bytes, err
}
