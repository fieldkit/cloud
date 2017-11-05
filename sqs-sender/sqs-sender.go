package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func Print(msg *sqs.Message) error {
	fmt.Printf("%v", *msg)
	return nil
}

func main() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	svc := sqs.New(sess)

	resultUrl, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String("messages-test"),
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to find queue: %v", err)
		os.Exit(1)
	}

	result, err := svc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"Title": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String("The Whistler"),
			},
			"Author": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String("John Grisham"),
			},
			"WeeksOn": &sqs.MessageAttributeValue{
				DataType:    aws.String("Number"),
				StringValue: aws.String("6"),
			},
		},
		MessageBody: aws.String("Information about current NY Times fiction bestseller for week of 12/11/2016."),
		QueueUrl:    resultUrl.QueueUrl,
	})

	if err != nil {
		fmt.Println("Error", err)
		return
	}

	fmt.Println("Success", *result.MessageId)
}
