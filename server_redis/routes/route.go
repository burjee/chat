package routes

import (
	"server_redis/libs"
	"server_redis/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rcrowley/go-metrics"
	"github.com/redis/go-redis/v9"
)

func Setup(g *gin.Engine, redis_client *redis.Client, information *libs.Information, epoll *libs.Epoll, hub *libs.Hub, validate *validator.Validate, meter metrics.Meter) {
	service := service.New(redis_client, information, epoll, hub, validate, meter)
	go service.StartReadWebsocket()

	g.GET("/ws", service.WebsocketHandler)
	g.Static("/assets", "./web/assets")
	g.StaticFile("/", "./web/index.html")
}
