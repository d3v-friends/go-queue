package fnRabbit

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/d3v-friends/pure-go/fnLogger"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

type Queue[BODY any] struct {
	channel  *amqp.Channel
	queueNm  string
	exchange ExchangeKind
}

func NewQueue[BODY any](
	conn *Connection,
	i *IDeclareQueue,
	exchange ExchangeKind,
) (res *Queue[BODY], err error) {
	res = &Queue[BODY]{
		channel:  conn.channel,
		queueNm:  i.Name,
		exchange: exchange,
	}

	if _, err = res.channel.QueueDeclare(
		i.Name,
		i.Durable,
		i.AutoDelete,
		i.Exclusive,
		i.NoWait,
		i.Args,
	); err != nil {
		return
	}

	if err = res.channel.QueueBind(
		i.Name,
		i.Name,
		exchange.String(),
		false,
		nil,
	); err != nil {
		return
	}

	return
}

func (x *Queue[BODY]) Publish(ctx context.Context, v *BODY) (err error) {
	var body []byte
	if body, err = json.Marshal(v); err != nil {
		return
	}

	var headers = make(amqp.Table)
	err = x.channel.PublishWithContext(
		ctx,
		x.exchange.String(),
		x.queueNm,
		false,
		false,
		amqp.Publishing{
			Headers:         headers,
			ContentType:     "application/json",
			ContentEncoding: "utf8",
			DeliveryMode:    0,
			Priority:        0,
			CorrelationId:   "",
			ReplyTo:         "",
			Expiration:      "",
			MessageId:       "",
			Timestamp:       time.Now(),
			Type:            "",
			UserId:          "",
			AppId:           "",
			Body:            body,
		},
	)

	return
}

func (x *Queue[BODY]) Consume(
	logger fnLogger.IfLogger,
	fn FnConsumer[BODY],
) {
	logger = logger.WithFields(fnLogger.Fields{
		"queue": x.queueNm,
	})

	var err error
	var args = make(amqp.Table)
	var msg <-chan amqp.Delivery
	if msg, err = x.channel.Consume(
		x.queueNm,
		fmt.Sprintf("%s-consumer", x.queueNm),
		true,
		false,
		false,
		false,
		args,
	); err != nil {
		logger.Error("fail consume: queue=%s, err=%s", x.queueNm, err.Error())
		return
	}

	for item := range msg {
		var body *BODY
		if err = json.Unmarshal(item.Body, body); err != nil {
			panic("invalid queue message")
		}

		if err = fn(body); err != nil {
			panic("fail consume")
		}
	}

}

type IDeclareQueue struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

type FnConsumer[T any] func(body *T) (err error)
