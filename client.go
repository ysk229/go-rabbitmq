package rabbitmq

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ysk229/go-rabbitmq/channels"
	"github.com/ysk229/go-rabbitmq/connections"
	"github.com/ysk229/go-rabbitmq/consumers"
	"github.com/ysk229/go-rabbitmq/producers"
	"log"
	"time"
)

// Client  struct
type Client struct {
	*connections.Connect
	*channels.Channel
	*producers.Producer
	closeConnectionChan chan error
	channels            []*channels.Channel
}

// NewClient new client
func NewClient(url string) *Client {
	conn := connections.NewConnect().Open(url)
	c := &Client{Connect: conn, Channel: channels.NewChannel(conn.Connection), closeConnectionChan: make(chan error)}
	go c.healthCheck()
	go c.ReConnection(url)
	return c
}

// ReConnection Re Connection
func (c *Client) ReConnection(url string) {
	for {
		if c.GetConnect().IsClosed() || c.GetChannel().IsClosed() {
			<-c.ReConnectionChan()
			log.Print("reconnecting ...")
			c.Connect = connections.NewConnect().Open(url)
			c.Channel = channels.NewChannel(c.Connection)
			c.ReChannels()
			log.Println("Connect is closed ", c.GetConnect().IsClosed(), ",Channel is closed ", c.GetChannel().IsClosed())

		}
		time.Sleep(time.Second * time.Duration(5))
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

// ReConnectionChan Re Connection Chan
func (c *Client) ReConnectionChan() chan error {
	return c.closeConnectionChan
}

// GetChannels Get Channels
func (c *Client) GetChannels() []*channels.Channel {
	return c.channels
}

// ReChannels Re Channels
func (c *Client) ReChannels() {
	for index := range c.channels {
		log.Println("reconnecting channel", index)
		channels.NewChannel(c.Connect.Connection)
	}
}

// GetConnect Get Connect
func (c *Client) GetConnect() *connections.Connect {
	return c.Connect
}

// GetChan Get Chan
func (c *Client) GetChan() *channels.Channel {
	return c.Channel
}

// NewChan New Chan
func (c *Client) NewChan() *channels.Channel {
	ch := channels.NewChannel(c.Connect.Connection)
	c.channels = append(c.channels, ch)
	return ch
}

// GetProducer Get Producer
func (c *Client) GetProducer() *producers.Producer {
	return producers.NewProducer(c.Channel)
}

// NewChanProducer New Chan Producer
func (c *Client) NewChanProducer() *producers.Producer {
	return producers.NewProducer(c.NewChan())
}

// GetConsumer  Get Consumer
func (c *Client) GetConsumer() *consumers.Consumer {
	return consumers.NewConsumer(c.Channel)
}

// NewChanConsumer New Chan Consumer
func (c *Client) NewChanConsumer() *consumers.Consumer {
	return consumers.NewConsumer(c.NewChan())
}
