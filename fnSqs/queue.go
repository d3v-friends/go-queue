package fnSqs

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/d3v-friends/pure-go/fnParams"
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

func (x *Queue) Send(ctx context.Context, msg string, delaySec ...int32) (err error) {
	var delay = fnParams.Get(delaySec)
	_, err = x.client.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody:  &msg,
		QueueUrl:     x.url.QueueUrl,
		DelaySeconds: delay,
	})
	return
}

func (x *Queue) ReceiverWithChannel(fn FnReceiver, chErr chan<- error) {
	var err error
	var ctx = context.Background()

	for {
		var message *sqs.ReceiveMessageOutput
		if message, err = x.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl: x.url.QueueUrl,
			//MaxNumberOfMessages: 1,
			WaitTimeSeconds: 1,
		}); err != nil {
			chErr <- err
			continue
		}

		for _, item := range message.Messages {
			if err = fn(item); err != nil {
				chErr <- err
				continue
			}

			if _, err = x.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
				QueueUrl:      x.url.QueueUrl,
				ReceiptHandle: item.ReceiptHandle,
			}); err != nil {
				chErr <- err
				continue
			}
		}
	}
}

func (x *Queue) Receiver(fn FnReceiver) {
	var err error
	var ctx = context.Background()

	for {
		var message *sqs.ReceiveMessageOutput
		if message, err = x.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl: x.url.QueueUrl,
			//MaxNumberOfMessages: 1,
			WaitTimeSeconds: 1,
		}); err != nil {
			panic(err)
		}

		for _, item := range message.Messages {
			if err = fn(item); err != nil {
				panic(err)
			}

			if _, err = x.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
				QueueUrl:      x.url.QueueUrl,
				ReceiptHandle: item.ReceiptHandle,
			}); err != nil {
				panic(err)
			}
		}
	}
}
