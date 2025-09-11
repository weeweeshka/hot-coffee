package main

import (
	"github.com/gin-gonic/gin"
	handlers "github.com/weeweeshka/hot-coffee/internal/transport/handler"
)

func main() {
	router := gin.Default()
	groupOrder := router.Group("/orders")
	{
		groupOrder.POST("/", handlers.CreateOrder)
		groupOrder.GET("/")
		groupOrder.GET("/{id}")
		groupOrder.PUT("/{id}")
		groupOrder.DELETE("/{id}")
		groupOrder.POST("/{id}/close")
	}
}
