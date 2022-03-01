package producers

import (
	"fmt"
	"github.com/ysk229/go-rabbitmq/channels"
	"github.com/ysk229/go-rabbitmq/connections"
	"github.com/ysk229/go-rabbitmq/lib"
	"github.com/ysk229/go-rabbitmq/msg"
	"log"
	"testing"
)

func TestProducer(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	//new client mq
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", "admin", "123456", "rabbitmq", 5672, "")
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
	select {}
}
