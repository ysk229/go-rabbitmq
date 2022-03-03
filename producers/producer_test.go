package producers

import (
	"fmt"
	"github.com/ysk229/go-rabbitmq/channels"
	"github.com/ysk229/go-rabbitmq/connections"
	"github.com/ysk229/go-rabbitmq/lib"
	"github.com/ysk229/go-rabbitmq/msg"
	"log"
	"testing"
	"time"
)

func TestProducer(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	//new client mq
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", "admin", "123456", "127.0.0.1", 5672, "")
	conn := connections.NewConnect().Open(url)

	d := make(chan string, 10)
	go func() {
		for i := 0; i < 3; i++ {
			d <- fmt.Sprintf("this is chan %d", i)
			//time.Sleep(3*time.Second)
		}
		close(d)
	}()
	exchangeName := "go-test"
	routeKey := "go-test"
	p := NewProducer(channels.NewChannel(conn.Connection))
	for c := range d {
		go p.Producer(
			msg.NewMessage(
				msg.WithOptionsChannel(channels.NewChannel(conn.Connection)),
				msg.WithOptionsBody("sdfdsfd"+c),
			),
			WithOptionsProducer(&ProducerOpt{
				Exchange:     exchangeName,
				ExchangeType: lib.Topic,
				RouteKey:     routeKey,
				Mandatory:    true,
				ResendNum:    2,
			}),
			WithOptionsProducerCallBack(&CallBack{Fnc: func(ret msg.Ret) {
				log.Printf("call back %+v", ret)

			}}),
		)

	}
	//select {}
}

func TestProducerChannel(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	//new client mq
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", "admin", "123456", "127.0.0.1", 5672, "")
	conn := connections.NewConnect().Open(url)
	p := NewProducer(channels.NewChannel(conn.Connection))
	job := make(chan string, 15)
	//10 worker
	for i := 0; i < 10; i++ {
		go func(job <-chan string) {
			exchangeName := "go-test"
			routeKey := "go-test"
			// new 10 mq channel
			c := channels.NewChannel(conn.Connection)
			for body := range job {
				p.Producer(
					msg.NewMessage(
						msg.WithOptionsChannel(c),
						msg.WithOptionsBody(body),
					),
					WithOptionsProducer(&ProducerOpt{
						Exchange:     exchangeName,
						ExchangeType: lib.Topic,
						RouteKey:     routeKey,
						Mandatory:    true,
						ResendNum:    2,
					}),
					WithOptionsProducerCallBack(&CallBack{Fnc: func(ret msg.Ret) {
						log.Printf("call back %+v", ret)

					}}),
				)
			}

		}(job)
	}

	for i := 0; i < 15; i++ {
		job <- fmt.Sprintf("this is chan %d", i)
	}
	close(job)
	<-time.After(30 * time.Second)
}
