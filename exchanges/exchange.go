package exchanges

import (
	"github.com/ysk229/go-rabbitmq/options"
)

// Exchange
type Exchange struct {
	*options.Exchange
}

// NewExchange
func NewExchange(e *options.Exchange) *Exchange {
	e.Durable = true
	return &Exchange{Exchange: e}
}
