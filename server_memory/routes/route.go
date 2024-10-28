package routes

import (
	"server_memory/libs"
	"server_memory/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rcrowley/go-metrics"
)

func Setup(g *gin.Engine, information *libs.Information, epoll *libs.Epoll, hub *libs.Hub, validate *validator.Validate, meter metrics.Meter) {
	service := service.New(information, epoll, hub, validate, meter)
	go service.StartReadWebsocket()

	g.GET("/ws", service.WebsocketHandler)
	g.Static("/assets", "./web/assets")
	g.StaticFile("/", "./web/index.html")
}
