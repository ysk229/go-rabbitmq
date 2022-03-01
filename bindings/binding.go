package bindings

import "github.com/ysk229/go-rabbitmq/options"

type Binding struct {
	*options.QueueBind
}

func NewBinding(q *options.QueueBind) *Binding {
	return &Binding{QueueBind: q}
}
