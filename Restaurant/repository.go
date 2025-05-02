package Restaurant

import (
	"SEProject/Restaurant/types"
	"context"
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAll(ctx context.Context) ([]types.Restaurant, error) {
	// ✅ cuisine alanı eklendi
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, location, cuisine FROM restaurants")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var restaurants []types.Restaurant
	for rows.Next() {
		var res types.Restaurant
		// ✅ cuisine alanı Scan'e eklendi
		if err := rows.Scan(&res.ID, &res.Name, &res.Location, &res.Cuisine); err != nil {
			return nil, err
		}
		restaurants = append(restaurants, res)
	}
	return restaurants, nil
}

func (r *Repository) Create(ctx context.Context, rest *types.Restaurant) error {
	// ✅ cuisine alanı INSERT'e eklendi
	_, err := r.db.ExecContext(ctx, "INSERT INTO restaurants (name, location, cuisine) VALUES ($1, $2, $3)", rest.Name, rest.Location, rest.Cuisine)
	return err
}

func (r *Repository) Update(ctx context.Context, id int, rest *types.Restaurant) error {
	// ✅ cuisine alanı UPDATE'e eklendi
	_, err := r.db.ExecContext(ctx, "UPDATE restaurants SET name=$1, location=$2, cuisine=$3 WHERE id=$4", rest.Name, rest.Location, rest.Cuisine, id)
	return err
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM restaurants WHERE id=$1", id)
	return err
}
