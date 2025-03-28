package rabbitmq

import (
	"context"
	"github.com/shopspring/decimal"
)

type ApiQueue struct {
	Arguments                  map[string]any  `json:"arguments"`
	AutoDelete                 bool            `json:"auto_delete"`
	ConsumerCapacity           decimal.Decimal `json:"consumer_capacity"`
	ConsumerUtilisation        decimal.Decimal `json:"consumer_utilisation"`
	Consumers                  decimal.Decimal `json:"consumers"`
	Durable                    bool            `json:"durable"`
	EffectivePolicyDefinition  map[string]any  `json:"effective_policy_definition"`
	Exclusive                  bool            `json:"exclusive"`
	Memory                     decimal.Decimal `json:"memory"`
	MessageBytes               decimal.Decimal `json:"message_bytes"`
	MessageBytesPagedOut       decimal.Decimal `json:"message_bytes_paged_out"`
	MessageBytesRam            decimal.Decimal `json:"message_bytes_ram"`
	MessageBytesReady          decimal.Decimal `json:"message_bytes_ready"`
	MessageBytesUnacknowledged decimal.Decimal `json:"message_bytes_unacknowledged"`
	Messages                   decimal.Decimal `json:"messages"`
	MessagesDetails            struct {
		Rate decimal.Decimal
	} `json:"messages_details"`
	MessagesPagedOut    decimal.Decimal `json:"messages_paged_out"`
	MessagesPersistent  decimal.Decimal `json:"messages_persistent"`
	MessagesRam         decimal.Decimal `json:"messages_ram"`
	MessagesReady       decimal.Decimal `json:"messages_ready"`
	MessageReadyDetails struct {
		Rate decimal.Decimal
	}
	MessageReadyRam              decimal.Decimal `json:"message_ready_ram"`
	MessageUnacknowledged        decimal.Decimal `json:"message_unacknowledged"`
	MessageUnacknowledgedDetails struct {
		Rate decimal.Decimal
	} `json:"message_unacknowledged_details"`
	MessageUnacknowledgedRam decimal.Decimal `json:"message_unacknowledged_ram"`
	Name                     string          `json:"name"`
	Node                     string          `json:"node"`
	Reductions               decimal.Decimal `json:"reductions"`
	ReductionsDetails        struct {
		Rate decimal.Decimal
	} `json:"reductions_details"`
	State          string          `json:"state"`
	StorageVersion decimal.Decimal `json:"storage_version"`
	Type           string          `json:"type"`
	Vhost          string          `json:"vhost"`
}

type ApiQueues []*ApiQueue

func (x ApiQueues) HasByName(name string) bool {
	for _, queue := range x {
		if queue.Name == name {
			return true
		}
	}
	return false
}

type ApiExchange struct {
	Arguments              map[string]string `json:"arguments"`
	AutoDelete             bool              `json:"auto_delete"`
	Durable                bool              `json:"durable"`
	Internal               bool              `json:"internal"`
	Name                   string            `json:"name"`
	Type                   string            `json:"type"`
	UserWhoPerformedAction string            `json:"user_who_performed_action"`
	Vhost                  string            `json:"vhost"`
}

type ApiExchanges []*ApiExchange

func (x ApiExchanges) HasByName(name string) bool {
	for _, exchange := range x {
		if exchange.Name == name {
			return true
		}
	}
	return false
}

type Consumer[T any] interface {
	Consume(ctx context.Context, v *T) (err error)
}
