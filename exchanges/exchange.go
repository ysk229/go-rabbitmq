package exchanges

import (
	"github.com/ysk229/go-rabbitmq/v2/options"
)

// Exchange exchange
type Exchange struct {
	*options.Exchange
}

// NewExchange New Exchange
func NewExchange(e *options.Exchange) *Exchange {
	e.Durable = true
	return &Exchange{Exchange: e}
}
