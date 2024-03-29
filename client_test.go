package rabbitmq

import (
	"fmt"
	"github.com/ysk229/go-rabbitmq/v2/bindings"
	"github.com/ysk229/go-rabbitmq/v2/channels"
	"github.com/ysk229/go-rabbitmq/v2/consumers"
	"github.com/ysk229/go-rabbitmq/v2/exchanges"
	"github.com/ysk229/go-rabbitmq/v2/lib"
	"github.com/ysk229/go-rabbitmq/v2/msg"
	"github.com/ysk229/go-rabbitmq/v2/options"
	"github.com/ysk229/go-rabbitmq/v2/producers"
	"github.com/ysk229/go-rabbitmq/v2/queues"
	"log"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	//new client mq
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", "admin", "123456", "127.0.0.1", 5672, "")
	mq := NewClient(url)
	d := make(chan string, 10)
	go func() {
		for i := 0; i < 3; i++ {
			d <- fmt.Sprintf("this is chan %d", i)
			//time.Sleep(3*time.Second)
		}
		close(d)
	}()

	p := mq.GetProducer()
	for c := range d {
		go p.Producer(
			msg.NewMessage(
				msg.WithOptionsChannel(channels.NewChannel(mq.Connection)),
				msg.WithOptionsBody("sdfdsfd"+c),
			),
			producers.WithOptionsProducer(&producers.ProducerOpt{
				Exchange:     "exchangeName",
				ExchangeType: lib.Topic,
				RouteKey:     "routeKey",
				Mandatory:    true,
				ResendNum:    2,
			}),
			producers.WithOptionsProducerCallBack(&producers.CallBack{Fnc: func(ret msg.Ret) {
				log.Printf("call back %+v", ret)

			}}),
		)

	}

	//select {}
}

func TestConsumer(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	//new client mq
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", "admin", "123456", "127.0.0.1", 5672, "")
	mq := NewClient(url)
	//new mq channel
	channelClient := mq.GetChan()
	exchangeName, err := channelClient.Exchange(exchanges.NewExchange(&options.Exchange{ExchangeName: "exchangeName", ExchangeType: lib.Topic}))
	if err != nil {
		log.Println(err)
	}
	q, err := channelClient.Queue(queues.NewQueue(&options.Queue{QueueName: "testestet"}))
	if err != nil {
		log.Println(err)
	}
	routeKey := "routeKey"
	err = channelClient.BindingQueue(bindings.NewBinding(&options.QueueBind{Queue: q.Name, Exchange: exchangeName, RoutingKey: routeKey}))
	if err != nil {
		log.Println(err)
	}
	go func() {
		mq.GetConsumer().Consumer(
			channelClient,
			consumers.WithOptionsConsumer(
				&consumers.ConsumerOpt{QueueName: q.Name, RoutingKey: routeKey, Exchange: exchangeName, ExchangeType: lib.Topic},
			),
			consumers.WithOptionsConsumerCallBack(
				&consumers.CallBack{Fnc: func(delivery consumers.Delivery) error {
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

func TestClientExchange(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	//new client mq
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", "admin", "123456", "127.0.0.1", 5672, "")
	mq := NewClient(url)
	//new mq channel
	channelClient := mq.GetChan()

	// new channel
	channelClient2 := mq.NewChan()

	////exchanges
	//err := channelClient.ExchangeDeclare("test-111111", "fanout", true, false, false, true, nil)
	//err := channelClient.Exchange(&exchanges.Exchange{Exchange:&options.Exchange{ExchangeName: "test-22zzzz", ExchangeType: lib.Direct,Durable: true,AutoDelete: false}})
	exchangeName, err := channelClient.Exchange(exchanges.NewExchange(&options.Exchange{ExchangeName: "go-test", ExchangeType: lib.Topic}))
	if err != nil {
		log.Println(err)
	}
	//log.Println(exchangeName)
	////queues
	q, err := channelClient.Queue(queues.NewQueue(&options.Queue{QueueName: "go-test"}))
	if err != nil {
		log.Println(err)
	}
	log.Println(q.Name)

	q2, err := channelClient2.Queue(queues.NewQueue(&options.Queue{QueueName: "testestet2"}))
	if err != nil {
		log.Println(err)
	}
	routeKey := "go-test"
	err = channelClient.BindingQueue(bindings.NewBinding(&options.QueueBind{Queue: q.Name, Exchange: exchangeName, RoutingKey: routeKey}))
	if err != nil {
		log.Println(err)
	}
	log.Println(q2.Name)
	//////more channel
	log.Println(channelClient)
	log.Println(channelClient2)
}

func TestClientProducers(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	//new client mq
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", "admin", "123456", "127.0.0.1", 5672, "")
	mq := NewClient(url)
	d := make(chan string, 10)
	go func() {
		for i := 0; i < 10; i++ {
			d <- fmt.Sprintf("this is chan %d", i)
			//time.Sleep(3*time.Second)
		}
		close(d)
	}()

	p := mq.NewChanProducer()
	for c := range d {
		p.Producer(
			msg.NewMessage(
				//msg.WithOptionsChannel(channels.NewChannel(mq.Connection)),
				msg.WithOptionsChannel(p.Channel),
				msg.WithOptionsBody("sdfdsfd"+c),
			),
			producers.WithOptionsProducer(&producers.ProducerOpt{
				Exchange:     "exchangeName",
				ExchangeType: lib.Topic,
				RouteKey:     "routeKey",
				Mandatory:    true,
				ResendNum:    2,
			}),
			producers.WithOptionsProducerCallBack(&producers.CallBack{Fnc: func(ret msg.Ret) {
				log.Printf("call back %+v", ret)

			}}),
		)

	}

	//select {}
}
