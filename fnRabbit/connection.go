package fnRabbit

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Connection struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

type ExchangeKind string

func (x ExchangeKind) String() string {
	return string(x)
}

type ExchangeNm string

func (x ExchangeNm) String() string {
	return string(x)
}

const (
	ExchangeKindDefault ExchangeKind = "exchange"
	ExchangeKindDelayed ExchangeKind = "dlr-delayed-exchange"
)

const DelayedExchange = "dlr-delayed-exchange"
const DefaultExchange = "exchange"
