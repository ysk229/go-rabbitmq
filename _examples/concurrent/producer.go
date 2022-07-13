package main

import (
	"flag"
	"fmt"
	"github.com/ysk229/go-rabbitmq"
	"github.com/ysk229/go-rabbitmq/channels"
	"github.com/ysk229/go-rabbitmq/lib"
	"github.com/ysk229/go-rabbitmq/msg"
	"github.com/ysk229/go-rabbitmq/producers"
	"log"
	"time"
)

var (
	uri          = flag.String("uri", "amqp://admin:123456@10.1.2.7:5672", "AMQP URI")
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
	job := make(chan string, 15)
	//10 worker
	for i := 0; i < 10; i++ {
		go func(job <-chan string) {
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
		}(job)
	}
	for i := 0; i < 15; i++ {
		job <- fmt.Sprintf("this is chan %d", i)
	}
	close(job)
	<-time.After(30 * time.Second)
}
