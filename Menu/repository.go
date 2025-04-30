package Menu

import (
	"SEProject/Menu/types"
	"context"
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetMenuByRestaurantID(ctx context.Context, restaurantID int) ([]types.Menu, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, restaurant_id, name, price FROM menu WHERE restaurant_id = $1`, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var menu []types.Menu
	for rows.Next() {
		var item types.Menu
		if err := rows.Scan(&item.ID, &item.RestaurantID, &item.Name, &item.Price); err != nil {
			return nil, err
		}
		menu = append(menu, item)
	}
	return menu, nil
}

func (r *Repository) CreateMenuItem(ctx context.Context, item *types.Menu) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO menu (restaurant_id, name, price) VALUES ($1, $2, $3)`,
		item.RestaurantID, item.Name, item.Price)
	return err
}

func (r *Repository) UpdateMenuItem(ctx context.Context, id int, item *types.Menu) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE menu SET name=$1, price=$2 WHERE id=$3`,
		item.Name, item.Price, id)
	return err
}

func (r *Repository) DeleteMenuItem(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM menu WHERE id=$1`, id)
	return err
}
