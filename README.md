# go-rabbitmq

Wrapper of [rabbitmq/amqp091-go](https://github.com/rabbitmq/amqp091-go) that provides reconnection logic and sane defaults. Hit the project with a star if you find it useful ⭐


[![Github License](https://img.shields.io/github/license/ysk229/go-rabbitmq.svg?style=flat)](https://github.com/ysk229/go-rabbitmq/blob/master/LICENSE)
[![Go Doc](https://godoc.org/github.com/ysk229/go-rabbitmq?status.svg)](https://pkg.go.dev/github.com/ysk229/go-rabbitmq)
[![Go Report](https://goreportcard.com/badge/github.com/ysk229/go-rabbitmq)](https://goreportcard.com/report/github.com/ysk229/go-rabbitmq)
[![Github Latest Release](https://img.shields.io/github/release/ysk229/go-rabbitmq.svg?style=flat)](https://github.com/ysk229/go-rabbitmq/releases/latest)
[![Github Latest Tag](https://img.shields.io/github/tag/ysk229/go-rabbitmq.svg?style=flat)](https://github.com/sliveryou/arabbitmq/tags)
[![Github Stars](https://img.shields.io/github/stars/ysk229/go-rabbitmq.svg?style=flat)](https://github.com/ysk229/go-rabbitmq/stargazers)

## Motivation

[Streadway's AMQP](https://github.com/rabbitmq/amqp091-go) library is currently the most robust and well-supported Go client I'm aware of. It's a fantastic option and I recommend starting there and seeing if it fulfills your needs. Their project has made an effort to stay within the scope of the AMQP protocol, as such, no reconnection logic and few ease-of-use abstractions are provided.

### Goal

The goal with `go-rabbitmq` is to still provide most all of the nitty-gritty functionality of AMQP, but to make it easier to work with via a higher-level API. Particularly:

* Automatic reconnection
* Multithreaded consumers via a handler function
* Reasonable defaults
* Flow control handling

## ⚙️ Installation

1. Create rabbitmq docker container by using:
```bash
$ docker run --name rabbitmq --hostname rabbitmq-test-node-1 -p 15672:15672 -p 5672:5672 -e RABBITMQ_DEFAULT_USER=root -e RABBITMQ_DEFAULT_PASS=123123 -d rabbitmq:3.8.5-management
```

2. Download rabbitmq package by using:
```bash
go get github.com/ysk229/go-rabbitmq
```

## 🚀 Quick Start Consumer

### Default options

```go
mq := NewClient("amqp://user:pass@localhost") 
exchangeName := "go-test"
routeKey := "go-test"
q := "go-test"
////消费者，确认是否消费成功
//chanClient.Consumer().CallBack()
mq.GetConsumer().Consumer(
    mq.GetChan(),
    consumers.WithOptionsConsumer(
        &consumers.ConsumerOpt{
			QueueName: q,
			RoutingKey: routeKey,
			Exchange: exchangeName,
			ExchangeType: lib.Topic},
    ),
    consumers.WithOptionsConsumerCallBack(
        &consumers.CallBack{
			Fnc: func(delivery consumers.Delivery) {
				 log.Printf("%+v", delivery)
            },
        },
    ),
)
select {}
```


## 🚀 Quick Start Publisher

### Default options

```go

NewClient("amqp://user:pass@localhost").GetProducer().Producer(
    msg.NewMessage(
        msg.WithOptionsChannel(channels.NewChannel(mq.Connection)),
        msg.WithOptionsBody("sdfdsfd"+c), 
	),
    producers.WithOptionsProducer(&producers.ProducerOpt{
    Exchange:     "exchangeName",
    ExchangeType: lib.Topic,
    RouteKey:     "routeKey",
    Mandatory:    true,
    ResendNum:    2}),
    producers.WithOptionsProducerCallBack(&producers.CallBack{
		Fnc: func(ret msg.Ret) {
            log.Printf("call back %+v", ret)
        }
	}),
)

}

select {}

```


## Other usage examples

See the [examples](_examples) directory for more ideas.



## Transient Dependencies

My goal is to keep dependencies limited to 1, [github.com/rabbitmq/amqp091-go](https://github.com/rabbitmq/amqp091-go).

## 👏 Contributing

I love help! Contribute by forking the repo and opening pull requests. Please ensure that your code passes the existing tests and linting, and write tests to test your changes if applicable.

All pull requests should be submitted to the `main` branch.
