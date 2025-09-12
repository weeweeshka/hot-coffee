package service

import (
	"context"
	"github.com/weeweeshka/hot-coffee/models"
	"log/slog"
)

type OrderImpl struct {
	logr *slog.Logger
	repo OrderRepo
}

type OrderRepo interface {
	Save(ctx context.Context, data models.Order) (string, error)
	GetAll(ctx context.Context) ([]models.Order, error)
	Get(ctx context.Context, id string) (models.Order, error)
	Update(ctx context.Context, id string, order models.Order) (models.Order, error)
	Delete(ctx context.Context, id string) error
}

func NewOrderService(repo OrderRepo, logr *slog.Logger) *OrderImpl {
	return &OrderImpl{repo: repo, logr: logr}
}

func (o *OrderImpl) CreateOrder(ctx context.Context, data models.Order) (string, error) {

	id, err := o.repo.Save(ctx, data)
	if err != nil {
		o.logr.Info("Failed to save order", "err", err)
	}
	return id, err
}

func (o *OrderImpl) GetOrders(ctx context.Context) ([]models.Order, error) {

	var orders []models.Order
	orders, err := o.repo.GetAll(ctx)
	if err != nil {
		o.logr.Info("Failed to get all orders", "err", err)
	}
	return orders, err
}

func (o *OrderImpl) GetOrder(ctx context.Context, id string) (models.Order, error) {
	var order models.Order
	order, err := o.repo.Get(ctx, id)
	if err != nil {
		o.logr.Info("Failed to get order", "err", err)
	}
	return order, err
}

func (o *OrderImpl) UpdateOrder(ctx context.Context, id string, order models.Order) (models.Order, error) {

	nOrder, err := o.repo.Update(ctx, id, order)
	if err != nil {
		o.logr.Info("Failed to update order", "err", err)
	}
	return nOrder, err
}

func (o *OrderImpl) DeleteOrder(ctx context.Context, id string) error {

	err := o.repo.Delete(ctx, id)
	if err != nil {
		o.logr.Info("Failed to delete order", "err", err)

	}
	return err
}
