package Order

import (
	"SEProject/Order/types"
	"context"
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Sepete ürün ekle
func (r *Repository) AddToCart(ctx context.Context, item *types.CartItem) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO cart (user_id, product_id, quantity)
		VALUES ($1, $2, $3)
	`, item.UserID, item.ProductID, item.Quantity)
	return err
}

// Sepeti getir
func (r *Repository) GetCart(ctx context.Context, userID int) ([]types.CartItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT user_id, product_id, quantity FROM cart
		WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cart []types.CartItem
	for rows.Next() {
		var item types.CartItem
		if err := rows.Scan(&item.UserID, &item.ProductID, &item.Quantity); err != nil {
			return nil, err
		}
		cart = append(cart, item)
	}
	return cart, nil
}

// Sipariş oluştur
func (r *Repository) PlaceOrder(ctx context.Context, req *types.PlaceOrderRequest) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO orders (user_id, restaurant_id, address, status)
		VALUES ($1, $2, $3, 'pending')
	`, req.UserID, req.RestaurantID, req.Address)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, `
		DELETE FROM cart WHERE user_id = $1
	`, req.UserID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// Restoranın siparişlerini getir
func (r *Repository) GetOrdersByRestaurantID(ctx context.Context, restaurantID int) ([]types.Order, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, restaurant_id, address, status
		FROM orders
		WHERE restaurant_id = $1
	`, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []types.Order
	for rows.Next() {
		var o types.Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.RestaurantID, &o.Address, &o.Status); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

// Siparişin durumunu güncelle (prepared/sent)
func (r *Repository) UpdateOrderStatus(ctx context.Context, orderID int, status string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE orders SET status = $1 WHERE id = $2
	`, status, orderID)
	return err
}
