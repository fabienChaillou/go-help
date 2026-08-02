package repo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DBTX interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

type User struct {
	ID   int
	Name string
}

type UserRepository struct {
	db DBTX
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*User, error) {
	row := r.db.QueryRow(ctx, "SELECT id, name FROM users WHERE id=$1", id)

	var u User
	err := row.Scan(&u.ID, &u.Name)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetByName(ctx context.Context, name string) ([]User, error) {
	rows, err := r.db.Query(ctx, "SELECT id, name FROM users WHERE name=$1", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}
