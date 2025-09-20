package postgres

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/weeweeshka/hot-coffee/internal/models"
	"log/slog"
	"time"
)

type Storage struct {
	db *pgx.Conn
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

	return &Storage{db: conn}, nil
}
func (s *Storage) SaveInventory(ctx context.Context, data models.InventoryItem) (int64, error) {

}
