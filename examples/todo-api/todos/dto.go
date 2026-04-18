package todos

import (
	"time"

	"github.com/nestgo/nestgo/common"
)

// Todo is the todo entity.
type Todo struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// CreateTodoDTO is the request body for creating a todo.
type CreateTodoDTO struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
}

// Validate validates the CreateTodoDTO.
func (d *CreateTodoDTO) Validate() error {
	return common.NewValidator().
		Required("title", d.Title).
		MinLength("title", d.Title, 1).
		MaxLength("title", d.Title, 200).
		MaxLength("description", d.Description, 1000).
		Validate()
}

// UpdateTodoDTO is the request body for updating a todo.
type UpdateTodoDTO struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Completed   *bool   `json:"completed,omitempty"`
}

// Validate validates the UpdateTodoDTO.
func (d *UpdateTodoDTO) Validate() error {
	v := common.NewValidator()
	if d.Title != nil {
		v.MinLength("title", *d.Title, 1).
			MaxLength("title", *d.Title, 200)
	}
	if d.Description != nil {
		v.MaxLength("description", *d.Description, 1000)
	}
	return v.Validate()
}
