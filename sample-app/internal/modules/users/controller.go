package users

import (
	"net/http"

	"github.com/Ashishkapoor1469/Nestgo/common"
)

// UsersController handles HTTP requests for users.
type UsersController struct {
	service *UsersService
}

// NewUsersController creates a new controller.
func NewUsersController(service *UsersService) *UsersController {
	return &UsersController{service: service}
}

// Prefix returns the route prefix.
func (c *UsersController) Prefix() string {
	return "/users"
}

// Routes returns the controller's route definitions.
func (c *UsersController) Routes() []common.Route {
	return []common.Route{
		{Method: http.MethodGet, Path: "/", Handler: c.FindAll, Summary: "List all users"},
		{Method: http.MethodGet, Path: "/{id}", Handler: c.FindOne, Summary: "Get a user by ID"},
		{Method: http.MethodPut, Path: "/{id}", Handler: c.Update, Summary: "Update a user"},
		{Method: http.MethodDelete, Path: "/{id}", Handler: c.Delete, Summary: "Delete a user"},
	}
}

// FindAll returns all users.
func (c *UsersController) FindAll(ctx *common.Context) error {
	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 20)

	items, total, err := c.service.FindAll(page, limit)
	if err != nil {
		return common.InternalError(err.Error())
	}

	return ctx.Paginated(items, total, page, limit)
}

// FindOne returns a single user.
func (c *UsersController) FindOne(ctx *common.Context) error {
	id := ctx.Param("id")

	item, err := c.service.FindByID(id)
	if err != nil {
		return common.NotFound("user not found")
	}

	return ctx.OK(item)
}


// Update updates a user.
func (c *UsersController) Update(ctx *common.Context) error {
	id := ctx.Param("id")

	var dto UpdateUserDTO
	if err := ctx.Bind(&dto); err != nil {
		return common.BadRequest(err.Error())
	}

	item, err := c.service.Update(id, dto)
	if err != nil {
		return common.NotFound("user not found")
	}

	return ctx.OK(item)
}

// Delete removes a user.
func (c *UsersController) Delete(ctx *common.Context) error {
	id := ctx.Param("id")

	if err := c.service.Delete(id); err != nil {
		return common.NotFound("user not found")
	}

	return ctx.NoContent()
}
