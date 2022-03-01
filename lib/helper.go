package lib

import (
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Table map[string]interface{}

// WrapTable eliminates the need for the amqp lib directly.
func WrapTable(table Table) amqp.Table {
	amqpTable := amqp.Table{}
	for k, v := range table {
		amqpTable[k] = v
	}

	return amqpTable
}

func GenerateReplyToQueueName() string {
	return "reply-to-" + RandomStr()
}

func RandomStr() string {
	return uuid.New().String()
}