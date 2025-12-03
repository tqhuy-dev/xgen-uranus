package main

import (
	"github.com/gin-gonic/gin"
	"github.com/tqhuy-dev/xgen-uranus/examples/pb/greet"
	"github.com/tqhuy-dev/xgen-uranus/transport/grpc"
)

func Register(server *grpc.Server) {
	var greetServer GreetServer
	greet.RegisterGreetServer(server, &greetServer)
}

func RegisterRouter(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
