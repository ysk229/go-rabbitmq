package consumers

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ysk229/go-rabbitmq/bindings"
	"github.com/ysk229/go-rabbitmq/channels"
	"github.com/ysk229/go-rabbitmq/exchanges"
	"github.com/ysk229/go-rabbitmq/lib"
	"github.com/ysk229/go-rabbitmq/options"
	"github.com/ysk229/go-rabbitmq/queues"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// ConsumerOpt consumer options
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

	ReReceiveDelay time.Duration //After the message consumption failure, how long seconds to resend,default is 3s
	ReReceiveNum   int           //Number of message reconsumption,default is 3 num ，-1 to disable the retry mechanism
	GracefulDelay  time.Duration //default is 10s

}

// CallBack call back consumer
type CallBack struct {
	Fnc func(Delivery) error
}

// Delivery consumer result data
type Delivery struct {
	amqp.Delivery
}

// ConsumerOption Consumer Option
type ConsumerOption func(*Consumer)

// Consumer consumer
type Consumer struct {
	*channels.Channel
	opt   *ConsumerOpt
	cb    *CallBack
	wg    *sync.WaitGroup
	close bool
}

// NewConsumer New Consumer
func NewConsumer(ch *channels.Channel) *Consumer {
	c := &Consumer{Channel: ch, wg: &sync.WaitGroup{}, opt: &ConsumerOpt{}}
	return c
}

// WithOptionsConsumer With Options Consumer
func WithOptionsConsumer(opt *ConsumerOpt) ConsumerOption {
	return func(c *Consumer) {
		c.opt = opt
		if c.opt.ReReceiveDelay == 0 {
			c.opt.ReReceiveDelay = 3
		}
		if c.opt.ReReceiveNum == 0 {
			c.opt.ReReceiveNum = 3
		}
		if c.opt.GracefulDelay == 0 {
			c.opt.GracefulDelay = 10
		}
		c.init()
	}
}

// WithOptionsConsumerCallBack With Options Consumer CallBack
func WithOptionsConsumerCallBack(cb *CallBack) ConsumerOption {
	return func(c *Consumer) {
		c.cb = cb
	}
}

// Consumer consumer
func (c *Consumer) Consumer(ch *channels.Channel, opts ...ConsumerOption) {
	for _, opt := range opts {
		opt(c)
	}
	c.wg.Add(1)
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
	defer c.wg.Done()
	log.Printf("subscribed...")
	for d := range deliveries {
		log.Printf(
			"got %dB delivery: [%v] %q",
			len(d.Body),
			d.DeliveryTag,
			d.Body,
		)
		if c.close {
			log.Printf("handle: deliveries channel closed")
			// No new messages are processed after shutdown, and the message queue is notified to re-deliver the currently received messages
			_ = d.Nack(true, true)
			break
		}

		if c.cb != nil {
			err = c.cb.Fnc(Delivery{d})
			if err != nil && c.opt.ReReceiveNum > 1 {
				for i := 1; i < c.opt.ReReceiveNum; i++ {
					err := c.cb.Fnc(Delivery{d})
					if err == nil {
						break
					}
					if i == c.opt.ReReceiveNum-1 {
						_ = d.Nack(false, true)
					}
					time.Sleep(c.opt.ReReceiveDelay)
				}
			}
		}
	}

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

func (c *Consumer) GracefulShutdown() {
	if c.opt.GracefulDelay == 0 {
		c.opt.GracefulDelay = 10
	}
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	s := <-sc
	log.Println("Got signal:", s)
	c.close = true
	done := make(chan struct{})
	go func() {
		c.wg.Wait()
		done <- struct{}{}
	}()
	ctx, cancel := context.WithTimeout(context.Background(), c.opt.GracefulDelay*time.Second)
	defer cancel()
	select {
	case <-ctx.Done():
		_ = c.GetChannel().Close()
		log.Printf("RabbitMQ consumer did not shut down after %d seconds \n", c.opt.GracefulDelay)
	case <-done:
		_ = c.GetChannel().Close()
	}
}
