package Menu

import (
	"SEProject/Menu/types"
	"context"
	"database/sql"
)

type Repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) *Repository { return &Repository{db} }

// SELECT
func (r *Repository) GetMenuByRestaurantID(ctx context.Context, rid int) ([]types.Menu, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, restaurant_id, name, price, image_url
		FROM menu WHERE restaurant_id=$1`, rid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]types.Menu, 0)
	for rows.Next() {
		var m types.Menu
		if err := rows.Scan(&m.ID, &m.RestaurantID, &m.Name, &m.Price, &m.ImageURL); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, nil
}

// INSERT → id döndür
func (r *Repository) CreateMenuItem(ctx context.Context, m *types.Menu) (int, error) {
	var id int
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO menu (restaurant_id, name, price)
		VALUES ($1,$2,$3) RETURNING id`,
		m.RestaurantID, m.Name, m.Price).Scan(&id)
	return id, err
}

// UPDATE (name/price)
func (r *Repository) UpdateMenuItem(ctx context.Context, id int, m *types.Menu) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE menu SET name=$1, price=$2 WHERE id=$3`,
		m.Name, m.Price, id)
	return err
}

func (r *Repository) DeleteMenuItem(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM menu WHERE id=$1`, id)
	return err
}

// image_url set
func (r *Repository) UpdateMenuItemImage(ctx context.Context, id int, url string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE menu SET image_url=$1 WHERE id=$2`, url, id)
	return err
}
