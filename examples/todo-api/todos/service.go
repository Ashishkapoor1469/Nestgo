package todos

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// TodoService provides business logic for todos.
type TodoService struct {
	mu    sync.RWMutex
	todos map[string]*Todo
	seq   int
}

// NewTodoService creates a new todo service.
func NewTodoService() *TodoService {
	svc := &TodoService{
		todos: make(map[string]*Todo),
	}

	// Seed with sample data.
	svc.Create(CreateTodoDTO{Title: "Learn NestGo framework", Description: "Explore modules, DI, and controllers"})
	svc.Create(CreateTodoDTO{Title: "Build a REST API", Description: "Create a production-ready API with NestGo"})
	svc.Create(CreateTodoDTO{Title: "Add WebSocket support", Description: "Implement real-time features"})

	return svc
}

// FindAll returns paginated todos.
func (s *TodoService) FindAll(page, limit int) ([]*Todo, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	all := make([]*Todo, 0, len(s.todos))
	for _, t := range s.todos {
		all = append(all, t)
	}

	// Sort by creation date descending.
	sort.Slice(all, func(i, j int) bool {
		return all[i].CreatedAt.After(all[j].CreatedAt)
	})

	total := len(all)
	start := (page - 1) * limit
	if start >= total {
		return []*Todo{}, total
	}
	end := start + limit
	if end > total {
		end = total
	}

	return all[start:end], total
}

// FindByID finds a todo by ID.
func (s *TodoService) FindByID(id string) (*Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	todo, ok := s.todos[id]
	if !ok {
		return nil, fmt.Errorf("todo not found: %s", id)
	}
	return todo, nil
}

// Create creates a new todo.
func (s *TodoService) Create(dto CreateTodoDTO) *Todo {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.seq++
	todo := &Todo{
		ID:          fmt.Sprintf("%d", s.seq),
		Title:       dto.Title,
		Description: dto.Description,
		Completed:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	s.todos[todo.ID] = todo
	return todo
}

// Update updates a todo.
func (s *TodoService) Update(id string, dto UpdateTodoDTO) (*Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	todo, ok := s.todos[id]
	if !ok {
		return nil, fmt.Errorf("todo not found: %s", id)
	}

	if dto.Title != nil {
		todo.Title = *dto.Title
	}
	if dto.Description != nil {
		todo.Description = *dto.Description
	}
	if dto.Completed != nil {
		todo.Completed = *dto.Completed
	}
	todo.UpdatedAt = time.Now()

	return todo, nil
}

// Toggle toggles a todo's completed status.
func (s *TodoService) Toggle(id string) (*Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	todo, ok := s.todos[id]
	if !ok {
		return nil, fmt.Errorf("todo not found: %s", id)
	}

	todo.Completed = !todo.Completed
	todo.UpdatedAt = time.Now()

	return todo, nil
}

// Delete removes a todo.
func (s *TodoService) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.todos[id]; !ok {
		return fmt.Errorf("todo not found: %s", id)
	}

	delete(s.todos, id)
	return nil
}
