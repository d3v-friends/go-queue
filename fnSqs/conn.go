package fnSqs

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/d3v-friends/pure-go/fnEnv"
	"github.com/d3v-friends/pure-go/fnReflect"
)

type Conn struct {
	cfg    aws.Config
	client *sqs.Client
}

type FnReceiver func(msg types.Message) (err error)

const (
	EnvAccessKey = "AWS_ACCESS_KEY_ID"
	EnvSecretKey = "AWS_SECRET_ACCESS_KEY"
	EnvRegion    = "AWS_REGION"
)

// NewConn aws 에서 환경변수 형식으로 데이터 입력하는것을 선호한다.
func NewConn() (res *Conn, err error) {
	var ctx = context.Background()

	_ = fnEnv.Read(EnvAccessKey)
	_ = fnEnv.Read(EnvSecretKey)
	_ = fnEnv.Read(EnvRegion)

	res = &Conn{}

	if res.cfg, err = config.LoadDefaultConfig(ctx); err != nil {
		return
	}

	res.cfg.Region = fnEnv.Read(EnvRegion)
	res.client = sqs.NewFromConfig(res.cfg)

	return
}

type INewConn struct {
	Access string
	Secret string
	Region string
}

func NewConnWithAuth(i *INewConn) (res *Conn, err error) {
	var ctx = context.Background()
	res = &Conn{}
	if res.cfg, err = config.LoadDefaultConfig(ctx, config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
		i.Access,
		i.Secret,
		"",
	))); err != nil {
		return
	}

	res.cfg.Region = i.Region
	res.client = sqs.NewFromConfig(res.cfg)

	return
}

func (x *Conn) Create(ctx context.Context, queueNm string) (err error) {
	_, err = x.client.CreateQueue(ctx, &sqs.CreateQueueInput{
		QueueName: fnReflect.ToPointer(queueNm),
	})
	return
}

func (x *Conn) List(ctx context.Context) (_ []string, err error) {
	var res *sqs.ListQueuesOutput
	if res, err = x.client.ListQueues(ctx, &sqs.ListQueuesInput{}); err != nil {
		return
	}
	return res.QueueUrls, nil
}

func (x *Conn) NewQueue(ctx context.Context, queueNm string) (res *Queue, err error) {
	var url *sqs.GetQueueUrlOutput
	if url, err = x.client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: fnReflect.ToPointer(queueNm),
	}); err != nil {
		return
	}
	return newQueue(url, x.client)
}
