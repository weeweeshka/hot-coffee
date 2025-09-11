package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weeweeshka/hot-coffee/models"
)

type InventoryBus interface {
	CreateInventory(ctx context.Context, inventory models.InventoryItem) (string, error)
	GetInventories(ctx context.Context) ([]models.InventoryItem, error)
	GetInventory(ctx context.Context, id string) (models.InventoryItem, error)
	UpdateInventory(ctx context.Context, id string, inventory models.InventoryItem) (models.InventoryItem, error)
	DeleteInventory(ctx context.Context, id string) error
}

func CreateInventory(logger *slog.Logger, bus InventoryBus) gin.HandlerFunc {
	return func(c *gin.Context) {
		var inventory models.InventoryItem
		if err := c.ShouldBindJSON(&inventory); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id, err := bus.CreateInventory(c.Request.Context(), inventory)
		if err != nil {
			writeError(c, http.StatusInternalServerError, err, logger, "CreateInventory: business error")
			return
		}

		logger.Info("Inventory created", "id", id)
		c.JSON(http.StatusCreated, gin.H{"id": id})
	}
}

func GetInventories(logger *slog.Logger, bus InventoryBus) gin.HandlerFunc {
	return func(c *gin.Context) {
		inventories, err := bus.GetInventories(c.Request.Context())
		if err != nil {
			writeError(c, http.StatusInternalServerError, err, logger, "GetInventories: business error")
			return
		}

		logger.Info("Inventories retrieved", "count", len(inventories))
		c.JSON(http.StatusOK, gin.H{"inventories": inventories})
	}
}

func GetInventory(logger *slog.Logger, bus InventoryBus) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
			return
		}

		inv, err := bus.GetInventory(c.Request.Context(), id)
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "inventory not found"})
				return
			}
			writeError(c, http.StatusInternalServerError, err, logger, "GetInventory: business error")
			return
		}

		logger.Info("Inventory retrieved", "id", id)
		c.JSON(http.StatusOK, gin.H{"id": id, "inventory": inv})
	}
}

func UpdateInventory(logger *slog.Logger, bus InventoryBus) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
			return
		}

		var inventory models.InventoryItem
		if err := c.ShouldBindJSON(&inventory); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updated, err := bus.UpdateInventory(c.Request.Context(), id, inventory)
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "inventory not found"})
				return
			}
			writeError(c, http.StatusInternalServerError, err, logger, "UpdateInventory: business error")
			return
		}

		logger.Info("Inventory updated", "id", id)
		c.JSON(http.StatusOK, gin.H{"id": id, "inventory": updated})
	}
}

func DeleteInventory(logger *slog.Logger, bus InventoryBus) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
			return
		}

		if err := bus.DeleteInventory(c.Request.Context(), id); err != nil {
			if errors.Is(err, models.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "inventory not found"})
				return
			}
			writeError(c, http.StatusInternalServerError, err, logger, "DeleteInventory: business error")
			return
		}

		logger.Info("Inventory deleted", "id", id)
		c.JSON(http.StatusNoContent, gin.H{})
	}
}
