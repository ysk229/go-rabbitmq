package queues

import (
	"github.com/ysk229/go-rabbitmq/options"
)

// Queue
type Queue struct {
	*options.Queue
}

// NewQueue
func NewQueue(q *options.Queue) *Queue {
	if q.Args == nil {
		q.Args = map[string]interface{}{"x-ha-policy": "all"}
	}
	q.Durable = true
	return &Queue{Queue: q}
}
