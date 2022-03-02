package options

import (
	"github.com/ysk229/go-rabbitmq/lib"
)

// Exchange Exchange
type Exchange struct {
	ExchangeName string
	// Exchange type the routing algorithm used depends on this
	// See https://www.rabbitmq.com/tutorials/amqp-concepts.html
	ExchangeType lib.ExchangeType
	// Durable sets the durable flag. Durable exchanges survives server restart.
	Durable bool
	// AutoDelete sets the auto-delete flag, this ensures the exchange is
	// deleted when it isn't bound to any more.
	AutoDelete bool
	Internal   bool
	NoWait     bool
	// Args a function that sets the binding exchange declare arguments
	// that are specific to the server's implementation of the exchange declare
	Args lib.Table
}

// QueueBind QueueBind
type QueueBind struct {
	Queue      string
	Exchange   string
	RoutingKey string
	NoWait     bool
	Args       lib.Table
}

// Queue Queue
type Queue struct {
	QueueName string
	// Durable queues remain active when a server restarts. Non-durable queues (transient queues) are purged if/when a server restarts.
	// Note that durable queues do not necessarily hold persistent messages, although it does not make sense to
	// send persistent messages to a transient queue.
	Durable bool
	// AutoDelete the queue is deleted when all consumers have finished using it.
	// The last consumer can be cancelled either explicitly or because its channel is closed.
	// If there was no consumer ever on the queue, it won't be deleted. Applications can
	// explicitly delete auto-delete queues using the Delete method as normal.
	AutoDelete bool
	// Exclusive queues may only be accessed by the current connection,
	// and are deleted when that connection closes. Passive declaration of an exclusive queue by other connections are not allowed.
	Exclusive bool
	NoWait    bool
	// Args a function that sets the binding queue declare arguments
	// that are specific to the server's implementation of the queue declare
	Args lib.Table
}
