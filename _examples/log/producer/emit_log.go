package main

import (
	"github.com/ysk229/go-rabbitmq/v2"
	"github.com/ysk229/go-rabbitmq/v2/channels"
	"github.com/ysk229/go-rabbitmq/v2/lib"
	"github.com/ysk229/go-rabbitmq/v2/msg"
	"github.com/ysk229/go-rabbitmq/v2/producers"
	"log"
	"os"
	"strings"
)

func main() {
	mq := rabbitmq.NewClient("amqp://admin:123456@127.0.0.1:5672")
	p := mq.GetProducer()
	p.Producer(
		msg.NewMessage(
			msg.WithOptionsChannel(channels.NewChannel(mq.Connection)),
			msg.WithOptionsBody(bodyFromEmit(os.Args)),
		),
		producers.WithOptionsProducer(&producers.ProducerOpt{
			Exchange:     "logs",
			ExchangeType: lib.Fanout,
			RouteKey:     "",
			Mandatory:    true,
			ResendNum:    2,
		}),
		producers.WithOptionsProducerCallBack(&producers.CallBack{Fnc: func(ret msg.Ret) {
			log.Printf("call back %+v", ret)

		}}),
	)
}

func bodyFromEmit(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}
