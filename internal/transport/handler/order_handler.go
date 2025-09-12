package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weeweeshka/hot-coffee/models"
)

type OrderBus interface {
	CreateOrder(ctx context.Context, data models.Order) (string, error)
	GetOrders(ctx context.Context) ([]models.Order, error)
	GetOrder(ctx context.Context, id string) (models.Order, error)
	UpdateOrder(ctx context.Context, id string, order models.Order) (models.Order, error)
	DeleteOrder(ctx context.Context, id string) error
	CloseOrder(ctx context.Context, id string) error
}

func writeError(c *gin.Context, code int, err error, logger *slog.Logger, msg string) {
	logger.Error(msg, "error", err)
	c.JSON(code, gin.H{"error": err.Error()})
}

func CreateOrder(logger *slog.Logger, bus OrderBus) gin.HandlerFunc {
	return func(c *gin.Context) {
		var order models.Order
		if err := c.ShouldBindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id, err := bus.CreateOrder(c.Request.Context(), order)
		if err != nil {
			writeError(c, http.StatusInternalServerError, err, logger, "CreateOrder: business error")
			return
		}

		logger.Info("Order created", "order_id", id)
		c.JSON(http.StatusCreated, gin.H{"id": id, "status": "created"})
	}
}

func GetOrders(logger *slog.Logger, bus OrderBus) gin.HandlerFunc {
	return func(c *gin.Context) {
		orders, err := bus.GetOrders(c.Request.Context())
		if err != nil {
			writeError(c, http.StatusInternalServerError, err, logger, "GetOrders: business error")
			return
		}

		logger.Info("Orders retrieved", "count", len(orders))
		c.JSON(http.StatusOK, gin.H{"orders": orders})
	}
}

func GetOrder(logger *slog.Logger, bus OrderBus) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
			return
		}

		order, err := bus.GetOrder(c.Request.Context(), id)
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
				return
			}
			writeError(c, http.StatusInternalServerError, err, logger, "GetOrder: business error")
			return
		}

		logger.Info("Order retrieved", "id", id)
		c.JSON(http.StatusOK, gin.H{"order": order})
	}
}

func UpdateOrder(logger *slog.Logger, bus OrderBus) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
			return
		}

		var order models.Order
		if err := c.ShouldBindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		nOrder, err := bus.UpdateOrder(c.Request.Context(), id, order)
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
				return
			}
			writeError(c, http.StatusInternalServerError, err, logger, "UpdateOrder: business error")
			return
		}

		logger.Info("Order updated", "id", nOrder.ID)
		c.JSON(http.StatusOK, gin.H{"id": nOrder.ID, "order": nOrder})
	}
}

func DeleteOrder(logger *slog.Logger, bus OrderBus) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
			return
		}

		if err := bus.DeleteOrder(c.Request.Context(), id); err != nil {
			if errors.Is(err, models.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
				return
			}
			writeError(c, http.StatusInternalServerError, err, logger, "DeleteOrder: business error")
			return
		}

		logger.Info("Order deleted", "id", id)
		c.JSON(http.StatusNoContent, gin.H{})
	}
}

func CloseOrder(logger *slog.Logger, bus OrderBus) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
			return
		}

		if err := bus.CloseOrder(c.Request.Context(), id); err != nil {
			if errors.Is(err, models.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
				return
			}
			writeError(c, http.StatusInternalServerError, err, logger, "CloseOrder: business error")
			return
		}

		logger.Info("Order closed", "id", id)
		c.JSON(http.StatusOK, gin.H{"id": id, "status": "closed"})
	}
}
