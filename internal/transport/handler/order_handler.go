package handler

import (
	"context"
	"errors"
	"github.com/weeweeshka/hot-coffee/internal/models"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	bus  OrderBus
	logr *slog.Logger
}

type OrderBus interface {
	CreateOrder(ctx context.Context, data models.OrderRequest) (int64, error)
	GetOrders(ctx context.Context) ([]models.Order, error)
	GetOrder(ctx context.Context, id int64) (models.Order, error)
	UpdateOrder(ctx context.Context, id int64, order models.OrderRequest) (models.Order, error)
	DeleteOrder(ctx context.Context, id int64) error
	CloseOrder(ctx context.Context, id int64) error
}

func writeError(c *gin.Context, code int, err error, logger *slog.Logger, msg string) {
	logger.Error(msg, "error", err)
	c.JSON(code, gin.H{"error": err.Error()})
}

func (h *OrderHandler) CreateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var orderReq models.OrderRequest
		if err := c.ShouldBindJSON(&orderReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id, err := h.bus.CreateOrder(c.Request.Context(), orderReq)
		if err != nil {
			writeError(c, http.StatusInternalServerError, err, h.logr, "CreateOrder: business error")
			return
		}

		h.logr.Info("Order created", "order_id", id)
		c.JSON(http.StatusCreated, gin.H{"id": id, "status": "created"})
	}
}

func (h *OrderHandler) GetOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		orders, err := h.bus.GetOrders(c.Request.Context())
		if err != nil {
			writeError(c, http.StatusInternalServerError, err, h.logr, "GetOrders: business error")
			return
		}

		h.logr.Info("Orders retrieved", "count", len(orders))
		c.JSON(http.StatusOK, gin.H{"orders": orders})
	}
}

func (h *OrderHandler) GetOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		idstr := c.Param("id")
		id, err := strconv.Atoi(idstr)
		if err != nil {
			writeError(c, http.StatusInternalServerError, err, h.logr, "GetInventory: invalid id")
		}

		order, err := h.bus.GetOrder(c.Request.Context(), int64(id))
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
				return
			}
			writeError(c, http.StatusInternalServerError, err, h.logr, "GetOrder: business error")
			return
		}

		h.logr.Info("Order retrieved", "id", id)
		c.JSON(http.StatusOK, gin.H{"order": order})
	}
}

func (h *OrderHandler) UpdateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		idstr := c.Param("id")
		id, err := strconv.Atoi(idstr)
		if err != nil {
			writeError(c, http.StatusInternalServerError, err, h.logr, "GetInventory: invalid id")
		}

		var order models.OrderRequest
		if err := c.ShouldBindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		nOrder, err := h.bus.UpdateOrder(c.Request.Context(), int64(id), order)
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
				return
			}
			writeError(c, http.StatusInternalServerError, err, h.logr, "UpdateOrder: business error")
			return
		}

		h.logr.Info("Order updated", "id", nOrder.ID)
		c.JSON(http.StatusOK, gin.H{"id": nOrder.ID, "order": nOrder})
	}
}

func (h *OrderHandler) DeleteOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		idstr := c.Param("id")
		id, err := strconv.Atoi(idstr)
		if err != nil {
			writeError(c, http.StatusInternalServerError, err, h.logr, "GetInventory: invalid id")
		}

		if err := h.bus.DeleteOrder(c.Request.Context(), int64(id)); err != nil {
			if errors.Is(err, models.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
				return
			}
			writeError(c, http.StatusInternalServerError, err, h.logr, "DeleteOrder: business error")
			return
		}

		h.logr.Info("Order deleted", "id", id)
		c.JSON(http.StatusNoContent, gin.H{})
	}
}

func (h *OrderHandler) CloseOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		idstr := c.Param("id")
		id, err := strconv.Atoi(idstr)
		if err != nil {
			writeError(c, http.StatusInternalServerError, err, h.logr, "GetInventory: invalid id")
		}

		if err := h.bus.CloseOrder(c.Request.Context(), int64(id)); err != nil {
			if errors.Is(err, models.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
				return
			}
			writeError(c, http.StatusInternalServerError, err, h.logr, "CloseOrder: business error")
			return
		}

		h.logr.Info("Order closed", "id", id)
		c.JSON(http.StatusOK, gin.H{"id": id, "status": "closed"})
	}
}
