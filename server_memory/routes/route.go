package routes

import (
	"server_memory/libs"
	"server_memory/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Setup(g *gin.Engine, information *libs.Information, epoll *libs.Epoll, hub *libs.Hub, validate *validator.Validate, counter *libs.Counter) {
	service := service.New(information, epoll, hub, validate, counter)
	go service.StartReadWebsocket()

	g.GET("/ws", service.WebsocketHandler)
	g.Static("/assets", "./web/assets")
	g.StaticFile("/", "./web/index.html")
}
