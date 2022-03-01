package consumers

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ysk229/go-rabbitmq/bindings"
	"github.com/ysk229/go-rabbitmq/channels"
	"github.com/ysk229/go-rabbitmq/exchanges"
	"github.com/ysk229/go-rabbitmq/lib"
	"github.com/ysk229/go-rabbitmq/options"
	"github.com/ysk229/go-rabbitmq/queues"
	"log"
	"time"
)

type ConsumerOpt struct {
	QueueName    string
	RoutingKey   string
	Exchange     string
	ExchangeType lib.ExchangeType
	// AutoAck sets the auto-ack flag. When this is set to false, you must
	// manually ack any deliveries. This is always true for the Client when
	// consuming replies
	AutoAck bool
	// Exclusive sets the exclusive flag. When this is set to true, no other
	// instances can consume from a given queue. This has no affect on the
	// Client when consuming replies where it's always set to true so that no
	// two clients can consume from the same reply-to queue.
	Exclusive bool
	// NoWait it does not wait for the server to
	// confirm the request and immediately begin deliveries.
	NoLocal bool
	NoWait  bool
	Args    lib.Table

	ResendDelay time.Duration //消息发送失败后，多久秒重发,默认是3s
	ResendNum   int           //消息重发次数

}
type CallBack struct {
	Fnc func(Delivery)
}
type Delivery struct {
	amqp.Delivery
}
type ConsumerOption func(*Consumer)

type Consumer struct {
	*channels.Channel
	opt *ConsumerOpt
	cb  *CallBack
}

func NewConsumer(ch *channels.Channel) *Consumer {
	c := &Consumer{Channel: ch}
	return c
}

func WithOptionsConsumer(opt *ConsumerOpt) ConsumerOption {
	return func(c *Consumer) {
		c.opt = opt
		if c.opt.ResendDelay == 0 {
			c.opt.ResendDelay = 3
		}
		c.init()
	}
}

func WithOptionsConsumerCallBack(cb *CallBack) ConsumerOption {
	return func(c *Consumer) {
		c.cb = cb
	}
}

func (c *Consumer) Consumer(ch *channels.Channel, opts ...ConsumerOption) {
	for _, opt := range opts {
		opt(c)
	}
	c.subscribe(ch)
}

//err := msg.Ack(false)
//err := msg.Nack(false, true)
func (c *Consumer) subscribe(ch *channels.Channel) {
	deliveries, err := ch.GetChannel().Consume(
		c.opt.QueueName,
		"",
		c.opt.AutoAck,
		c.opt.Exclusive,
		c.opt.NoLocal,
		c.opt.NoWait,
		lib.WrapTable(c.opt.Args),
	)
	if err != nil {
		log.Printf("cannot consume from: %q, %v", c.opt.QueueName, err)
		return
	}
	//
	log.Printf("subscribed...")
	for d := range deliveries {
		log.Printf(
			"got %dB delivery: [%v] %q",
			len(d.Body),
			d.DeliveryTag,
			d.Body,
		)
		if c.cb != nil {
			c.cb.Fnc(Delivery{d})
		}
	}
	log.Printf("handle: deliveries channel closed")
}

func (c *Consumer) init() {
	//	create queue queueBind exchange
	exchangeName, err := c.Channel.Exchange(exchanges.NewExchange(&options.Exchange{ExchangeName: c.opt.Exchange, ExchangeType: c.opt.ExchangeType}))
	if err != nil {
		log.Println(err)
	}
	//log.Println(exchangeName)
	////queues
	q, err := c.Channel.Queue(queues.NewQueue(&options.Queue{QueueName: c.opt.QueueName}))
	if err != nil {
		log.Println(err)
	}
	//log.Println(q.Name)
	err = c.Channel.BindingQueue(bindings.NewBinding(&options.QueueBind{Queue: q.Name, Exchange: exchangeName, RoutingKey: c.opt.RoutingKey}))
	if err != nil {
		log.Println(err)
	}
}
