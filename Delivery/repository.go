package Delivery

import (
	"SEProject/Delivery/types"
	"context"
	"database/sql"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Siparişleri kuryeye göre getirir
func (r *Repository) GetOrdersByDeliveryPersonID(ctx context.Context, deliveryPersonID int) ([]types.DeliveryOrder, error) {
	query := `
		SELECT o.id, u.username, r.name, o.address, o.status
		FROM orders o
		JOIN users u ON o.user_id = u.id
		JOIN restaurants r ON o.restaurant_id = r.id
		WHERE o.delivery_person_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, deliveryPersonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []types.DeliveryOrder
	for rows.Next() {
		var order types.DeliveryOrder
		if err := rows.Scan(
			&order.OrderID,
			&order.CustomerName,
			&order.RestaurantName,
			&order.Address,
			&order.Status,
		); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// Sipariş durumunu günceller
func (r *Repository) UpdateOrderStatus(ctx context.Context, orderID int, status string) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE orders
		SET status = $1
		WHERE id = $2
	`, status, orderID)

	if err != nil {
		fmt.Println("❌ SQL HATASI:", err)
		return err
	}

	rows, _ := res.RowsAffected()
	fmt.Println("📝 Güncellenen satır sayısı:", rows)

	if rows == 0 {
		return fmt.Errorf("hiçbir kayıt güncellenmedi (id=%d olabilir mi?)", orderID)
	}

	return nil
}
