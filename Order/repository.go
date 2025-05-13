package Order

import (
	"SEProject/Order/types"
	"context"
	"database/sql"
	"log"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) AddToCart(ctx context.Context, it *types.CartItem) error {
	_, err := r.db.ExecContext(ctx, `insert into cart (user_id,product_id,quantity) values ($1,$2,$3)`, it.UserID, it.ProductID, it.Quantity)
	return err
}

func (r *Repository) GetCart(ctx context.Context, uid int) ([]types.CartItem, error) {
	rows, err := r.db.QueryContext(ctx, `select user_id,product_id,quantity from cart where user_id=$1`, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []types.CartItem
	for rows.Next() {
		var it types.CartItem
		if err = rows.Scan(&it.UserID, &it.ProductID, &it.Quantity); err != nil {
			return nil, err
		}
		list = append(list, it)
	}
	return list, nil
}

func (r *Repository) PlaceOrder(ctx context.Context, uid int, req *types.PlaceOrderRequest) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	// 1) role_id = 4 olan tüm kullanıcı id’lerini JSON dizi (string) olarak al
	var deliveryIDs string
	err = tx.QueryRowContext(ctx, `
        SELECT COALESCE(json_agg(id)::text, '[]')
        FROM users
        WHERE role_id = 4
    `).Scan(&deliveryIDs)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 2) Siparişi oluştur
	var oid int
	err = tx.QueryRowContext(ctx, `
        INSERT INTO orders
            (user_id, restaurant_id, address, status, total, delivery_person_id)
        VALUES
            ($1, $2, $3, 'pending', 0, $4::jsonb)
        RETURNING id
    `, uid, req.RestaurantID, req.Address, deliveryIDs).Scan(&oid)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 3) Ürünleri ekle ve toplamı hesapla
	var total float64
	for _, it := range req.Items {
		var price float64
		if err = tx.QueryRowContext(ctx,
			`SELECT price FROM menu WHERE id = $1`,
			it.ProductID,
		).Scan(&price); err != nil {
			tx.Rollback()
			return err
		}

		if _, err = tx.ExecContext(ctx, `
            INSERT INTO order_items (order_id, product_id, quantity, unit_price)
            VALUES ($1, $2, $3, $4)
        `, oid, it.ProductID, it.Quantity, price); err != nil {
			tx.Rollback()
			return err
		}
		total += price * float64(it.Quantity)
	}

	// 4) Toplamı güncelle
	if _, err = tx.ExecContext(ctx,
		`UPDATE orders SET total = $1 WHERE id = $2`,
		total, oid,
	); err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Printf("tx commit error: %v", err)
		return err
	}
	return nil
}

func (r *Repository) GetOrdersByUser(ctx context.Context, uid int) ([]types.OrderResponse, error) {
	rows, err := r.db.QueryContext(ctx, `select o.id,u.username,r.name,o.address,o.status,o.total from orders o join users u on o.user_id=u.id join restaurants r on o.restaurant_id=r.id where o.user_id=$1`, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []types.OrderResponse
	for rows.Next() {
		var o types.OrderResponse
		if err = rows.Scan(&o.ID, &o.User, &o.Restaurant, &o.Address, &o.Status, &o.Total); err != nil {
			return nil, err
		}
		itRows, err := r.db.QueryContext(ctx, `select m.name,oi.quantity from order_items oi join menu m on oi.product_id=m.id where oi.order_id=$1`, o.ID)
		if err != nil {
			return nil, err
		}
		for itRows.Next() {
			var ir types.ItemResponse
			if err = itRows.Scan(&ir.Name, &ir.Quantity); err != nil {
				itRows.Close()
				return nil, err
			}
			o.Items = append(o.Items, ir)
		}
		itRows.Close()
		out = append(out, o)
	}
	return out, nil
}

func (r *Repository) GetOrdersByRestaurant(ctx context.Context, rid int) ([]types.OrderResponse, error) {
	rows, err := r.db.QueryContext(ctx, `select o.id,u.username,r.name,o.address,o.status,o.total from orders o join users u on o.user_id=u.id join restaurants r on o.restaurant_id=r.id where o.restaurant_id=$1`, rid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []types.OrderResponse
	for rows.Next() {
		var o types.OrderResponse
		if err = rows.Scan(&o.ID, &o.User, &o.Restaurant, &o.Address, &o.Status, &o.Total); err != nil {
			return nil, err
		}
		itRows, err := r.db.QueryContext(ctx, `select m.name,oi.quantity from order_items oi join menu m on oi.product_id=m.id where oi.order_id=$1`, o.ID)
		if err != nil {
			return nil, err
		}
		for itRows.Next() {
			var ir types.ItemResponse
			if err = itRows.Scan(&ir.Name, &ir.Quantity); err != nil {
				itRows.Close()
				return nil, err
			}
			o.Items = append(o.Items, ir)
		}
		itRows.Close()
		out = append(out, o)
	}
	return out, nil
}

func (r *Repository) UpdateOrderStatus(ctx context.Context, id int, status string) error {
	_, err := r.db.ExecContext(ctx, `update orders set status=$1 where id=$2`, status, id)
	return err
}

func (r *Repository) GetAllOrders(ctx context.Context) ([]types.OrderResponse, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT  o.id,
		        u.username,
		        r.name,
		        o.address,
		        o.status,
		        o.total
		FROM    orders o
		JOIN    users        u ON u.id = o.user_id
		JOIN    restaurants  r ON r.id = o.restaurant_id
		ORDER BY o.id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]types.OrderResponse, 0)
	for rows.Next() {
		var o types.OrderResponse
		if err = rows.Scan(&o.ID, &o.User, &o.Restaurant, &o.Address, &o.Status, &o.Total); err != nil {
			return nil, err
		}

		itRows, err := r.db.QueryContext(ctx, `
			SELECT m.name, oi.quantity
			FROM   order_items oi
			JOIN   menu m ON m.id = oi.product_id
			WHERE  oi.order_id = $1
		`, o.ID)
		if err != nil {
			return nil, err
		}
		items := make([]types.ItemResponse, 0)
		for itRows.Next() {
			var ir types.ItemResponse
			if err = itRows.Scan(&ir.Name, &ir.Quantity); err != nil {
				itRows.Close()
				return nil, err
			}
			items = append(items, ir)
		}
		itRows.Close()
		o.Items = items

		out = append(out, o)
	}
	return out, nil
}
