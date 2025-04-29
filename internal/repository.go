package internal

import (
	"SEProject/internal/types"
	"context"
	"database/sql"
	"errors"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) GetByUsername(ctx context.Context, username string) (*types.User, error) {
	row := r.db.QueryRow("SELECT id, username, password FROM users WHERE username=$1", username)
	var u types.User
	if err := row.Scan(&u.ID, &u.Username, &u.Password); err != nil {
		return nil, err
	}
	return &u, nil
}
func (r *Repository) Create(ctx context.Context, user *types.User) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO users (username,password) 
		VALUES ($1, $2)
	`, user.Username, user.Password)

	return err
}

func (r *Repository) Update(ctx context.Context, id int, user *types.User) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE users 
		SET username = $1, password = $2 
		WHERE id = $3
	`, user.Username, user.Password, id)

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("kullanıcı bulunamadı")
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, `
		DELETE FROM users WHERE id = $1
	`, id)

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("kullanıcı bulunamadı")
	}

	return nil
}
