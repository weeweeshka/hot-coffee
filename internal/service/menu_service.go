package service

import (
	"context"
	"github.com/weeweeshka/hot-coffee/internal/models"
	"log/slog"
)

type MenuImpl struct {
	logr *slog.Logger
	repo MenuRepo
}

type MenuRepo interface {
	SaveMenu(ctx context.Context, data models.MenuItem) (int64, error)
	GetAllMenus(ctx context.Context) ([]models.MenuItem, error)
	GetMenu(ctx context.Context, id int64) (models.MenuItem, error)
	UpdateMenu(ctx context.Context, id int64, order models.MenuItem) (models.MenuItem, error)
	DeleteMenu(ctx context.Context, id int64) error
}

func NewMenuService(logr *slog.Logger, repo MenuRepo) *MenuImpl {
	return &MenuImpl{
		logr: logr,
		repo: repo,
	}
}

func (m *MenuImpl) CreateMenu(ctx context.Context, menu models.MenuItem) (int64, error) {

	id, err := m.repo.SaveMenu(ctx, menu)
	if err != nil {
		m.logr.Info("Menu Create Error", "err", err)
		return 0, err
	}

	return id, nil
}

func (m *MenuImpl) GetMenus(ctx context.Context) ([]models.MenuItem, error) {

	var menus []models.MenuItem

	menus, err := m.repo.GetAllMenus(ctx)
	if err != nil {
		m.logr.Info("Menu Get All Error", "err", err)
		return []models.MenuItem{}, err
	}

	return menus, nil
}

func (m *MenuImpl) GetMenu(ctx context.Context, id int64) (models.MenuItem, error) {

	var menu models.MenuItem
	menu, err := m.repo.GetMenu(ctx, id)
	if err != nil {
		m.logr.Info("Menu Get Error", "err", err)
		return models.MenuItem{}, err
	}
	return menu, nil
}

func (m *MenuImpl) UpdateMenu(ctx context.Context, id int64, menu models.MenuItem) (models.MenuItem, error) {

	var nMenu models.MenuItem
	nMenu, err := m.repo.UpdateMenu(ctx, id, menu)
	if err != nil {
		m.logr.Info("Menu Update Error", "err", err)
	}
	return nMenu, err
}

func (m *MenuImpl) DeleteMenu(ctx context.Context, id int64) error {
	err := m.repo.DeleteMenu(ctx, id)
	if err != nil {
		m.logr.Info("Menu Delete Error", "err", err)
	}
	return err
}
