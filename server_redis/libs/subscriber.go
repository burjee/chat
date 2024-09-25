package libs

import (
	"context"
	"encoding/base64"
	"log"

	"github.com/redis/go-redis/v9"
)

type Subscriber struct {
	redis_client *redis.Client
	hub          *Hub
}

func NewSubscriber(redis_client *redis.Client, hub *Hub) *Subscriber {
	return &Subscriber{redis_client, hub}
}

func (s *Subscriber) Sub() {
	pubsub := s.redis_client.Subscribe(context.Background(), "broadcast")
	if _, err := pubsub.Receive(context.Background()); err != nil {
		panic("subscriber error")
	}

	go func() {
		for msg := range pubsub.Channel() {
			decoded_data, err := base64.StdEncoding.DecodeString(msg.Payload)
			if err != nil {
				log.Printf("Failed to decode message: %v", err)
				continue
			}

			s.hub.Broadcast <- &BroadcastPack{HandleFunc: func(c *Client) func() {
				return func() { c.Write(decoded_data) }
			}}
		}
	}()
}
