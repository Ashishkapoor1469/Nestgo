package products

import (
	"github.com/Ashishkapoor1469/Nestgo/common"
)

// CreateProductDTO is the request body for creating a product.
type CreateProductDTO struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// Validate validates the DTO.
func (d *CreateProductDTO) Validate() error {
	return common.NewValidator().
		Required("name", d.Name).
		MinLength("name", d.Name, 2).
		MaxLength("name", d.Name, 100).
		Validate()
}

// UpdateProductDTO is the request body for updating a product.
type UpdateProductDTO struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// Validate validates the DTO.
func (d *UpdateProductDTO) Validate() error {
	v := common.NewValidator()
	if d.Name != nil {
		v.MinLength("name", *d.Name, 2).
			MaxLength("name", *d.Name, 100)
	}
	return v.Validate()
}
