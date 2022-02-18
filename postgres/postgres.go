package postgres

import (
	"database/sql"

	"github.com/chingizof/golang-test-assignment/cmd/simple-service/main.go"
	_ "github.com/lib/pq"
)

type UserService struct {
	DB *sql.DB
}

func (s *UserService) User(id int) (*main.Item, error) {
	var u main.User
	row := db.QueryRow(`SELECT id, name FROM users WHERE id = $1`, id)
	if row.Scan(&u.ID, &u.Name); err != nil {
		return nil, err
	}
	return &u, nil
}
