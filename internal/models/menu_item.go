package models

type MenuItem struct {
	ID          int64                `json:"product_id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Price       float64              `json:"price"`
	Ingredients []MenuItemIngredient `json:"ingredients"`
}

type MenuItemIngredient struct {
	IngredientID int64   `json:"ingredient_id"`
	Quantity     float64 `json:"quantity"`
}
