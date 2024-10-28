package main

import (
	"log"
	"os"
	_ "server_redis/config"
	"server_redis/libs"
	"server_redis/routes"
	"server_redis/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rcrowley/go-metrics"
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

	meter := metrics.NewMeter()
	metrics.Register("ws-meter", meter)
	go metrics.Log(metrics.DefaultRegistry, time.Second, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))

	subscriber := libs.NewSubscriber(redis_client, hub)
	subscriber.Sub()

	g := gin.Default()
	routes.Setup(g, redis_client, information, epoll, hub, validate, meter)

	log.Printf("start server ws://0.0.0.0:8000")
	log.Fatal(g.Run(":8000"))
}
