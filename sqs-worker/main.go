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

	worker, err := NewSqsWorker(svc, "messages-test", HandlerFunc(Print))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create SQS worker: %v", err)
		os.Exit(1)
	}

	worker.Start()
}
