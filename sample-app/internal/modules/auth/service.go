package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/Ashishkapoor1469/Nestgo/database"
)

// User represents a user entity in PostgreSQL.
type User struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// UserResponse is the public user representation.
type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

// ToResponse converts a Postgres User to a public response.
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}

// AuthService provides PostgreSQL authentication logic.
type AuthService struct {
	db        *database.Database
	jwtSecret string
}

// NewAuthService creates a new auth service. Takes *database.Database automatically via DI.
func NewAuthService(db *database.Database) *AuthService {
	// Taking the secret from env would be ideal in prod!
	return &AuthService{
		db:        db,
		jwtSecret: "super-secret-key-in-production",
	}
}

// Register creates a new user account in Postgres.
func (s *AuthService) Register(name, email, password string) (*User, error) {
	if s.db == nil {
		return nil, fmt.Errorf("database connection not established")
	}

	hash := hashPassword(password)
	
	var user User
	var idInt int
	err := s.db.DB().QueryRow(
		`INSERT INTO users (name, email, password_hash) 
		 VALUES ($1, $2, $3) 
		 RETURNING id, name, email, created_at, updated_at`,
		name, email, hash,
	).Scan(&idInt, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("email may already be registered: %v", err)
	}
	
	user.ID = fmt.Sprintf("%d", idInt)
	return &user, nil
}

// Login authenticates a user by email and password using Postgres.
func (s *AuthService) Login(email, password string) (*User, error) {
	if s.db == nil {
		return nil, fmt.Errorf("database connection not established")
	}

	var user User
	var idInt int
	err := s.db.DB().QueryRow(
		`SELECT id, name, email, password_hash, created_at, updated_at 
		 FROM users WHERE email = $1`,
		email,
	).Scan(&idInt, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("invalid credentials")
	} else if err != nil {
		return nil, err
	}

	if user.PasswordHash != hashPassword(password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	user.ID = fmt.Sprintf("%d", idInt)
	return &user, nil
}

// FindByID returns a user by Postgres ID.
func (s *AuthService) FindByID(id string) (*User, error) {
	if s.db == nil {
		return nil, fmt.Errorf("database connection not established")
	}

	var user User
	var idInt int
	err := s.db.DB().QueryRow(
		`SELECT id, name, email, created_at, updated_at 
		 FROM users WHERE id = $1`,
		id,
	).Scan(&idInt, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	} else if err != nil {
		return nil, err
	}

	user.ID = fmt.Sprintf("%d", idInt)
	return &user, nil
}

// GenerateToken creates a JWT token for the user.
func (s *AuthService) GenerateToken(userID, email string) (string, error) {
	return GenerateJWT(s.jwtSecret, userID, email)
}

// ValidateToken validates a JWT token and returns claims.
func (s *AuthService) ValidateToken(token string) (*Claims, error) {
	return ValidateJWT(s.jwtSecret, token)
}

// hashPassword creates a SHA256 hash of the password with a salt.
func hashPassword(password string) string {
	salt := make([]byte, 16)
	_, _ = rand.Read(salt)
	h := sha256.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}
