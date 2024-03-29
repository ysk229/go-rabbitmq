package connections

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ysk229/go-rabbitmq/v2/channels"
	"log"
)

// Connect connect struct
type Connect struct {
	*amqp.Connection
	*channels.Channel
	//连接异常结束
	ConnNotifyClose chan *amqp.Error
	//通道异常接收
	ChNotifyClose chan *amqp.Error
}

// NewConnect new connect
func NewConnect() *Connect {
	return &Connect{}
}

// Open mq connections
// url 	"amqp://user:password@host:port/vhost"
func (c *Connect) Open(url string) *Connect {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalf("cannot dial integration server. Is the rabbitmq-server service running? %s", err)
	}
	c.Connection = conn
	return c
}

//Close connect
func (c *Connect) Close() {
	err := c.Connection.Close()
	if err != nil {
		log.Fatalf("connection close: %s", err)
		return
	}
}
