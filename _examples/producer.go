package main

import (
	"flag"
	"github.com/ysk229/go-rabbitmq/v2"
	"github.com/ysk229/go-rabbitmq/v2/channels"
	"github.com/ysk229/go-rabbitmq/v2/lib"
	"github.com/ysk229/go-rabbitmq/v2/msg"
	"github.com/ysk229/go-rabbitmq/v2/producers"
	"log"
)

var (
	uri          = flag.String("uri", "amqp://admin:123456@127.0.0.1:5672", "AMQP URI")
	exchangeName = flag.String("exchange", "test-exchange7", "Durable AMQP exchange name")
	routingKey   = flag.String("key", "test-key2", "AMQP routing key")
	body         = flag.String("body", "foobar", "Body of message")
)

func init() {
	flag.Parse()
}

func main() {
	mq := rabbitmq.NewClient(*uri)
	p := mq.GetProducer()
	p.Producer(
		msg.NewMessage(
			msg.WithOptionsChannel(channels.NewChannel(mq.Connection)),
			msg.WithOptionsBody(*body),
		),
		producers.WithOptionsProducer(&producers.ProducerOpt{
			Exchange:     *exchangeName,
			ExchangeType: lib.Topic,
			RouteKey:     *routingKey,
			Mandatory:    true,
			ResendNum:    2,
		}),
		producers.WithOptionsProducerCallBack(&producers.CallBack{Fnc: func(ret msg.Ret) {
			log.Printf("call back %+v", ret)

		}}),
	)
}
