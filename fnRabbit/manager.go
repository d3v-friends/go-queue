package fnRabbit

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Manager struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

type IConfig struct {
	Host     string
	Username string
	Password string
}

type IDeclareQueue struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

func (x *IConfig) Dial() string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s",
		x.Username,
		x.Password,
		x.Host,
	)
}

func NewManager(i *IConfig) (mng *Manager, err error) {
	mng = &Manager{}
	if mng.conn, err = amqp.Dial(i.Dial()); err != nil {
		return
	}

	defer func() {
		if err != nil {
			_ = mng.conn.Close()
		}
	}()

	if mng.channel, err = mng.conn.Channel(); err != nil {
		return
	}

	return
}

func (x *Manager) CreateQueue(i *IDeclareQueue) (err error) {
	_, err = x.channel.QueueDeclare(
		i.Name,
		i.Durable,
		i.AutoDelete,
		i.Exclusive,
		i.NoWait,
		i.Args,
	)
	return
}

// CreateExchange 모두 Delay queue exchange 에 입력하되, delay 값을 조절
func (x *Manager) CreateExchange(i *IDeclareExchange) (err error) {
	err = x.channel.ExchangeDeclare(
		i.Name,
		ExchangeKindDelayed.String(),
		i.Durable,
		i.AutoDelete,
		i.Internal,
		i.NoWait,
		i.Args,
	)
	return
}

type IDeclareExchange struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table
}
