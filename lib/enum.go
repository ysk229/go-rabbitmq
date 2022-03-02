package lib

import amqp "github.com/rabbitmq/amqp091-go"

const (
	// Transient Transient
	Transient = amqp.Transient
	// Persistent Persistent
	Persistent = amqp.Persistent
)

// AcknowledgementMode AcknowledgementMode
type AcknowledgementMode string

const (
	// Ack AcknowledgementMode
	Ack AcknowledgementMode = "ACK"
	// Nack AcknowledgementMode
	Nack AcknowledgementMode = "NACK"
)

// IsNack AcknowledgementMode
func (a AcknowledgementMode) IsNack() bool {
	return a == Nack
}

// IsAck AcknowledgementMode
func (a AcknowledgementMode) IsAck() bool {
	return a == Ack
}

// ExchangeType  ExchangeType
type ExchangeType string

const (
	// Direct exchange delivers messages to queues based on the message routing key.
	// A direct exchange is ideal for the unicast routing of messages (although they can be used for multicast routing as well).
	// Here is how it works
	Direct ExchangeType = "direct"
	// Topic exchanges route messages to one or many queues based on matching between a message routing key and
	// the pattern that was used to bind a queue to an exchange. The topic exchange type is often used to implement
	// various publish/subscribe pattern variations. Topic exchanges are commonly used for the multicast routing of messages.
	Topic ExchangeType = "topic"
	// Headers exchange is designed for routing on multiple attributes that are more easily expressed as message headers than a routing key.
	// Header exchanges ignore the routing key attribute. Instead, the attributes used for routing are taken from the headers attribute.
	// A message is considered matching if the value of the header equals the value specified upon binding.
	Headers ExchangeType = "headers"
	// Fanout Faut exchange routes messages to all of the queues that are bound to it and the routing key is ignored.
	// If N queues are bound to a fanout exchange, when a new message is published to that exchange a copy of the message is delivered to
	// all N queues. Fanout exchanges are ideal for the broadcast routing of messages.
	Fanout ExchangeType = "fanout"
)

//String
func (e ExchangeType) String() string {
	return string(e)
}
