package main

import (
	"github.com/ysk229/go-rabbitmq"
	"github.com/ysk229/go-rabbitmq/consumers"
	"github.com/ysk229/go-rabbitmq/lib"
	"log"
	"os"
	"time"
)

func main() {
	mq := rabbitmq.NewClient("amqp://admin:123456@rabbitmq:5672")
	if len(os.Args) < 2 {
		log.Printf("Usage: %s [binding_key]...", os.Args[0])
		os.Exit(0)
	}
	mq.GetConsumer().Consumer(
		mq.GetChan(),
		consumers.WithOptionsConsumer(
			&consumers.ConsumerOpt{QueueName: "", RoutingKey: os.Args[1:][0], Exchange: "logs_topic", ExchangeType: lib.Topic},
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
	forever := make(chan bool)
	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
