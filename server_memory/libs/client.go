package libs

import (
	"encoding/json"
	"log"
	"net"
	"server_memory/structure"
	"sync"

	"github.com/gobwas/ws/wsutil"
)

type Client struct {
	Lock       *sync.Mutex
	Connection net.Conn
	Name       string
	unregister chan net.Conn
	counter    *Counter
}

func NewClient(connection net.Conn, unregister chan net.Conn, counter *Counter) *Client {
	return &Client{Lock: &sync.Mutex{}, Name: "", Connection: connection, unregister: unregister, counter: counter}
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
	c.counter.Increment()
}
