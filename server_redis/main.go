package main

import (
	"log"
	_ "server_redis/config"
	"server_redis/libs"
	"server_redis/routes"
	"server_redis/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	validate := libs.NewValidate()
	utils.RegisterValidation(validate)

	redis_client := libs.NewRedis()
	information := libs.NewInformation(redis_client)

	go_pool := libs.NewGoPool()
	defer go_pool.Release()

	epoll := libs.NewEpoll()
	defer epoll.Close()

	hub := libs.NewHub(information, epoll, go_pool, redis_client)
	go hub.Run()
	defer hub.Close()

	counter := libs.NewCounter(10000)
	go counter.Start()
	defer counter.Close()

	subscriber := libs.NewSubscriber(redis_client, hub)
	subscriber.Sub()

	g := gin.Default()
	routes.Setup(g, redis_client, information, epoll, hub, validate, counter)

	log.Printf("start server ws://0.0.0.0:8000")
	log.Fatal(g.Run(":8000"))
}
