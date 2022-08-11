package msg

import (
	"github.com/ysk229/go-rabbitmq/v2/channels"
	"github.com/ysk229/go-rabbitmq/v2/lib"
	"time"
)

// Message struct
type Message struct {
	// Application or exchange specific fields,
	// the headers exchange will inspect this field.
	Headers lib.Table

	// Properties
	ContentType     string    // MIME content type
	ContentEncoding string    // MIME content encoding
	DeliveryMode    uint8     // Transient (0 or 1) or Persistent (2)
	Priority        uint8     // 0 to 9
	CorrelationId   string    // correlation identifier
	ReplyTo         string    // address to to reply to (ex: RPC)
	Expiration      string    // message expiration spec
	MessageId       string    // message identifier
	Timestamp       time.Time // message timestamp
	Type            string    // message type name
	UserId          string    // creating user id - ex: "guest"
	AppId           string    // creating application id

	// The application specific payload of the message
	Body []byte

	Channel *channels.Channel
}

// MsgOption Msg Option
type MsgOption func(*Message)

// NewMessage New Message
func NewMessage(opts ...MsgOption) *Message {
	m := &Message{}
	for _, opt := range opts {
		opt(m)
	}
	m.Timestamp = time.Now()
	m.DeliveryMode = lib.Persistent
	return m
}

// WithOptionsChannel With Options Channel
func WithOptionsChannel(ch *channels.Channel) MsgOption {
	return func(m *Message) {
		m.Channel = ch
	}
}

// WithOptionsDeliveryMode With Options Delivery Mode
func WithOptionsDeliveryMode(DeliveryMode uint8) MsgOption {
	return func(m *Message) {
		m.DeliveryMode = DeliveryMode
	}
}

// WithOptionsBody With Options Body
func WithOptionsBody(Body string) MsgOption {
	return func(m *Message) {
		m.Body = []byte(Body)
	}
}
