package main

import (
	"flag"
	"github.com/ysk229/go-rabbitmq"
	"github.com/ysk229/go-rabbitmq/consumers"
	"github.com/ysk229/go-rabbitmq/lib"
	"log"
	"time"
)

var (
	url        = flag.String("uri", "amqp://admin:123456@10.1.2.7:5672/", "AMQP URI")
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
			&consumers.CallBack{Fnc: func(delivery consumers.Delivery) {
				time.Sleep(3 * time.Second)
				if delivery.DeliveryTag == 1 {
					_ = delivery.Ack(false)
				} else {
					_ = delivery.Nack(false, false)
				}
				log.Printf("%+v", delivery)
			},
			},
		),
	)
}
