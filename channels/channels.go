package channels

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ysk229/go-rabbitmq/bindings"
	"github.com/ysk229/go-rabbitmq/exchanges"
	"github.com/ysk229/go-rabbitmq/lib"
	"github.com/ysk229/go-rabbitmq/queues"
	"log"
)

type Channel struct {
	*amqp.Channel
	conn        *amqp.Connection
	confirmChan chan amqp.Confirmation
	returnChan  chan amqp.Return
}

func NewChannel(conn *amqp.Connection) *Channel {
	ch, _ := conn.Channel()

	return &Channel{Channel: ch, conn: conn}
}

func (c *Channel) GetChannel() *amqp.Channel {
	return c.Channel
}

func (c *Channel) NewChannel() *Channel {
	ch, err := c.conn.Channel()
	if err != nil {
		log.Println(err)
	}
	c.Channel = ch
	return c
}

func (c *Channel) ConfirmOne(confirm chan amqp.Confirmation) chan amqp.Confirmation {
	if err := c.GetChannel().Confirm(false); err != nil {
		close(confirm) // confirms not supported, simulate by always nacking
	} else {
		c.GetChannel().NotifyPublish(confirm)
	}
	return confirm
}

func (c *Channel) GetConfirmChan() chan amqp.Confirmation {
	return c.confirmChan
}

func (c *Channel) SetConfirmChan(confirm chan amqp.Confirmation) {
	c.confirmChan = confirm
}

func (c *Channel) GetReturnChan() chan amqp.Return {
	return c.returnChan
}

func (c *Channel) SetReturnChan(notify chan amqp.Return) {
	c.returnChan = notify
}

func (c *Channel) Close() {
	err := c.GetChannel().Close()
	if err != nil {
		log.Printf("unexpected error during connection close: %v", err)
		return
	}
}

func (c *Channel) CloseAll() {
	_ = c.GetChannel().Close()
	_ = c.conn.Close()
}

func (c *Channel) Exchange(exchange *exchanges.Exchange) (string, error) {
	err := c.ExchangeDeclare(exchange.ExchangeName, exchange.ExchangeType.String(), exchange.Durable, exchange.AutoDelete, exchange.Internal, exchange.NoWait, lib.WrapTable(exchange.Args))
	if err != nil {
		return "", err
	}
	return exchange.ExchangeName, nil
}

func (c *Channel) Queue(q *queues.Queue) (amqp.Queue, error) {
	return c.QueueDeclare(q.QueueName, q.Durable, q.AutoDelete, q.Exclusive, q.NoWait, lib.WrapTable(q.Args))
}

func (c *Channel) BindingQueue(q *bindings.Binding) error {
	return c.QueueBind(q.Queue, q.RoutingKey, q.Exchange, q.NoWait, lib.WrapTable(q.Args))
}
