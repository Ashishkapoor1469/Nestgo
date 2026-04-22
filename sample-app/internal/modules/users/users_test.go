package users

import (
	"testing"
)

func TestNewUsersService(t *testing.T) {
	service := NewUsersService()
	if service == nil {
		t.Fatal("expected service to be non-nil")
	}
}

func TestUsersService_Create(t *testing.T) {
	service := NewUsersService()

	dto := CreateUserDTO{
		Name: "Test Item",
	}

	item, err := service.Create(dto)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if item.ID == "" {
		t.Fatal("expected ID to be set")
	}
}

func TestUsersService_FindAll(t *testing.T) {
	service := NewUsersService()

	// Create test data.
	for i := 0; i < 5; i++ {
		_, _ = service.Create(CreateUserDTO{Name: "Item"})
	}

	items, total, err := service.FindAll(1, 10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 5 {
		t.Fatalf("expected 5 total, got %d", total)
	}
	if len(items) != 5 {
		t.Fatalf("expected 5 items, got %d", len(items))
	}
}

func TestUsersService_FindByID(t *testing.T) {
	service := NewUsersService()

	created, _ := service.Create(CreateUserDTO{Name: "Test"})

	found, err := service.FindByID(created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if found.ID != created.ID {
		t.Fatalf("expected ID %s, got %s", created.ID, found.ID)
	}
}

func TestUsersService_Delete(t *testing.T) {
	service := NewUsersService()

	created, _ := service.Create(CreateUserDTO{Name: "Test"})
	err := service.Delete(created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = service.FindByID(created.ID)
	if err == nil {
		t.Fatal("expected error after deletion")
	}
}
