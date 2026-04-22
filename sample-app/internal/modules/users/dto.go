package users

import (
	"github.com/Ashishkapoor1469/Nestgo/common"
)

// CreateUserDTO is the request body for creating a user.
type CreateUserDTO struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// Validate validates the DTO.
func (d *CreateUserDTO) Validate() error {
	return common.NewValidator().
		Required("name", d.Name).
		MinLength("name", d.Name, 2).
		MaxLength("name", d.Name, 100).
		Validate()
}

// UpdateUserDTO is the request body for updating a user.
type UpdateUserDTO struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// Validate validates the DTO.
func (d *UpdateUserDTO) Validate() error {
	v := common.NewValidator()
	if d.Name != nil {
		v.MinLength("name", *d.Name, 2).
			MaxLength("name", *d.Name, 100)
	}
	return v.Validate()
}
