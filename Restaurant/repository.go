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
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, location, cuisine, description, avg_price, rating FROM restaurants")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var restaurants []types.Restaurant
	for rows.Next() {
		var res types.Restaurant
		if err := rows.Scan(&res.ID, &res.Name, &res.Location, &res.Cuisine, &res.Description, &res.AvgPrice, &res.Rating); err != nil {
			return nil, err
		}
		restaurants = append(restaurants, res)
	}
	return restaurants, nil
}

func (r *Repository) Create(ctx context.Context, rest *types.Restaurant) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO restaurants (name, description, location, cuisine, avg_price, rating)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, rest.Name, rest.Description, rest.Location, rest.Cuisine, rest.AvgPrice, rest.Rating)
	return err
}

func (r *Repository) Update(ctx context.Context, id int, rest *types.Restaurant) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE restaurants
		SET name = $1,
			description = $2,
			location = $3,
			cuisine = $4,
			avg_price = $5,
			rating = $6
		WHERE id = $7
	`, rest.Name, rest.Description, rest.Location, rest.Cuisine, rest.AvgPrice, rest.Rating, id)
	return err
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM restaurants WHERE id=$1", id)
	return err
}
