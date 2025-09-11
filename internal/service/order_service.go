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
}

func NewOrderService(repo OrderRepo, logr *slog.Logger) *OrderImpl {
	return &OrderImpl{repo: repo, logr: logr}
}

func (o *OrderImpl) CreateOrder(ctx context.Context, data models.Order) (string, error) {
}
