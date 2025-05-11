package Cart

import (
	"SEProject/Cart/types"
	"context"
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) AddItem(ctx context.Context, item *types.CartItem) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO cart (user_id, product_id, quantity)
		VALUES ($1, $2, $3)
	`, item.UserID, item.ProductID, item.Quantity)
	return err
}

func (r *Repository) RemoveItem(ctx context.Context, userID, productID int) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM cart
		WHERE user_id = $1 AND product_id = $2
	`, userID, productID)
	return err
}

func (r *Repository) GetItems(ctx context.Context, userID int) ([]*types.CartItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, product_id, quantity
		FROM cart
		WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*types.CartItem
	for rows.Next() {
		var item types.CartItem
		if err := rows.Scan(&item.ID, &item.UserID, &item.ProductID, &item.Quantity); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, nil
}
