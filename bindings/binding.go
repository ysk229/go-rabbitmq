package bindings

import "github.com/ysk229/go-rabbitmq/v2/options"

// Binding binding struct
type Binding struct {
	*options.QueueBind
}

// NewBinding  new binding
func NewBinding(q *options.QueueBind) *Binding {
	return &Binding{QueueBind: q}
}
