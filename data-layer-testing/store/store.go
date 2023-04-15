package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/amanbolat/pkg/sql"
)

type User struct {
	ID   string
	Name string
}

type Store struct {
	conn pkgsql.Database
}

func NewStore(conn *sql.DB) *Store {
	return &Store{
		conn: conn,
	}
}

func (s Store) CreateUser(ctx context.Context, user User) error {
	res, err := s.conn.ExecContext(ctx, `INSERT INTO "user" (id, name) VALUES ($1, $2);`, user.ID, user.Name)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected != 1 {
		return fmt.Errorf("expected 1 row to be affected, got %d", affected)
	}

	return nil
}

func (s Store) GetUserByID(ctx context.Context, id string) (User, error) {
	var user User

	rows, err := s.conn.QueryContext(ctx, `SELECT id, name FROM "user" WHERE id = $1`, id)
	if err != nil {
		return User{}, err
	}

	defer rows.Close()

	if !rows.Next() {
		return User{}, fmt.Errorf("user with id %s not found", id)
	}

	err = rows.Scan(&user.ID, &user.Name)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
