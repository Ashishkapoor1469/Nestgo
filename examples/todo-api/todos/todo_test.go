package todos

import (
	"testing"
)

func TestNewTodoService(t *testing.T) {
	service := NewTodoService()
	if service == nil {
		t.Fatal("expected service to be non-nil")
	}
}

func TestTodoService_Create(t *testing.T) {
	service := NewTodoService()
	initialCount := len(service.todos)

	todo := service.Create(CreateTodoDTO{
		Title:       "Test Todo",
		Description: "A test todo item",
	})

	if todo.ID == "" {
		t.Fatal("expected ID to be set")
	}
	if todo.Title != "Test Todo" {
		t.Fatalf("expected title 'Test Todo', got '%s'", todo.Title)
	}
	if todo.Completed {
		t.Fatal("expected new todo to not be completed")
	}
	if len(service.todos) != initialCount+1 {
		t.Fatalf("expected %d todos, got %d", initialCount+1, len(service.todos))
	}
}

func TestTodoService_FindAll(t *testing.T) {
	service := &TodoService{todos: make(map[string]*Todo)}

	for i := 0; i < 15; i++ {
		service.Create(CreateTodoDTO{Title: "Todo"})
	}

	// Page 1.
	items, total := service.FindAll(1, 10)
	if total != 15 {
		t.Fatalf("expected total 15, got %d", total)
	}
	if len(items) != 10 {
		t.Fatalf("expected 10 items on page 1, got %d", len(items))
	}

	// Page 2.
	items, _ = service.FindAll(2, 10)
	if len(items) != 5 {
		t.Fatalf("expected 5 items on page 2, got %d", len(items))
	}
}

func TestTodoService_FindByID(t *testing.T) {
	service := &TodoService{todos: make(map[string]*Todo)}
	created := service.Create(CreateTodoDTO{Title: "Find me"})

	found, err := service.FindByID(created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if found.Title != "Find me" {
		t.Fatalf("expected title 'Find me', got '%s'", found.Title)
	}

	// Not found.
	_, err = service.FindByID("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent ID")
	}
}

func TestTodoService_Update(t *testing.T) {
	service := &TodoService{todos: make(map[string]*Todo)}
	created := service.Create(CreateTodoDTO{Title: "Original"})

	newTitle := "Updated"
	updated, err := service.Update(created.ID, UpdateTodoDTO{Title: &newTitle})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if updated.Title != "Updated" {
		t.Fatalf("expected title 'Updated', got '%s'", updated.Title)
	}
}

func TestTodoService_Toggle(t *testing.T) {
	service := &TodoService{todos: make(map[string]*Todo)}
	created := service.Create(CreateTodoDTO{Title: "Toggle me"})

	if created.Completed {
		t.Fatal("expected initial state to be not completed")
	}

	toggled, err := service.Toggle(created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !toggled.Completed {
		t.Fatal("expected todo to be completed after toggle")
	}

	toggled, _ = service.Toggle(created.ID)
	if toggled.Completed {
		t.Fatal("expected todo to be not completed after second toggle")
	}
}

func TestTodoService_Delete(t *testing.T) {
	service := &TodoService{todos: make(map[string]*Todo)}
	created := service.Create(CreateTodoDTO{Title: "Delete me"})

	err := service.Delete(created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = service.FindByID(created.ID)
	if err == nil {
		t.Fatal("expected error after deletion")
	}

	// Delete nonexistent.
	err = service.Delete("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent ID")
	}
}

func TestCreateTodoDTO_Validate(t *testing.T) {
	// Valid.
	dto := CreateTodoDTO{Title: "Valid title"}
	if err := dto.Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Empty title.
	dto = CreateTodoDTO{Title: ""}
	if err := dto.Validate(); err == nil {
		t.Fatal("expected validation error for empty title")
	}
}
