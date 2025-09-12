package service

import (
	"context"
	"github.com/weeweeshka/hot-coffee/models"
	"log/slog"
)

type MenuImpl struct {
	logr *slog.Logger
	repo MenuRepo
}

type MenuRepo interface {
	Save(ctx context.Context, data models.MenuItem) (string, error)
	GetAll(ctx context.Context) ([]models.MenuItem, error)
	Get(ctx context.Context, id string) (models.MenuItem, error)
	Update(ctx context.Context, id string, order models.MenuItem) (models.MenuItem, error)
	Delete(ctx context.Context, id string) error
}

func NewMenuService(logr *slog.Logger, repo MenuRepo) *MenuImpl {
	return &MenuImpl{
		logr: logr,
		repo: repo,
	}
}

func (m *MenuImpl) CreateMenu(ctx context.Context, menu models.MenuItem) (string, error) {

	id, err := m.repo.Save(ctx, menu)
	if err != nil {
		m.logr.Info("Menu Create Error", "err", err)
		return "", err
	}

	return id, nil
}

func (m *MenuImpl) GetMenus(ctx context.Context) ([]models.MenuItem, error) {

	var menus []models.MenuItem

	menus, err := m.repo.GetAll(ctx)
	if err != nil {
		m.logr.Info("Menu Get All Error", "err", err)
		return []models.MenuItem{}, err
	}

	return menus, nil
}

func (m *MenuImpl) GetMenu(ctx context.Context, id string) (models.MenuItem, error) {

	var menu models.MenuItem
	menu, err := m.repo.Get(ctx, id)
	if err != nil {
		m.logr.Info("Menu Get Error", "err", err)
		return models.MenuItem{}, err
	}
	return menu, nil
}

func (m *MenuImpl) UpdateMenu(ctx context.Context, id string, menu models.MenuItem) (models.MenuItem, error) {

	var nMenu models.MenuItem
	nMenu, err := m.repo.Update(ctx, id, menu)
	if err != nil {
		m.logr.Info("Menu Update Error", "err", err)
	}
	return nMenu, err
}

func (m *MenuImpl) DeleteMenu(ctx context.Context, id string) error {
	err := m.repo.Delete(ctx, id)
	if err != nil {
		m.logr.Info("Menu Delete Error", "err", err)
	}
	return err
}
