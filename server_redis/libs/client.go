package libs

import (
	"encoding/json"
	"log"
	"net"
	"server_redis/structure"
	"sync"

	"github.com/gobwas/ws/wsutil"
	"github.com/rcrowley/go-metrics"
)

type Client struct {
	Lock       *sync.Mutex
	Connection net.Conn
	Name       string
	meter      metrics.Meter
}

func NewClient(connection net.Conn, meter metrics.Meter) *Client {
	return &Client{Lock: &sync.Mutex{}, Name: "", Connection: connection, meter: meter}
}

func (c *Client) Response(error_text string, nonce string) {
	type_of := structure.NotificationResponse
	notification := &structure.Notification{Type: type_of, Error: error_text, Nonce: nonce}
	message, err := json.Marshal(notification)
	if err != nil {
		log.Println("SERVER ERROR: json.Marshal")
		return
	}

	c.Write(message)
}

func (c *Client) Write(message []byte) {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	wsutil.WriteServerText(c.Connection, message)
	c.meter.Mark(1)
}
