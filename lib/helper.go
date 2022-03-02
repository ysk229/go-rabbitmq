package lib

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

// Table interface
type Table map[string]interface{}

// WrapTable eliminates the need for the amqp lib directly.
func WrapTable(table Table) amqp.Table {
	amqpTable := amqp.Table{}
	for k, v := range table {
		amqpTable[k] = v
	}

	return amqpTable
}
