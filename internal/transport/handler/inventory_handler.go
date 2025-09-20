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

type InventoryHandler struct {
	bus  InventoryBus
	logr *slog.Logger
}

type InventoryBus interface {
	CreateInventory(ctx context.Context, inventory models.InventoryItem) (int64, error)
	GetInventories(ctx context.Context) ([]models.InventoryItem, error)
	GetInventory(ctx context.Context, id int64) (models.InventoryItem, error)
	UpdateInventory(ctx context.Context, id int64, inventory models.InventoryItem) (models.InventoryItem, error)
	DeleteInventory(ctx context.Context, id int64) error
}

func (h *InventoryHandler) CreateInventory() gin.HandlerFunc {
	return func(c *gin.Context) {
		var inventory models.InventoryItem
		if err := c.ShouldBindJSON(&inventory); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id, err := h.bus.CreateInventory(c.Request.Context(), inventory)
		if err != nil {
			writeError(c, http.StatusInternalServerError, err, h.logr, "CreateInventory: business error")
			return
		}

		h.logr.Info("Inventory created", "id", id)
		c.JSON(http.StatusCreated, gin.H{"id": id})
	}
}

func (h *InventoryHandler) GetInventories() gin.HandlerFunc {
	return func(c *gin.Context) {
		inventories, err := h.bus.GetInventories(c.Request.Context())
		if err != nil {
			writeError(c, http.StatusInternalServerError, err, h.logr, "GetInventories: business error")
			return
		}

		h.logr.Info("Inventories retrieved")
		c.JSON(http.StatusOK, gin.H{"inventories": inventories})
	}
}

func (h *InventoryHandler) GetInventory() gin.HandlerFunc {
	return func(c *gin.Context) {
		idstr := c.Param("id")
		id, err := strconv.Atoi(idstr)
		if err != nil {
			writeError(c, http.StatusInternalServerError, err, h.logr, "GetInventory: invalid id")
		}

		if id == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
			return
		}

		inv, err := h.bus.GetInventory(c.Request.Context(), int64(id))
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "inventory not found"})
				return
			}
			writeError(c, http.StatusInternalServerError, err, h.logr, "GetInventory: business error")
			return
		}

		h.logr.Info("Inventory retrieved", "id", id)
		c.JSON(http.StatusOK, gin.H{"id": id, "inventory": inv})
	}
}

func (h *InventoryHandler) UpdateInventory() gin.HandlerFunc {
	return func(c *gin.Context) {
		idstr := c.Param("id")
		id, err := strconv.Atoi(idstr)
		if err != nil {
			writeError(c, http.StatusInternalServerError, err, h.logr, "GetInventory: invalid id")
		}

		if id == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
			return
		}

		var inventory models.InventoryItem
		if err := c.ShouldBindJSON(&inventory); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updated, err := h.bus.UpdateInventory(c.Request.Context(), int64(id), inventory)
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "inventory not found"})
				return
			}
			writeError(c, http.StatusInternalServerError, err, h.logr, "UpdateInventory: business error")
			return
		}

		h.logr.Info("Inventory updated", "id", id)
		c.JSON(http.StatusOK, gin.H{"id": id, "inventory": updated})
	}
}

func (h *InventoryHandler) DeleteInventory() gin.HandlerFunc {
	return func(c *gin.Context) {
		idstr := c.Param("id")
		id, err := strconv.Atoi(idstr)
		if err != nil {
			writeError(c, http.StatusInternalServerError, err, h.logr, "GetInventory: invalid id")
		}

		if err := h.bus.DeleteInventory(c.Request.Context(), int64(id)); err != nil {
			if errors.Is(err, models.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "inventory not found"})
				return
			}
			writeError(c, http.StatusInternalServerError, err, h.logr, "DeleteInventory: business error")
			return
		}

		h.logr.Info("Inventory deleted", "id", id)
		c.JSON(http.StatusNoContent, gin.H{})
	}
}
