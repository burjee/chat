package main

import (
	"log"
	"os"
	_ "server_memory/config"
	"server_memory/libs"
	"server_memory/routes"
	"server_memory/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rcrowley/go-metrics"
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

	meter := metrics.NewMeter()
	metrics.Register("meter", meter)
	go metrics.Log(metrics.DefaultRegistry, time.Second, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))

	g := gin.Default()
	routes.Setup(g, information, epoll, hub, validate, meter)

	log.Printf("start server ws://0.0.0.0:8000")
	log.Fatal(g.Run(":8000"))
}
