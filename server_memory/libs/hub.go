package libs

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"server_memory/structure"
	"time"

	"github.com/panjf2000/ants/v2"
	"golang.org/x/sys/unix"
)

type JoinPack struct {
	Client *Client
	Name   string
	Nonce  string
}

type RequestPack struct {
	Connection net.Conn
	HandleFunc func(*Client) func()
}

type BroadcastPack struct {
	HandleFunc func(*Client) func()
}

type Hub struct {
	information *Information
	epoll       *Epoll
	go_pool     *ants.Pool

	connections map[net.Conn]*Client

	Register   chan *Client
	Unregister chan net.Conn
	Join       chan *JoinPack
	Process    chan *RequestPack
	Broadcast  chan *BroadcastPack
}

func NewHub(information *Information, epoll *Epoll, go_pool *ants.Pool) *Hub {
	return &Hub{
		information: information,
		epoll:       epoll,
		go_pool:     go_pool,

		connections: make(map[net.Conn]*Client),

		Register:   make(chan *Client, 1000),
		Unregister: make(chan net.Conn, 1000),
		Join:       make(chan *JoinPack, 1000),
		Process:    make(chan *RequestPack, 1000),
		Broadcast:  make(chan *BroadcastPack, 1000),
	}
}

func (h *Hub) Run() {
	for {
		select {

		case client := <-h.Register:
			h.connections[client.Connection] = client

		case connection := <-h.Unregister:
			if client, ok := h.connections[connection]; ok {
				h.removeFromEpoll(client.Connection)
				h.removeConnection(client)
				h.sendLeaveMessage(client.Name)
			}

		case join_pack := <-h.Join:
			if client, ok := h.connections[join_pack.Client.Connection]; ok {
				if !h.information.AddUser(join_pack.Name) {
					client.Response(structure.NameError, join_pack.Nonce)
					break
				}

				client.Name = join_pack.Name

				content := &structure.Content{From: client.Name}
				notification := &structure.Notification{Type: structure.NotificationJoin, Content: content, Nonce: join_pack.Nonce}
				if err := h.marshalAndBroadcast(notification); err != nil {
					break
				}

				client.Response(structure.Success, notification.Nonce)
			}

		case broadcast_pack := <-h.Broadcast:
			h.broadcast(broadcast_pack)

		case message_pack := <-h.Process:
			if client, ok := h.connections[message_pack.Connection]; ok {
				h.go_pool.Submit(message_pack.HandleFunc(client))
			}
		}
	}
}

func (h *Hub) sendLeaveMessage(name string) {
	if name == "" {
		return
	}

	content := &structure.Content{From: name}
	notification := &structure.Notification{Type: structure.NotificationLeave, Content: content, Nonce: ""}
	h.marshalAndBroadcast(notification)
}

func (h *Hub) marshalAndBroadcast(notification *structure.Notification) error {
	message, err := json.Marshal(notification)
	if err != nil {
		log.Println("SERVER ERROR: json.Marshal")
		return err
	}

	broadcast_pack := &BroadcastPack{HandleFunc: func(c *Client) func() {
		return func() { c.Write(message) }
	}}

	h.broadcast(broadcast_pack)

	return nil
}

func (h *Hub) broadcast(broadcast_pack *BroadcastPack) {
	for _, client := range h.connections {
		h.go_pool.Submit(broadcast_pack.HandleFunc(client))
	}
}

func (h *Hub) removeFromEpoll(connection net.Conn) {
	max_retries := 10
	for retries := 0; retries < max_retries; retries += 1 {
		if err := h.epoll.Remove(connection); err == nil || errors.Is(err, unix.ENOENT) || errors.Is(err, unix.EBADF) {
			break
		} else {
			log.Printf("Failed to remove from epoll %v (attempt %d)", err, retries+1)
		}
		time.Sleep(time.Millisecond * 10)
	}
}

func (h *Hub) removeConnection(client *Client) {
	h.information.RemUser(client.Name)
	delete(h.connections, client.Connection)
	client.Connection.Close()
}

func (h *Hub) Close() {
	for connection := range h.connections {
		connection.Close()
	}
}
