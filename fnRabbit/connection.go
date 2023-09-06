package fnRabbit

import (
	"fmt"
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

func NewConnection(i *IConnection) (res *Connection, err error) {
	res = &Connection{}
	if res.conn, err = amqp.Dial(i.Dial()); err != nil {
		return
	}

	defer func() {
		if err != nil {
			_ = res.conn.Close()
		}
	}()

	if res.channel, err = res.conn.Channel(); err != nil {
		return
	}

	return
}

type IExchangeDeclare struct {
	Name       string
	Kind       ExchangeKind
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table
}

func (x *Connection) DeclareExchange(
	i *IExchangeDeclare,
) (err error) {
	err = x.channel.ExchangeDeclare(
		i.Name,
		i.Kind.String(),
		i.Durable,
		i.AutoDelete,
		i.Internal,
		i.NoWait,
		i.Args,
	)
	return
}

type IConnection struct {
	Host     string
	Username string
	Password string
}

func (x *IConnection) Dial() string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s",
		x.Username,
		x.Password,
		x.Host,
	)
}
