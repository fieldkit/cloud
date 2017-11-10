package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
	"sync"
)

type SqsHandlerFunc func(msg *sqs.Message) error

func (f SqsHandlerFunc) HandleMessage(msg *sqs.Message) error {
	return f(msg)
}

type SqsHandler interface {
	HandleMessage(msg *sqs.Message) error
}

type SqsWorker struct {
	queueUrl *string
	svc      *sqs.SQS
	handler  SqsHandler
}

func NewSqsWorker(svc *sqs.SQS, name string, handler SqsHandler) (w *SqsWorker, err error) {
	resultUrl, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(name),
	})
	if err != nil {
		return nil, err
	}

	w = &SqsWorker{
		queueUrl: resultUrl.QueueUrl,
		svc:      svc,
		handler:  handler,
	}
	return
}

func (w *SqsWorker) Start() {
	log.Println("worker: Start polling...")

	for {
		messages, err := w.svc.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl: w.queueUrl,
			AttributeNames: aws.StringSlice([]string{
				"SentTimestamp",
			}),
			MaxNumberOfMessages: aws.Int64(5),
			MessageAttributeNames: aws.StringSlice([]string{
				"All",
			}),
			WaitTimeSeconds: aws.Int64(20),
		})
		if err != nil {
			log.Println(err)
			continue
		}
		if len(messages.Messages) > 0 {
			w.run(messages.Messages)
		}
	}
}

func (w *SqsWorker) run(messages []*sqs.Message) {
	numberOfMessages := len(messages)
	log.Printf("worker: Received %d messages", numberOfMessages)

	var wg sync.WaitGroup
	wg.Add(numberOfMessages)
	for i := range messages {
		go func(m *sqs.Message) {
			defer wg.Done()
			if err := w.handleMessage(m); err != nil {
				log.Println(err)
			}
		}(messages[i])
	}

	wg.Wait()
}

func (w *SqsWorker) handleMessage(m *sqs.Message) error {
	err := w.handler.HandleMessage(m)
	if err != nil {
		return err
	}
	_, err = w.svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      w.queueUrl,
		ReceiptHandle: m.ReceiptHandle,
	})
	if err != nil {
		return err
	}
	return nil
}
