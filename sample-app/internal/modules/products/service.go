package products

import (
	"fmt"
	"sync"
	"time"
)

// ProductsService provides business logic for products.
type ProductsService struct {
	mu    sync.RWMutex
	items map[string]*Product
	seq   int
}

// NewProductsService creates a new service.
func NewProductsService() *ProductsService {
	return &ProductsService{
		items: make(map[string]*Product),
	}
}

// FindAll returns paginated products.
func (s *ProductsService) FindAll(page, limit int) ([]*Product, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	all := make([]*Product, 0, len(s.items))
	for _, item := range s.items {
		all = append(all, item)
	}

	total := len(all)
	start := (page - 1) * limit
	if start >= total {
		return []*Product{}, total, nil
	}

	end := start + limit
	if end > total {
		end = total
	}

	return all[start:end], total, nil
}

// FindByID returns a product by ID.
func (s *ProductsService) FindByID(id string) (*Product, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.items[id]
	if !ok {
		return nil, fmt.Errorf("product not found: %s", id)
	}
	return item, nil
}

// Create creates a new product.
func (s *ProductsService) Create(dto CreateProductDTO) (*Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.seq++
	item := &Product{
		ID:        fmt.Sprintf("%d", s.seq),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.items[item.ID] = item
	return item, nil
}

// Update updates a product.
func (s *ProductsService) Update(id string, dto UpdateProductDTO) (*Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.items[id]
	if !ok {
		return nil, fmt.Errorf("product not found: %s", id)
	}

	item.UpdatedAt = time.Now()
	return item, nil
}

// Delete removes a product.
func (s *ProductsService) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.items[id]; !ok {
		return fmt.Errorf("product not found: %s", id)
	}

	delete(s.items, id)
	return nil
}
