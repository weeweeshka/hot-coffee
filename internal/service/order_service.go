package service

import (
	"context"
	"github.com/weeweeshka/hot-coffee/internal/models"
	"log/slog"
)

type OrderImpl struct {
	logr *slog.Logger
	repo OrderRepo
}

type OrderRepo interface {
	SaveOrder(ctx context.Context, data models.OrderRequest) (int64, error)
	GetAllOrders(ctx context.Context) ([]models.Order, error)
	GetOrder(ctx context.Context, id int64) (models.Order, error)
	UpdateOrder(ctx context.Context, id int64, order models.Order) (models.Order, error)
	DeleteOrder(ctx context.Context, id int64) error
}

func NewOrderService(repo OrderRepo, logr *slog.Logger) *OrderImpl {
	return &OrderImpl{repo: repo, logr: logr}
}

func (o *OrderImpl) CreateOrder(ctx context.Context, data models.OrderRequest) (int64, error) {

	id, err := o.repo.SaveOrder(ctx, data)
	if err != nil {
		o.logr.Info("Failed to save order", "err", err)
		return 0, err
	}

	return id, err
}

func (o *OrderImpl) GetOrders(ctx context.Context) ([]models.Order, error) {

	var orders []models.Order
	orders, err := o.repo.GetAllOrders(ctx)
	if err != nil {
		o.logr.Info("Failed to get all orders", "err", err)
		return []models.Order{}, err
	}
	return orders, nil
}

func (o *OrderImpl) GetOrder(ctx context.Context, id int64) (models.Order, error) {
	var order models.Order
	order, err := o.repo.GetOrder(ctx, id)
	if err != nil {
		o.logr.Info("Failed to get order", "err", err)
		return models.Order{}, err
	}
	return order, nil
}

func (o *OrderImpl) UpdateOrder(ctx context.Context, id int64, order models.Order) (models.Order, error) {

	nOrder, err := o.repo.UpdateOrder(ctx, id, order)
	if err != nil {
		o.logr.Info("Failed to update order", "err", err)
		return models.Order{}, err
	}
	return nOrder, nil
}

func (o *OrderImpl) DeleteOrder(ctx context.Context, id int64) error {

	err := o.repo.DeleteOrder(ctx, id)
	if err != nil {
		o.logr.Info("Failed to delete order", "err", err)
		return err
	}

	return nil
}

func (o *OrderImpl) CloseOrder(ctx context.Context, id int64) error {
	var order models.Order

	order, err := o.repo.GetOrder(ctx, id)
	if err != nil {
		o.logr.Info("Failed to get order", "err", err)
		return err
	}

	order.Status = "closed"
	_, err = o.repo.UpdateOrder(ctx, id, order)
	if err != nil {
		o.logr.Info("Failed to update order", "err", err)
		return err
	}

	return err
}
