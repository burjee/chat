package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net"
	"server_redis/libs"
	"server_redis/structure"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sys/unix"
)

type service struct {
	redis_client *redis.Client
	information  *libs.Information
	epoll        *libs.Epoll
	hub          *libs.Hub
	validate     *validator.Validate
	counter      *libs.Counter
}

func New(redis_client *redis.Client, information *libs.Information, epoll *libs.Epoll, hub *libs.Hub, validate *validator.Validate, counter *libs.Counter) *service {
	return &service{redis_client, information, epoll, hub, validate, counter}
}

func (s *service) initWebsocketError(connection net.Conn, error_text string) {
	connection.Close()
	log.Printf("Failed to " + error_text)
}

func (s *service) WebsocketHandler(c *gin.Context) {
	connection, _, _, err := ws.UpgradeHTTP(c.Request, c.Writer)
	if err != nil {
		s.initWebsocketError(connection, "upgrade connection "+err.Error())
		return
	}

	room_list, user_list, err := s.information.GetList()
	if err != nil {
		s.initWebsocketError(connection, "get list "+err.Error())
		return
	}

	if err := wsutil.WriteServerText(connection, room_list); err != nil {
		s.initWebsocketError(connection, "write room list "+err.Error())
		return
	}

	if err := wsutil.WriteServerText(connection, user_list); err != nil {
		s.initWebsocketError(connection, "write user list "+err.Error())
		return
	}

	if err := s.epoll.Add(connection); err != nil {
		s.initWebsocketError(connection, "add connection "+err.Error())
		return
	}

	client := libs.NewClient(connection, s.hub.Unregister, s.counter)
	s.hub.Register <- client
}

func (s *service) StartReadWebsocket() {
	for {

		connections, err := s.epoll.Wait()
		if err != nil {
			if !errors.Is(err, unix.EINTR) {
				log.Printf("Failed to epoll wait %v", err)
			}
			continue
		}

		for _, connection := range connections {
			if connection == nil {
				break
			}

			if bytes, op, err := wsutil.ReadClientData(connection); err != nil {
				s.hub.Unregister <- connection
			} else {
				switch op {
				case ws.OpClose:
					s.hub.Unregister <- connection
				case ws.OpText:
					s.hub.Process <- &libs.RequestPack{Connection: connection, HandleFunc: func(client *libs.Client) func() { return func() { s.handleRequest(client, bytes) } }}
				}
			}
		}
	}
}

func (s *service) handleRequest(client *libs.Client, bytes []byte) {
	// ignore bad requests
	var req structure.Request
	if err := json.Unmarshal(bytes, &req); err != nil {
		return
	}

	if err := s.validate.Struct(req); err != nil {
		return
	}

	if req.Method == structure.RequestJoin && client.Name == "" {
		s.hub.Join <- &libs.JoinPack{Client: client, Name: *req.Name, Nonce: req.Nonce}
	} else if req.Method == structure.RequestMessage && client.Name != "" {
		s.handleMessage(client, &req)
	} else {
		client.Response(structure.MethodError, req.Nonce)
	}
}

func (s *service) handleMessage(client *libs.Client, req *structure.Request) {
	has, err := s.information.HasRoom(req.Message.Room)
	if err != nil || !has {
		client.Response(structure.RoomError, req.Nonce)
		return
	}

	date_time := time.Now()
	content := &structure.Content{From: client.Name, Message: req.Message, DateTime: &date_time}
	notification := &structure.Notification{Type: structure.NotificationMessage, Content: content, Nonce: req.Nonce}

	message, err := json.Marshal(notification)
	if err != nil {
		log.Println("SERVER ERROR: json.Marshal")
		return
	}

	encoded_data := base64.StdEncoding.EncodeToString(message)
	if err := s.redis_client.Publish(context.Background(), "broadcast", encoded_data).Err(); err != nil {
		log.Println("SERVER ERROR: redis.Publish")
		return
	}

	client.Response(structure.Success, notification.Nonce)
}
