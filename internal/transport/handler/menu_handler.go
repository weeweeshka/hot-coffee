package handler

import (
	"context"
	"errors"
	models2 "github.com/weeweeshka/hot-coffee/internal/models"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MenuHandler struct {
	bus  MenuBus
	logr *slog.Logger
}

type MenuBus interface {
	CreateMenu(ctx context.Context, menu models2.MenuItem) (int64, error)
	GetMenus(ctx context.Context) ([]models2.MenuResponse, error)
	GetMenu(ctx context.Context, id int64) (models2.MenuResponse, error)
	UpdateMenu(ctx context.Context, id int64, menu models2.MenuItem) (models2.MenuItem, error)
	DeleteMenu(ctx context.Context, id int64) error
}

func (h *MenuHandler) CreateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		var menu models2.MenuItem
		if err := c.ShouldBindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id, err := h.bus.CreateMenu(c.Request.Context(), menu)
		if err != nil {
			writeError(c, http.StatusInternalServerError, err, h.logr, "CreateMenu: business error")
			return
		}

		h.logr.Info("Menu created", "id", id)
		c.JSON(http.StatusCreated, gin.H{"id": id, "status": "created"})
	}
}

func (h *MenuHandler) GetMenus() gin.HandlerFunc {
	return func(c *gin.Context) {
		menus, err := h.bus.GetMenus(c.Request.Context())
		if err != nil {
			writeError(c, http.StatusInternalServerError, err, h.logr, "GetMenus: business error")
			return
		}

		h.logr.Info("Menus retrieved")
		c.JSON(http.StatusOK, gin.H{"menus": menus})
	}
}

func (h *MenuHandler) GetMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		idstr := c.Param("id")
		id, err := strconv.Atoi(idstr)
		if err != nil {
			writeError(c, http.StatusInternalServerError, err, h.logr, "GetInventory: invalid id")
		}

		menu, err := h.bus.GetMenu(c.Request.Context(), int64(id))
		if err != nil {
			if errors.Is(err, models2.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "menu not found"})
				return
			}
			writeError(c, http.StatusInternalServerError, err, h.logr, "GetMenu: business error")
			return
		}

		h.logr.Info("Menu retrieved", "id", id)
		c.JSON(http.StatusOK, gin.H{"menu": menu})
	}
}

func (h *MenuHandler) UpdateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		idstr := c.Param("id")
		id, err := strconv.Atoi(idstr)
		if err != nil {
			writeError(c, http.StatusInternalServerError, err, h.logr, "GetInventory: invalid id")
		}

		var menu models2.MenuItem
		if err := c.ShouldBindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updated, err := h.bus.UpdateMenu(c.Request.Context(), int64(id), menu)
		if err != nil {
			if errors.Is(err, models2.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "menu not found"})
				return
			}
			writeError(c, http.StatusInternalServerError, err, h.logr, "UpdateMenu: business error")
			return
		}

		h.logr.Info("Menu updated", "id", id)
		c.JSON(http.StatusOK, gin.H{"id": id, "menu": updated})
	}
}

func (h *MenuHandler) DeleteMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		idstr := c.Param("id")
		id, err := strconv.Atoi(idstr)
		if err != nil {
			writeError(c, http.StatusInternalServerError, err, h.logr, "GetInventory: invalid id")
		}

		if err := h.bus.DeleteMenu(c.Request.Context(), int64(id)); err != nil {
			if errors.Is(err, models2.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "menu not found"})
				return
			}
			writeError(c, http.StatusInternalServerError, err, h.logr, "DeleteMenu: business error")
			return
		}

		h.logr.Info("Menu deleted", "id", id)
		c.JSON(http.StatusNoContent, gin.H{})
	}
}
