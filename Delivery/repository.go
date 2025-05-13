package Delivery

import (
	"context"
	"database/sql"

	ordertypes "SEProject/Order/types"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetOrdersByDeliveryPersonID(ctx context.Context, deliveryPersonID int) ([]ordertypes.OrderResponse, error) {
	rows, err := r.db.QueryContext(ctx, `
    SELECT  o.id,
            u.username,
            r.name,
            o.address,
            o.status,
            TO_CHAR(o.total, 'FM999999999.00')
    FROM    orders o
    JOIN    users        u ON u.id = o.user_id
    JOIN    restaurants  r ON r.id = o.restaurant_id
    WHERE   o.status = 'prepared'
      AND   o.delivery_person_id @> ('[' || $1::text || ']')::jsonb
`, deliveryPersonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]ordertypes.OrderResponse, 0)
	for rows.Next() {
		var o ordertypes.OrderResponse
		if err := rows.Scan(&o.ID, &o.User, &o.Restaurant, &o.Address, &o.Status, &o.Total); err != nil {
			return nil, err
		}

		itemRows, err := r.db.QueryContext(ctx, `
			SELECT m.name, oi.quantity
			FROM   order_items oi
			JOIN   menu m ON m.id = oi.product_id
			WHERE  oi.order_id = $1
		`, o.ID)
		if err != nil {
			return nil, err
		}
		items := make([]ordertypes.ItemResponse, 0)
		for itemRows.Next() {
			var it ordertypes.ItemResponse
			if err := itemRows.Scan(&it.Name, &it.Quantity); err != nil {
				itemRows.Close()
				return nil, err
			}
			items = append(items, it)
		}
		itemRows.Close()
		o.Items = items

		orders = append(orders, o)
	}
	return orders, nil
}
