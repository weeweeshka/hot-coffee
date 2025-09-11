package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/weeweeshka/hot-coffee/models"
	"log/slog"
	"net/http"
)

type OrderBus interface {
	CreateOrder(data models.Order) (string, error)
	GetOrders() ([]models.Order, error)
	GetOrder(id string) (models.Order, error)
	PutOrder(id string, order models.Order) (string, error)
	DeleteOrder(id string) error
	CloseOrder(id string) error
}

func CreateOrder(logr *slog.Logger, bus OrderBus) gin.HandlerFunc {
	return func(c *gin.Context) {
		var order models.Order

		if err := c.ShouldBindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id, err := bus.CreateOrder(order)
		if err != nil {
			logr.Error("Business error", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"id":     id,
			"status": "created",
		})

		logr.Info("Order created", "order_id", id)
	}
}

func GetOrders(logr *slog.Logger, bus OrderBus) gin.HandlerFunc {
	return func(c *gin.Context) {

		var orders []models.Order

		orders, err := bus.GetOrders()
		if err != nil {
			logr.Error("Business error", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		logr.Info("Successfully retrieved orders")
		c.JSON(http.StatusOK, gin.H{
			"orders": orders,
		})
	}
}

func GetOrder(id string, bus OrderBus, logr *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		order, err := bus.GetOrder(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		logr.Info("Successfully retrieved order", "id", id)
		c.JSON(http.StatusOK, gin.H{
			"order": order,
		})

	}
}

func PutOrder(id string, bus OrderBus, logr *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var order models.Order
		if err := c.ShouldBindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		_, err := bus.PutOrder(id, order)
		if err != nil {
			logr.Error("Business error", "error", err)
		}
		logr.Info("Successfully put order", "id", id)

		c.JSON(http.StatusOK, gin.H{
			"id":    id,
			"order": order,
		})

	}
}

func DeleteOrder(id string, bus OrderBus, logr *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := bus.DeleteOrder(id)
		if err != nil {
			logr.Error("Business error", "error", err)
		}
		logr.Info("Successfully deleted order", "id", id)
		c.JSON(http.StatusOK, gin.H{})
	}
}

func CloseOrder(id string, bus OrderBus, logr *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		err := bus.CloseOrder(id)
		if err != nil {
			logr.Error("Business error", "error", err)
		}
		logr.Info("Successfully closed order", "id", id)
	}
}
