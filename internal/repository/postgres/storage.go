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

func (s *Storage) SaveInventory(ctx context.Context, data models.InventoryItem) (int64, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("cannot begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	var inventoryId int64
	err = tx.QueryRow(ctx, `INSERT INTO inventory(name, quantity, unit) VALUES ($1, $2, $3) RETURNING id`, data.Name, data.Quantity, data.Unit).Scan(&inventoryId)
	if err != nil {
		return 0, fmt.Errorf("cannot insert into inventory: %v", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return 0, fmt.Errorf("cannot commit transaction: %v", err)
	}
	s.logr.Info("commit transaction")
	return inventoryId, nil
}

func (s *Storage) SaveMenu(ctx context.Context, data models.MenuItem) (int64, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("cannot begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var menuID int64
	err = tx.QueryRow(ctx, `
        INSERT INTO menus(name, description, price) 
        VALUES ($1, $2, $3) 
        RETURNING id
    `, data.Name, data.Description, data.Price).Scan(&menuID)

	if err != nil {
		return 0, fmt.Errorf("cannot save menu: %w", err)
	}

	for _, ingredient := range data.Ingredients {
		_, err = tx.Exec(ctx, `INSERT INTO menu_ingredients(ingredient_id, quantity) VALUES ($1, $2, $3)`, ingredient.IngredientID, ingredient.Quantity)
		if err != nil {
			return 0, fmt.Errorf("cannot save menu_ingredients: %w", err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("cannot commit transaction: %w", err)
	}
	return menuID, nil
}

func (s *Storage) SaveOrder(ctx context.Context, data models.OrderRequest) (int64, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("cannot begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var orderId int64
	err = s.db.QueryRow(ctx, `INSERT INTO orders(
                   customer_name) VALUES ($1) RETURNING id`, data.CustomerName).Scan(&orderId)
	if err != nil {
		return 0, fmt.Errorf("cannot insert into orders: %v", err)
	}

	for _, items := range data.Items {
		_, err = s.db.Exec(ctx, `INSERT INTO items(
                  order_id, menu_id ,quantity) VALUES ($1, $2, $3)`, orderId, items.ProductID, items.Quantity)
		if err != nil {
			return 0, fmt.Errorf("cannot insert into items: %v", err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("cannot commit transaction: %w", err)
	}
	s.logr.Info("order saved")

	return orderId, nil
}
