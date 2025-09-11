package main

import "github.com/gin-gonic/gin"

func main() {
	router := gin.Default()
	groupOrder := router.Group("/orders")
	{
		groupOrder.POST("/", order_handler.CreateOrder)
		groupOrder.GET("/")
		groupOrder.GET("/{id}")
		groupOrder.PUT("/{id}")
		groupOrder.DELETE("/{id}")
		groupOrder.POST("/{id}/close")
	}
}
