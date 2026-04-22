package users

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Ashishkapoor1469/Nestgo/database"
)

// UsersService provides PostgreSQL business logic for users.
type UsersService struct {
	db *database.Database
}

// NewUsersService creates a new service that automatically receives the active *database.Database from DI.
func NewUsersService(db *database.Database) *UsersService {
	return &UsersService{
		db: db,
	}
}

// FindAll returns paginated users from Postgres.
func (s *UsersService) FindAll(page, limit int) ([]*User, int, error) {
	if s.db == nil {
		return nil, 0, fmt.Errorf("database connection not established")
	}

	offset := (page - 1) * limit
	rows, err := s.db.DB().Query(
		`SELECT id, name, email, created_at, updated_at 
		 FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var all []*User
	for rows.Next() {
		var idInt int
		u := &User{}
		if err := rows.Scan(&idInt, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, 0, err
		}
		u.ID = fmt.Sprintf("%d", idInt)
		all = append(all, u)
	}

	var total int
	s.db.DB().QueryRow(`SELECT COUNT(*) FROM users`).Scan(&total)

	return all, total, nil
}

// FindByID returns a user by Postgres ID.
func (s *UsersService) FindByID(id string) (*User, error) {
	if s.db == nil {
		return nil, fmt.Errorf("database connection not established")
	}

	var u User
	var idInt int
	err := s.db.DB().QueryRow(
		`SELECT id, name, email, created_at, updated_at 
		 FROM users WHERE id = $1`,
		id,
	).Scan(&idInt, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found: %s", id)
	} else if err != nil {
		return nil, err
	}
	
	u.ID = fmt.Sprintf("%d", idInt)
	return &u, nil
}

// Update updates a user's details.
func (s *UsersService) Update(id string, dto UpdateUserDTO) (*User, error) {
	if s.db == nil {
		return nil, fmt.Errorf("database connection not established")
	}

	if dto.Name != nil {
		_, err := s.db.DB().Exec(
			`UPDATE users SET name = $1, updated_at = $2 WHERE id = $3`,
			*dto.Name, time.Now(), id,
		)
		if err != nil {
			return nil, err
		}
	}

	return s.FindByID(id)
}

// Delete removes a user.
func (s *UsersService) Delete(id string) error {
	if s.db == nil {
		return fmt.Errorf("database connection not established")
	}

	_, err := s.db.DB().Exec(`DELETE FROM users WHERE id = $1`, id)
	return err
}
