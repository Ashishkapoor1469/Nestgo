package products

import (
	"net/http"

	"github.com/Ashishkapoor1469/Nestgo/common"
)

// ProductsController handles HTTP requests for products.
type ProductsController struct {
	service *ProductsService
}

// NewProductsController creates a new controller.
func NewProductsController(service *ProductsService) *ProductsController {
	return &ProductsController{service: service}
}

// Prefix returns the route prefix.
func (c *ProductsController) Prefix() string {
	return "/products"
}

// Routes returns the controller's route definitions.
func (c *ProductsController) Routes() []common.Route {
	return []common.Route{
		{Method: http.MethodGet, Path: "/", Handler: c.FindAll, Summary: "List all products"},
		{Method: http.MethodGet, Path: "/{id}", Handler: c.FindOne, Summary: "Get a product by ID"},
		{Method: http.MethodPost, Path: "/", Handler: c.Create, Summary: "Create a new product"},
		{Method: http.MethodPut, Path: "/{id}", Handler: c.Update, Summary: "Update a product"},
		{Method: http.MethodDelete, Path: "/{id}", Handler: c.Delete, Summary: "Delete a product"},
	}
}

// FindAll returns all products.
func (c *ProductsController) FindAll(ctx *common.Context) error {
	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 20)

	items, total, err := c.service.FindAll(page, limit)
	if err != nil {
		return common.InternalError(err.Error())
	}

	return ctx.Paginated(items, total, page, limit)
}

// FindOne returns a single product.
func (c *ProductsController) FindOne(ctx *common.Context) error {
	id := ctx.Param("id")

	item, err := c.service.FindByID(id)
	if err != nil {
		return common.NotFound("product not found")
	}

	return ctx.OK(item)
}

// Create creates a new product.
func (c *ProductsController) Create(ctx *common.Context) error {
	var dto CreateProductDTO
	if err := ctx.Bind(&dto); err != nil {
		return common.BadRequest(err.Error())
	}

	item, err := c.service.Create(dto)
	if err != nil {
		return common.InternalError(err.Error())
	}

	return ctx.Created(item)
}

// Update updates a product.
func (c *ProductsController) Update(ctx *common.Context) error {
	id := ctx.Param("id")

	var dto UpdateProductDTO
	if err := ctx.Bind(&dto); err != nil {
		return common.BadRequest(err.Error())
	}

	item, err := c.service.Update(id, dto)
	if err != nil {
		return common.NotFound("product not found")
	}

	return ctx.OK(item)
}

// Delete removes a product.
func (c *ProductsController) Delete(ctx *common.Context) error {
	id := ctx.Param("id")

	if err := c.service.Delete(id); err != nil {
		return common.NotFound("product not found")
	}

	return ctx.NoContent()
}
