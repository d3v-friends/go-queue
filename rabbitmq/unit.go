package rabbitmq

import (
	"context"
	"encoding/json"
	"github.com/d3v-friends/go-tools/fnLogger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Unit[T any] struct {
	opt *UnitOptions
	ch  *amqp.Channel
}

type UnitOptions struct {
	ExchangeName string
	QueueName    string
}

func NewUnit[T any](mng *Manager, opt *UnitOptions) (res *Unit[T], err error) {
	res = &Unit[T]{
		opt: opt,
	}
	defer func() {
		if err != nil {
			res = nil
		}
	}()

	res.ch = mng.channel

	var has bool
	if has, err = mng.hasExchange(opt.ExchangeName); err != nil {
		return
	}

	if !has {
		if err = mng.createExchange(res.ch, opt.ExchangeName); err != nil {
			return
		}
	}

	if has, err = mng.hasQueue(opt.QueueName); err != nil {
		return
	}

	if !has {
		if err = mng.createQueue(res.ch, opt.ExchangeName, opt.QueueName); err != nil {
			return
		}
	}

	return
}

func (x *Unit[T]) Publish(v *T) (err error) {
	var body []byte
	if body, err = json.Marshal(v); err != nil {
		return
	}

	if err = x.ch.Publish(
		x.opt.ExchangeName,
		x.opt.QueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	); err != nil {
		return
	}

	return
}

func (x *Unit[T]) Consume(
	consumer Consumer[T],
	interceptors ...func(ctx context.Context, data amqp.Delivery) (context.Context, amqp.Delivery),
) {
	var delivery, err = x.ch.Consume(
		x.opt.QueueName,
		x.opt.QueueName,
		true,
		false,
		false,
		false,
		amqp.Table{},
	)

	if err != nil {
		panic(err)
	}

	for msg := range delivery {
		var ctx = context.TODO()
		ctx = fnLogger.SetID(ctx)
		ctx = fnLogger.SetLogger(ctx, fnLogger.NewLogger(fnLogger.LogLevelInfo))

		var logger = fnLogger.GetLogger(ctx)

		if len(interceptors) == 1 {
			ctx, msg = interceptors[0](ctx, msg)
		}

		var data = new(T)
		if err = json.Unmarshal(msg.Body, data); err != nil {
			logger.CtxError(ctx, err.Error())
			continue
		}

		if err = consumer.Consume(ctx, data); err != nil {
			logger.CtxError(ctx, err.Error())
			continue
		}

		ctx.Done()
	}
}
