package users

import (
	"fmt"
	"sync"
	"time"
)

// User represents a user entity.
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

// CreateUserDTO is the request body for creating a user.
type CreateUserDTO struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UsersService handles in-memory user storage.
type UsersService struct {
	mu    sync.RWMutex
	users map[string]*User
	seq   int
}

func NewUsersService() *UsersService {
	svc := &UsersService{users: make(map[string]*User)}
	// Seed with sample data
	svc.create("Alice Johnson", "alice@example.com")
	svc.create("Bob Smith", "bob@example.com")
	return svc
}

func (s *UsersService) create(name, email string) *User {
	s.seq++
	u := &User{
		ID:        fmt.Sprintf("%d", s.seq),
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
	}
	s.users[u.ID] = u
	return u
}

func (s *UsersService) FindAll() []*User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*User, 0, len(s.users))
	for _, u := range s.users {
		result = append(result, u)
	}
	return result
}

func (s *UsersService) FindByID(id string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if u, ok := s.users[id]; ok {
		return u, nil
	}
	return nil, fmt.Errorf("user %q not found", id)
}

func (s *UsersService) Create(dto CreateUserDTO) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if dto.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	u := s.create(dto.Name, dto.Email)
	return u, nil
}
