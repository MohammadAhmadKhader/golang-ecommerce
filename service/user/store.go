package user

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/mohammadahmadkhader/golang-ecommerce/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE email = ?", strings.ToLower(strings.TrimSpace(email)))
	if err != nil {
		return nil, err
	}

	user := new(types.User)
	for rows.Next() {
		user, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if user.ID == 0 {
		return nil, fmt.Errorf("user was not found")
	}

	return user, nil
}

func scanRowIntoUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)

	err := rows.Scan(userAllFieldsScanner(user))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Store) GetUserByID(id int) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	user := new(types.User)
	for rows.Next() {
		user, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if user.ID == 0 {
		return nil, fmt.Errorf("user was not found")
	}

	return user, nil
}

func (s *Store) CreateUser(user types.User) error {
	_, err := s.db.Exec(`
	INSERT INTO users (firstName, lastName, email, password)
	VALUES (?,?,?,?)`, strings.TrimSpace(user.FirstName), strings.TrimSpace(user.LastName), strings.ToLower(strings.TrimSpace(user.Email)), user.Password)
	if err != nil {
		return err
	}

	return nil
}

func userAllFieldsScanner(user *types.User) (*int, *string, *string, *string, *string, *time.Time, *time.Time) {
	return &user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt
}
