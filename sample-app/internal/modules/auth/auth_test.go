package auth

import (
	"testing"
)

func TestAuthService_Register(t *testing.T) {
	service := NewAuthService("test-secret")

	user, err := service.Register("John", "john@example.com", "password123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.Name != "John" {
		t.Fatalf("expected name John, got %s", user.Name)
	}
	if user.Email != "john@example.com" {
		t.Fatalf("expected email john@example.com, got %s", user.Email)
	}
}

func TestAuthService_DuplicateEmail(t *testing.T) {
	service := NewAuthService("test-secret")

	_, _ = service.Register("John", "john@example.com", "password123")
	_, err := service.Register("Jane", "john@example.com", "password456")
	if err == nil {
		t.Fatal("expected duplicate email error")
	}
}

func TestAuthService_GenerateAndValidateToken(t *testing.T) {
	service := NewAuthService("test-secret")

	token, err := service.GenerateToken("1", "john@example.com")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	claims, err := service.ValidateToken(token)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if claims.UserID != "1" {
		t.Fatalf("expected user ID 1, got %s", claims.UserID)
	}
	if claims.Email != "john@example.com" {
		t.Fatalf("expected email john@example.com, got %s", claims.Email)
	}
}

func TestAuthService_FindByID(t *testing.T) {
	service := NewAuthService("test-secret")

	created, _ := service.Register("John", "john@example.com", "password123")
	found, err := service.FindByID(created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if found.ID != created.ID {
		t.Fatalf("expected ID %s, got %s", created.ID, found.ID)
	}
}
