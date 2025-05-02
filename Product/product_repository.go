package Product

import (
	"context"
	"database/sql"
)

type ProductRepository struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{DB: db}
}

func (r *ProductRepository) Create(product *Product) error {
	query := `INSERT INTO products (name, price, in_stock) VALUES ($1, $2, $3)`
	_, err := r.DB.ExecContext(context.Background(), query, product.Name, product.Price, product.InStock)
	return err
}

func (r *ProductRepository) Update(product *Product) error {
	query := `UPDATE products SET name=$1, price=$2, in_stock=$3 WHERE id=$4`
	_, err := r.DB.ExecContext(context.Background(), query, product.Name, product.Price, product.InStock, product.ID)
	return err
}

func (r *ProductRepository) Delete(id int) error {
	query := `DELETE FROM products WHERE id=$1`
	_, err := r.DB.ExecContext(context.Background(), query, id)
	return err
}
