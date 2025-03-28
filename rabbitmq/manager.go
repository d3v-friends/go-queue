package rabbitmq

import (
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"net/http"
	"time"
)

type Options struct {
	Host        string
	QueuePort   int64
	ManagerPort int64
	Username    string
	Password    string
}

func (x *Options) host() string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%d",
		x.Username,
		x.Password,
		x.Host,
		x.QueuePort,
	)
}

type Manager struct {
	opts       *Options
	connection *amqp.Connection
	channel    *amqp.Channel
}

func NewManager(opt *Options) (mng *Manager, err error) {
	mng = &Manager{opts: opt}
	if mng.connection, err = amqp.Dial(opt.host()); err != nil {
		return
	}
	if mng.channel, err = mng.connection.Channel(); err != nil {
		return
	}
	return
}

func (x *Manager) hasQueue(queueName string) (res bool, err error) {
	var list ApiQueues
	if list, err = x.getQueueList(); err != nil {
		return
	}

	res = list.HasByName(queueName)
	return
}

func (x *Manager) getQueueList() (res ApiQueues, err error) {
	var request *http.Request
	if request, err = http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(
			"http://%s:%d/api/queues",
			x.opts.Host,
			x.opts.ManagerPort,
		),
		nil,
	); err != nil {
		return
	}
	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth(x.opts.Username, x.opts.Password)

	var client = &http.Client{
		Timeout: time.Second * 5,
	}

	var resp *http.Response
	if resp, err = client.Do(request); err != nil {
		return
	}

	res = make(ApiQueues, 0)
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return
	}

	return
}

func (x *Manager) hasExchange(exchangeName string) (res bool, err error) {
	var list ApiExchanges
	if list, err = x.getExchangeList(); err != nil {
		return
	}

	res = list.HasByName(exchangeName)
	return
}

func (x *Manager) getExchangeList() (res ApiExchanges, err error) {
	var request *http.Request
	if request, err = http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(
			"http://%s:%d/api/exchanges",
			x.opts.Host,
			x.opts.ManagerPort,
		),
		nil,
	); err != nil {
		return
	}
	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth(x.opts.Username, x.opts.Password)

	var client = &http.Client{
		Timeout: time.Second * 5,
	}

	var resp *http.Response
	if resp, err = client.Do(request); err != nil {
		return
	}

	res = make(ApiExchanges, 0)
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return
	}

	return
}

func (x *Manager) createExchange(
	channel *amqp.Channel,
	exchangeName string,
) (err error) {
	return channel.ExchangeDeclare(
		exchangeName,
		"direct",
		true,
		false,
		false,
		false,
		amqp.Table{},
	)
}

func (x *Manager) createQueue(
	channel *amqp.Channel,
	exchangeName string,
	queueName string,
) (err error) {
	if _, err = channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		amqp.Table{},
	); err != nil {
		return
	}

	if err = channel.QueueBind(
		queueName,
		queueName,
		exchangeName,
		false,
		amqp.Table{},
	); err != nil {
		return
	}

	return
}
