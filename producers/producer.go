package producers

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ysk229/go-rabbitmq/v2/channels"
	"github.com/ysk229/go-rabbitmq/v2/exchanges"
	"github.com/ysk229/go-rabbitmq/v2/lib"
	"github.com/ysk229/go-rabbitmq/v2/msg"
	"github.com/ysk229/go-rabbitmq/v2/options"
	"log"
	"time"
)

// ProducerOpt Producer Options
type ProducerOpt struct {
	// Mandatory makes the publishing mandatory, which means when a queue is not
	// bound to the routing key a message will be sent back on the returns channel for you to handle.
	Mandatory bool
	// Immediate makes the publishing immediate, which means when a consumer is not available
	// to immediately handle the new message, a message will be sent back on the returns channel for you to handle.
	Immediate    bool
	DeliveryMode bool
	Exchange     string
	ExchangeType lib.ExchangeType
	RouteKey     string
	ResendDelay  time.Duration //How many seconds to resend the message after it fails,default is 3s
	ResendNum    int           //Number of message resends, default is 3 times
}

// CallBack  Call Back
type CallBack struct {
	Fnc func(ret msg.Ret)
}

// ProducerOption Producer Option
type ProducerOption func(*Producer)

// Producer producer
type Producer struct {
	*channels.Channel
	opt *ProducerOpt
	cb  *CallBack
}

// NewProducer New Producer
func NewProducer(ch *channels.Channel) *Producer {
	p := &Producer{Channel: ch}

	return p
}

// WithOptionsProducer With Options Producer
func WithOptionsProducer(opt *ProducerOpt) ProducerOption {
	return func(p *Producer) {
		p.opt = opt
		if p.opt.ResendDelay == 0 {
			p.opt.ResendDelay = 3
		}
		if p.opt.ResendNum == 0 {
			p.opt.ResendNum = 3
		}
	}
}

// WithOptionsProducerCallBack With Options Producer CallBack
func WithOptionsProducerCallBack(cb *CallBack) ProducerOption {
	return func(p *Producer) {
		p.cb = cb
	}
}

// Producer producer
func (p *Producer) Producer(m *msg.Message, opts ...ProducerOption) {
	for _, opt := range opts {
		opt(p)
	}
	ch := m.Channel
	if ch.GetConfirmChan() == nil {
		ch.SetConfirmChan(ch.ConfirmOne(make(chan amqp.Confirmation, 1)))
		ch.SetReturnChan(ch.NotifyReturn(make(chan amqp.Return, 1)))
	}
	exchangeNum := 0
	queueNum := 0
	body := string(m.Body)
	Success := true
LOOKUP:
	for {
		if exchangeNum > p.opt.ResendNum || queueNum > p.opt.ResendNum {
			break LOOKUP
		}
		p.publish(m)
		cf := <-ch.GetConfirmChan()
		if !cf.Ack {
			log.Println("exchange error  ", body)
			exchangeNum++
			if exchangeNum == p.opt.ResendNum+1 {
				log.Println("exchange fail data", body)
				Success = false
			}
		} else {
			select {
			case data := <-ch.GetReturnChan():
				log.Println("queue error  data", body, ",ReplyCode ", data.ReplyCode, ",retry num", queueNum, ",return data ", string(data.Body))
				queueNum++
				if queueNum == p.opt.ResendNum+1 {
					log.Println("queue fail data", body)
					Success = false
				}
			default:
				break LOOKUP
			}
		}

		time.Sleep(p.opt.ResendDelay * time.Second)
	}
	ret := msg.Ret{Success: Success, Data: body, Exchange: p.opt.Exchange, RoutingKey: p.opt.RouteKey}
	if p.cb != nil {
		p.cb.Fnc(ret)
	} else {
		log.Printf("%+v", ret)
	}
}

func (p *Producer) publish(msg *msg.Message) {
	ch := msg.Channel
	_, _ = ch.Exchange(exchanges.NewExchange(&options.Exchange{ExchangeName: p.opt.Exchange, ExchangeType: p.opt.ExchangeType}))

	err := ch.Publish(p.opt.Exchange, p.opt.RouteKey,
		p.opt.Mandatory, // mandatory, true:若没有一个队列与交换器绑定，则将消息返还给生产者 , false:若交换器没有匹配到队列，消息直接丢弃
		p.opt.Immediate, // immediate , true:队列没有对应的消费者，则将消息返还给生产者,
		amqp.Publishing{
			ContentType:  msg.ContentType,
			Body:         msg.Body,
			DeliveryMode: msg.DeliveryMode,
			Headers:      lib.WrapTable(msg.Headers),
			Expiration:   msg.Expiration,
		},
	)

	if err != nil || ch.GetChannel().IsClosed() {
		log.Printf("%p Channel Publish %v err %v %p\n", ch, ch.GetChannel().IsClosed(), err, ch.GetChannel())
		p.reExchange(ch)
	}

}

func (p *Producer) reExchange(ch *channels.Channel) {
	if ch.GetChannel().IsClosed() {
		c := ch.GetChannel()
		_, _ = ch.NewChannel().Exchange(exchanges.NewExchange(&options.Exchange{ExchangeName: p.opt.Exchange, ExchangeType: p.opt.ExchangeType}))
		log.Println(fmt.Sprintf("old channel %p", c), "new channel", fmt.Sprintf("%p", ch.GetChannel()))
		go ch.SetConfirmChan(ch.ConfirmOne(make(chan amqp.Confirmation, 1)))
		go ch.SetReturnChan(ch.NotifyReturn(make(chan amqp.Return, 1)))
	}
}
