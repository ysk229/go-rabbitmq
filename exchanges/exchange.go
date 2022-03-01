package exchanges

import (
	"github.com/ysk229/go-rabbitmq/options"
)

type Exchange struct {
	*options.Exchange
}

func NewExchange(e *options.Exchange) *Exchange {
	e.Durable = true
	return &Exchange{Exchange: e}
}
