package User

import (
	"SEProject/User/types"
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
	row := r.db.QueryRow(`
		SELECT 
			u.id, u.username, u.password, 
			u.first_name, u.last_name, 
			u.created_at, 
			r.id, r.name
		FROM users u
		JOIN roles r ON u.role_id = r.id
		WHERE u.username = $1
	`, username)

	var u types.User
	if err := row.Scan(
		&u.ID,
		&u.Username,
		&u.Password,
		&u.FirstName,
		&u.LastName,
		&u.CreatedAt,
		&u.RoleID,
		&u.RoleName,
	); err != nil {
		return nil, err
	}
	return &u, nil
}
func (r *Repository) Create(ctx context.Context, user *types.User) error {
	_, err := r.db.ExecContext(ctx, `
  INSERT INTO users (username, password, role_id, first_name, last_name)
  VALUES ($1, $2, $3, $4, $5)
`, user.Username, user.Password, user.RoleID, user.FirstName, user.LastName)

	return err
}

func (r *Repository) UpdateAllowed(ctx context.Context, id int, u *types.UpdateUserRequest) error {
	res, err := r.db.ExecContext(ctx, `
	  UPDATE users 
	    SET username   = $1,
	        password   = COALESCE(NULLIF($2,''), password),
	        first_name = $3,
	        last_name  = $4,
	        role_id    = $5
	  WHERE id = $6
	`, u.Username, u.Password, u.FirstName, u.LastName, u.RoleID, id)
	if err != nil {
		return err
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
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

func (r *Repository) GetAllUsers(ctx context.Context) ([]*types.User, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT 
			u.id, u.username, u.password,
			u.first_name, u.last_name, u.created_at,
			r.id, r.name
		FROM users u
		JOIN roles r ON u.role_id = r.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*types.User
	for rows.Next() {
		var u types.User
		if err := rows.Scan(
			&u.ID,
			&u.Username,
			&u.Password,
			&u.FirstName,
			&u.LastName,
			&u.CreatedAt,
			&u.RoleID,
			&u.RoleName,
		); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, nil
}
