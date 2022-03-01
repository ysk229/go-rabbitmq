package rabbitmq

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ysk229/go-rabbitmq/channels"
	"github.com/ysk229/go-rabbitmq/connections"
	"github.com/ysk229/go-rabbitmq/consumers"
	"github.com/ysk229/go-rabbitmq/producers"
	"log"
)

type Client struct {
	*connections.Connect
	*channels.Channel
	*producers.Producer
	closeConnectionChan chan error
	channels            []*channels.Channel
}

func NewClient(url string) *Client {
	conn := connections.NewConnect().Open(url)
	c := &Client{Connect: conn, Channel: channels.NewChannel(conn.Connection), closeConnectionChan: make(chan error)}
	go c.healthCheck()
	go c.ReConnection(url)
	return c
}

func (c *Client) ReConnection(url string) {
	for {
		if c.GetConnect().IsClosed() || c.GetChannel().IsClosed() {
			select {
			case <-c.ReConnectionChan():
				log.Print("reconnecting ...")
				c.Connect = connections.NewConnect().Open(url)
				c.Channel = channels.NewChannel(c.Connection)
				c.ReChannels()
				log.Println("Connect is closed ", c.GetConnect().IsClosed(), ",Channel is closed ", c.GetChannel().IsClosed())
			}
		}
	}
}

func (c *Client) healthCheck() {
	closeChan := c.GetChan().GetChannel().NotifyClose(make(chan *amqp.Error))
	closeConnChan := c.GetConnect().NotifyClose(make(chan *amqp.Error))
	cancelChan := c.GetChan().GetChannel().NotifyCancel(make(chan string))

	defer close(closeChan)
	defer close(cancelChan)
	defer close(closeConnChan)

	for {
		select {
		case <-closeChan:
			log.Printf("closeChan the amqp server %s cannot connect", c.GetConnect().LocalAddr())
			c.closeConnectionChan <- fmt.Errorf("the amqp server %s cannot connect", c.GetConnect().LocalAddr())
		case <-closeConnChan:
			log.Printf("closeConnChan the amqp server %s cannot connect", c.GetConnect().LocalAddr())
			c.closeConnectionChan <- fmt.Errorf("the amqp server %s cannot connect", c.GetConnect().LocalAddr())
		case <-cancelChan:
			log.Printf("the queue of amqp has been canceled")
			c.closeConnectionChan <- fmt.Errorf("the queue of amqp has been canceled")
		}
	}
}
func (c *Client) ReConnectionChan() chan error {
	return c.closeConnectionChan
}

func (c *Client) GetChannels() []*channels.Channel {
	return c.channels
}
func (c *Client) ReChannels() {
	for index := range c.channels {
		log.Println("reconnecting channel", index)
		channels.NewChannel(c.Connect.Connection)
	}
}

func (c *Client) GetConnect() *connections.Connect {
	return c.Connect
}
func (c *Client) GetChan() *channels.Channel {
	return c.Channel
}

func (c *Client) NewChan() *channels.Channel {
	ch := channels.NewChannel(c.Connect.Connection)
	c.channels = append(c.channels, ch)
	return ch
}

func (c *Client) GetProducer() *producers.Producer {

	return producers.NewProducer(c.Channel)
}
func (c *Client) NewChanProducer() *producers.Producer {

	return producers.NewProducer(c.NewChan())
}

func (c *Client) GetConsumer() *consumers.Consumer {
	return consumers.NewConsumer(c.Channel)
}

func (c *Client) NewChanConsumer() *consumers.Consumer {
	return consumers.NewConsumer(c.NewChan())
}

func (c *Client) RPC() {
}
