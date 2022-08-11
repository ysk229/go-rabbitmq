package consumers

import (
	"errors"
	"fmt"
	"github.com/ysk229/go-rabbitmq/v2/channels"
	"github.com/ysk229/go-rabbitmq/v2/connections"
	"github.com/ysk229/go-rabbitmq/v2/lib"
	"log"
	"os"
	"syscall"
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
	go func() {
		NewConsumer(channelClient).Consumer(
			channelClient,
			WithOptionsConsumer(
				&ConsumerOpt{QueueName: q, RoutingKey: routeKey, Exchange: exchangeName, ExchangeType: lib.Topic},
			),
			WithOptionsConsumerCallBack(
				&CallBack{Fnc: func(delivery Delivery) error {

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
	}()

	log.Printf("running for %s", "10s")
	time.Sleep(10 * time.Second)
}
func TestRetryConsumer(t *testing.T) {
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
	go func() {
		NewConsumer(channelClient).Consumer(
			channelClient,
			WithOptionsConsumer(
				&ConsumerOpt{QueueName: q, RoutingKey: routeKey, Exchange: exchangeName, ExchangeType: lib.Topic},
			),
			WithOptionsConsumerCallBack(
				&CallBack{Fnc: func(delivery Delivery) error {

					time.Sleep(3 * time.Second)
					if delivery.DeliveryTag == 1 {
						_ = delivery.Ack(false)
					} else {
						_ = delivery.Nack(false, false)
					}
					log.Printf("%+v", delivery)
					return errors.New("error is nil then you can't retry, retry three times by default")
				},
				},
			),
		)
	}()

	log.Printf("running for %s", "10s")
	time.Sleep(10 * time.Second)
}
func TestGracefulShutdownConsumer(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	//new client mq
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", "admin", "123456", "127.0.0.1", 5672, "")
	conn := connections.NewConnect().Open(url)
	//new mq channel
	channelClient := channels.NewChannel(conn.Connection)
	exchangeName := "go-test"
	routeKey := "go-test"
	q := "go-test"
	c := NewConsumer(channelClient)
	//Graceful exit
	go func() {
		c.Consumer(
			channelClient,
			WithOptionsConsumer(
				&ConsumerOpt{QueueName: q, RoutingKey: routeKey, Exchange: exchangeName, ExchangeType: lib.Topic},
			),
			WithOptionsConsumerCallBack(
				&CallBack{Fnc: func(delivery Delivery) error {

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
	}()

	go func() {
		time.Sleep(5 * time.Second)
		p, _ := os.FindProcess(os.Getpid())
		log.Println(p.Pid)
		err := p.Signal(syscall.SIGINT)
		if err != nil {
			log.Fatalln(err)
		}
	}()

	c.GracefulShutdown()
}
func TestConcurrentConsumer(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	//new client mq
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", "admin", "123456", "127.0.0.1", 5672, "")
	conn := connections.NewConnect().Open(url)
	exchangeName := "go-test"
	routeKey := "go-test"
	q := "go-test"
	job := make(chan string, 15)
	//10 worker
	for i := 0; i < 10; i++ {
		go func(job <-chan string) {
			//new mq channel
			channelClient := channels.NewChannel(conn.Connection)
			NewConsumer(channelClient).Consumer(
				channelClient,
				WithOptionsConsumer(
					&ConsumerOpt{QueueName: q, RoutingKey: routeKey, Exchange: exchangeName, ExchangeType: lib.Topic},
				),
				WithOptionsConsumerCallBack(
					&CallBack{Fnc: func(delivery Delivery) error {

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
		}(job)
	}

	for i := 0; i < 15; i++ {
		job <- fmt.Sprintf("this is chan %d", i)
	}
	close(job)
	<-time.After(30 * time.Second)
}
