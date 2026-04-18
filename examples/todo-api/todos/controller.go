package todos

import (
	"net/http"

	"github.com/Ashishkapoor1469/Nestgo/common"
)

// TodoController handles todo HTTP endpoints.
type TodoController struct {
	service *TodoService
}

// NewTodoController creates a new TodoController.
func NewTodoController(service *TodoService) *TodoController {
	return &TodoController{service: service}
}

// Prefix returns the route prefix.
func (c *TodoController) Prefix() string {
	return "/todos"
}

// Routes defines the controller's routes.
func (c *TodoController) Routes() []common.Route {
	return []common.Route{
		{Method: http.MethodGet, Path: "/", Handler: c.FindAll, Summary: "List all todos"},
		{Method: http.MethodGet, Path: "/{id}", Handler: c.FindOne, Summary: "Get todo by ID"},
		{Method: http.MethodPost, Path: "/", Handler: c.Create, Summary: "Create a new todo"},
		{Method: http.MethodPut, Path: "/{id}", Handler: c.Update, Summary: "Update a todo"},
		{Method: http.MethodPatch, Path: "/{id}/toggle", Handler: c.Toggle, Summary: "Toggle todo completion"},
		{Method: http.MethodDelete, Path: "/{id}", Handler: c.Delete, Summary: "Delete a todo"},
	}
}

// FindAll returns all todos.
func (c *TodoController) FindAll(ctx *common.Context) error {
	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 20)

	items, total := c.service.FindAll(page, limit)
	return ctx.Paginated(items, total, page, limit)
}

// FindOne returns a single todo.
func (c *TodoController) FindOne(ctx *common.Context) error {
	id := ctx.Param("id")

	todo, err := c.service.FindByID(id)
	if err != nil {
		return common.NotFound("todo not found")
	}

	return ctx.OK(todo)
}

// Create creates a new todo.
func (c *TodoController) Create(ctx *common.Context) error {
	var dto CreateTodoDTO
	if err := ctx.Bind(&dto); err != nil {
		return common.BadRequest(err.Error())
	}

	todo := c.service.Create(dto)
	return ctx.Created(todo)
}

// Update updates a todo.
func (c *TodoController) Update(ctx *common.Context) error {
	id := ctx.Param("id")

	var dto UpdateTodoDTO
	if err := ctx.Bind(&dto); err != nil {
		return common.BadRequest(err.Error())
	}

	todo, err := c.service.Update(id, dto)
	if err != nil {
		return common.NotFound("todo not found")
	}

	return ctx.OK(todo)
}

// Toggle toggles a todo's completed status.
func (c *TodoController) Toggle(ctx *common.Context) error {
	id := ctx.Param("id")

	todo, err := c.service.Toggle(id)
	if err != nil {
		return common.NotFound("todo not found")
	}

	return ctx.OK(todo)
}

// Delete removes a todo.
func (c *TodoController) Delete(ctx *common.Context) error {
	id := ctx.Param("id")

	if err := c.service.Delete(id); err != nil {
		return common.NotFound("todo not found")
	}

	return ctx.NoContent()
}
