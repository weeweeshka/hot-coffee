package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/weeweeshka/hot-coffee/internal/models"
	"log/slog"
	"time"
)

type Storage struct {
	db   *pgx.Conn
	logr *slog.Logger
}

func NewStorage(logr *slog.Logger, connString string) (*Storage, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		panic(err)
	}
	logr.Info("connected to postgresql")

	return &Storage{
		logr: logr,
		db:   conn}, nil
}

func (s *Storage) SaveOrder(ctx context.Context, data models.OrderRequest) (int64, error) {

	tx, err := o.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("cannot begin transaction: %v", err)
	}

	defer tx.Rollback(ctx)

	_, err := s.db.Exec(ctx, "")
	if err != nil {
		return 0, fmt.Errorf("cannot insert order: %v", err)
	}

}
