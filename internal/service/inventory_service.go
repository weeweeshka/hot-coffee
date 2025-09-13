package service

import (
	"context"
	"github.com/weeweeshka/hot-coffee/models"
	"log/slog"
)

type InventoryImpl struct {
	logr *slog.Logger
	repo InventoryRepo
}

type InventoryRepo interface {
	SaveInventory(ctx context.Context, data models.InventoryItem) (string, error)
	GetAllInventories(ctx context.Context) ([]models.InventoryItem, error)
	GetInventory(ctx context.Context, id string) (models.InventoryItem, error)
	UpdateInventory(ctx context.Context, id string, inventory models.InventoryItem) (models.InventoryItem, error)
	DeleteInventory(ctx context.Context, id string) error
}

func NewInventoryService(logr *slog.Logger, repo InventoryRepo) *InventoryImpl {
	return &InventoryImpl{
		logr: logr,
		repo: repo,
	}
}

func (s *InventoryImpl) CreateInventory(ctx context.Context, inventory models.InventoryItem) (string, error) {

	id, err := s.repo.SaveInventory(ctx, inventory)
	if err != nil {
		s.logr.Info("Error creating inventory", err)
		return "", err
	}
	return id, nil
}

func (s *InventoryImpl) GetAllInventories(ctx context.Context) ([]models.InventoryItem, error) {

	var inventories []models.InventoryItem

	inventories, err := s.repo.GetAllInventories(ctx)
	if err != nil {
		s.logr.Info("Error getting inventory", err)
		return []models.InventoryItem{}, err
	}

	return inventories, nil
}

func (s *InventoryImpl) GetInventory(ctx context.Context, id string) (models.InventoryItem, error) {

	inventory, err := s.repo.GetInventory(ctx, id)
	if err != nil {
		s.logr.Info("Error getting inventory", err)
		return models.InventoryItem{}, err
	}
	return inventory, nil
}

func (s *InventoryImpl) UpdateInventory(ctx context.Context, id string, inventory models.InventoryItem) (models.InventoryItem, error) {
	var nInventory models.InventoryItem
	nInventory, err := s.repo.UpdateInventory(ctx, id, inventory)
	if err != nil {
		s.logr.Info("Error updating inventory", err)
		return models.InventoryItem{}, err
	}
	return nInventory, nil

}
func (s *InventoryImpl) DeleteInventory(ctx context.Context, id string) error {
	err := s.repo.DeleteInventory(ctx, id)
	if err != nil {
		s.logr.Info("Error deleting inventory", err)
		return err
	}
	return nil
}
