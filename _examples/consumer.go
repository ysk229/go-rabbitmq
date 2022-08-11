package main

import (
	"flag"
	"github.com/ysk229/go-rabbitmq/v2"
	"github.com/ysk229/go-rabbitmq/v2/consumers"
	"github.com/ysk229/go-rabbitmq/v2/lib"
	"log"
	"time"
)

var (
	url        = flag.String("uri", "amqp://admin:123456@127.0.0.1:5672", "AMQP URI")
	exchange   = flag.String("exchange", "test-exchange7", "Durable, non-auto-deleted AMQP exchange name")
	queue      = flag.String("queue", "test-queue", "Ephemeral AMQP queue name")
	bindingKey = flag.String("key", "test-key2", "AMQP binding key")
)

func init() {

	flag.Parse()
}

func main() {
	mq := rabbitmq.NewClient(*url)
	mq.GetConsumer().Consumer(
		mq.GetChan(),
		consumers.WithOptionsConsumer(
			&consumers.ConsumerOpt{QueueName: *queue, RoutingKey: *bindingKey, Exchange: *exchange, ExchangeType: lib.Topic},
		),
		consumers.WithOptionsConsumerCallBack(
			&consumers.CallBack{Fnc: func(delivery consumers.Delivery) error {
				time.Sleep(3 * time.Second)
				if delivery.DeliveryTag == 1 {
					_ = delivery.Ack(false)
				} else {
					_ = delivery.Nack(false, false)
				}
				log.Printf("%+v", delivery)
				return nil
			},
			},
		),
	)
}
