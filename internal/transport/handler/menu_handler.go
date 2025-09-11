package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weeweeshka/hot-coffee/models"
)

type MenuBus interface {
	CreateMenu(ctx context.Context, menu models.MenuItem) (string, error)
	GetMenus(ctx context.Context) ([]models.MenuItem, error)
	GetMenu(ctx context.Context, id string) (models.MenuItem, error)
	UpdateMenu(ctx context.Context, id string, menu models.MenuItem) (models.MenuItem, error)
	DeleteMenu(ctx context.Context, id string) error
}

func CreateMenu(logger *slog.Logger, bus MenuBus) gin.HandlerFunc {
	return func(c *gin.Context) {
		var menu models.MenuItem
		if err := c.ShouldBindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id, err := bus.CreateMenu(c.Request.Context(), menu)
		if err != nil {
			writeError(c, http.StatusInternalServerError, err, logger, "CreateMenu: business error")
			return
		}

		logger.Info("Menu created", "id", id)
		c.JSON(http.StatusCreated, gin.H{"id": id, "status": "created"})
	}
}

func GetMenus(logger *slog.Logger, bus MenuBus) gin.HandlerFunc {
	return func(c *gin.Context) {
		menus, err := bus.GetMenus(c.Request.Context())
		if err != nil {
			writeError(c, http.StatusInternalServerError, err, logger, "GetMenus: business error")
			return
		}

		logger.Info("Menus retrieved", "count", len(menus))
		c.JSON(http.StatusOK, gin.H{"menus": menus})
	}
}

func GetMenu(logger *slog.Logger, bus MenuBus) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
			return
		}

		menu, err := bus.GetMenu(c.Request.Context(), id)
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "menu not found"})
				return
			}
			writeError(c, http.StatusInternalServerError, err, logger, "GetMenu: business error")
			return
		}

		logger.Info("Menu retrieved", "id", id)
		c.JSON(http.StatusOK, gin.H{"id": id, "menu": menu})
	}
}

func UpdateMenu(logger *slog.Logger, bus MenuBus) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
			return
		}

		var menu models.MenuItem
		if err := c.ShouldBindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updated, err := bus.UpdateMenu(c.Request.Context(), id, menu)
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "menu not found"})
				return
			}
			writeError(c, http.StatusInternalServerError, err, logger, "UpdateMenu: business error")
			return
		}

		logger.Info("Menu updated", "id", id)
		c.JSON(http.StatusOK, gin.H{"id": id, "menu": updated})
	}
}

func DeleteMenu(logger *slog.Logger, bus MenuBus) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
			return
		}

		if err := bus.DeleteMenu(c.Request.Context(), id); err != nil {
			if errors.Is(err, models.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "menu not found"})
				return
			}
			writeError(c, http.StatusInternalServerError, err, logger, "DeleteMenu: business error")
			return
		}

		logger.Info("Menu deleted", "id", id)
		c.JSON(http.StatusNoContent, gin.H{})
	}
}
