package main

import (
	"log"
	_ "server_memory/config"
	"server_memory/libs"
	"server_memory/routes"
	"server_memory/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	validate := libs.NewValidate()
	utils.RegisterValidation(validate)

	information := libs.NewInformation()

	go_pool := libs.NewGoPool()
	defer go_pool.Release()

	epoll := libs.NewEpoll()
	defer epoll.Close()

	hub := libs.NewHub(information, epoll, go_pool)
	go hub.Run()
	defer hub.Close()

	counter := libs.NewCounter(10000)
	go counter.Start()
	defer counter.Close()

	g := gin.Default()
	routes.Setup(g, information, epoll, hub, validate, counter)

	log.Printf("start server ws://0.0.0.0:8000")
	log.Fatal(g.Run(":8000"))
}
