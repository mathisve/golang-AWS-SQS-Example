package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var dynamo *dynamodb.DynamoDB
var table, region string

func init() {
	table = os.Getenv("TABLE")
	region = os.Getenv("REGION")

	dynamo = dynamodb.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	})))
}

func main() {
	lambda.Start(Handler)
}

// Handler for Lambda runtime
func Handler(ctx context.Context, sqsEvent events.SQSEvent) (err error) {
	for _, message := range sqsEvent.Records {
		err = putEvent(message)
		if err != nil {
			return err
		}
	}
	return err
}

func putEvent(message events.SQSMessage) error {
	_, err := dynamo.PutItem(&dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"Link": {
				S: message.MessageAttributes["Link"].StringValue,
			},
			"Key": {
				S: message.MessageAttributes["Key"].StringValue,
			},
			"S3Link": {
				S: message.MessageAttributes["S3Link"].StringValue,
			},
		},
		TableName: aws.String(table),
	})

	return err
}
