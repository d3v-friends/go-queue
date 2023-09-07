package fnSqs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/d3v-friends/go-pure/fnParams"
	"github.com/d3v-friends/go-pure/fnReflect"
	"reflect"
)

type Queue[BODY any] struct {
	url    *sqs.GetQueueUrlOutput
	client *sqs.Client
}

func NewQueue[BODY any](ctx context.Context, conn *Conn, queueNm string) (queue *Queue[BODY], err error) {
	var url *sqs.GetQueueUrlOutput
	if url, err = conn.QueueURL(ctx, queueNm); err != nil {
		return
	}

	var kind reflect.Kind
	if kind = reflect.TypeOf(*new(BODY)).Kind(); kind != reflect.Struct {
		err = fmt.Errorf("invalud BODY type: BODY.Kind()=%s", kind.String())
		return
	}

	queue = &Queue[BODY]{
		url:    url,
		client: conn.client,
	}

	return
}

func (x *Queue[BODY]) Send(ctx context.Context, msg *BODY, delaySec ...int32) (err error) {
	var delay = fnParams.Get(delaySec)
	var body []byte
	if body, err = json.Marshal(msg); err != nil {
		return
	}

	_, err = x.client.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody:  fnReflect.ToPointer(string(body)),
		QueueUrl:     x.url.QueueUrl,
		DelaySeconds: delay,
	})
	return
}

func (x *Queue[BODY]) Receiver(fn FnReceiver[BODY]) {
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
			if item.Body == nil {
				err = fmt.Errorf("empty body")
				return
			}

			var body = new(BODY)
			if err = json.Unmarshal([]byte(*item.Body), body); err != nil {
				return
			}

			if err = fn(body); err != nil {
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
