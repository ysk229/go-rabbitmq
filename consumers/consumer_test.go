package consumers

import (
	"fmt"
	"github.com/ysk229/go-rabbitmq/channels"
	"github.com/ysk229/go-rabbitmq/connections"
	"github.com/ysk229/go-rabbitmq/lib"
	"log"
	"testing"
	"time"
)

func TestConsumer(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	//new client mq
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", "admin", "123456", "127.0.0.1", 5672, "")
	conn := connections.NewConnect().Open(url)
	//new mq channel
	channelClient := channels.NewChannel(conn.Connection)
	exchangeName := "go-test"
	routeKey := "go-test"
	q := "go-test"
	//i := 0
	NewConsumer(channelClient).Consumer(
		channelClient,
		WithOptionsConsumer(
			&ConsumerOpt{QueueName: q, RoutingKey: routeKey, Exchange: exchangeName, ExchangeType: lib.Topic},
		),
		WithOptionsConsumerCallBack(
			&CallBack{Fnc: func(delivery Delivery) {

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
	//select {}
}
