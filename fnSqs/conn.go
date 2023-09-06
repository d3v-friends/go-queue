package fnSqs

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/aws/smithy-go/ptr"
	"github.com/d3v-friends/pure-go/fnEnv"
)

type Conn struct {
	cfg    aws.Config
	client *sqs.Client
	url    *sqs.GetQueueUrlOutput
}

type FnReceiver func(msg types.Message) (err error)

const (
	EnvAccessKey = "AWS_ACCESS_KEY_ID"
	EnvSecretKEy = "AWS_SECRET_ACCESS_KEY"
)

// NewQueue aws 에서 환경변수 형식으로 데이터 입력하는것을 선호한다.
func NewQueue(queueNm string) (res *Conn, err error) {
	var ctx = context.Background()

	_ = fnEnv.Read(EnvAccessKey)
	_ = fnEnv.Read(EnvSecretKEy)

	res = &Conn{}

	if res.cfg, err = config.LoadDefaultConfig(ctx); err != nil {
		return
	}

	res.client = sqs.NewFromConfig(res.cfg)
	if res.url, err = res.client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: ptr.String(queueNm),
	}); err != nil {
		return
	}

	return
}

func (x *Conn) CreateQueue(ctx context.Context, queueNm string) (err error) {

}

func (x *Conn) Sender(ctx context.Context, msg string) (err error) {
	_, err = x.client.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody: &msg,
		QueueUrl:    x.url.QueueUrl,
	})
	return
}

func (x *Conn) Receiver(fn FnReceiver, chErr chan<- error) {
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
