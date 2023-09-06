package fnSqs

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Queue struct {
	url    *sqs.GetQueueUrlOutput
	client *sqs.Client
}

func newQueue(url *sqs.GetQueueUrlOutput, client *sqs.Client) (res *Queue, err error) {
	res = &Queue{
		url:    url,
		client: client,
	}
	return
}

func (x *Queue) Send(ctx context.Context, msg string) (err error) {
	_, err = x.client.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody: &msg,
		QueueUrl:    x.url.QueueUrl,
	})
	return
}

func (x *Queue) Receiver(fn FnReceiver, chErr chan<- error) {
	var err error
	var ctx = context.Background()

	for {
		var message *sqs.ReceiveMessageOutput
		if message, err = x.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            x.url.QueueUrl,
			MaxNumberOfMessages: 1,
			WaitTimeSeconds:     30,
		}); err != nil {
			chErr <- err
		}

		for _, item := range message.Messages {
			if err = fn(item); err != nil {
				chErr <- err
			}

			if _, err = x.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
				QueueUrl:      x.url.QueueUrl,
				ReceiptHandle: item.ReceiptHandle,
			}); err != nil {
				chErr <- err
			}
		}
	}
}
