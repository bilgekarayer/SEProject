// Restaurant/repository.go
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
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, location, cuisine, description, avg_price, rating, image_url
		FROM restaurants`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []types.Restaurant
	for rows.Next() {
		var res types.Restaurant
		if err := rows.Scan(&res.ID, &res.Name, &res.Location, &res.Cuisine,
			&res.Description, &res.AvgPrice, &res.Rating, &res.ImageURL); err != nil {
			return nil, err
		}
		list = append(list, res)
	}
	return list, nil
}

func (r *Repository) Create(ctx context.Context, rest *types.Restaurant) (int, error) {
	var id int
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO restaurants
		  (name, description, location, cuisine, avg_price, rating)
		VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING id`,
		rest.Name, rest.Description, rest.Location, rest.Cuisine,
		rest.AvgPrice, rest.Rating).Scan(&id)
	return id, err
}

func (r *Repository) Update(ctx context.Context, id int, rest *types.Restaurant) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE restaurants SET
		    name=$1, description=$2, location=$3, cuisine=$4,
		    avg_price=$5, rating=$6, image_url=$7
		WHERE id=$8`,
		rest.Name, rest.Description, rest.Location, rest.Cuisine,
		rest.AvgPrice, rest.Rating, rest.ImageURL, id)
	return err
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM restaurants WHERE id=$1`, id)
	return err
}

func (r *Repository) UpdateRestaurantImage(ctx context.Context, id int, url string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE restaurants SET image_url=$1 WHERE id=$2`, url, id)
	return err
}
