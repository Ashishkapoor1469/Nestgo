package users

import (
	"github.com/Ashishkapoor1469/Nestgo/common"
)

type UsersController struct {
	service *UsersService
}

func NewUsersController(svc *UsersService) *UsersController {
	return &UsersController{service: svc}
}

func (c *UsersController) Prefix() string { return "/users" }

func (c *UsersController) Routes() []common.Route {
	return []common.Route{
		{Method: "GET", Path: "/", Handler: c.FindAll, Summary: "List all users"},
		{Method: "GET", Path: "/{id}", Handler: c.FindOne, Summary: "Get user by ID"},
		{Method: "POST", Path: "/", Handler: c.Create, Summary: "Create a user"},
	}
}

func (c *UsersController) FindAll(ctx *common.Context) error {
	users := c.service.FindAll()
	return ctx.OK(map[string]any{
		"data":  users,
		"total": len(users),
	})
}

func (c *UsersController) FindOne(ctx *common.Context) error {
	id := ctx.Param("id")
	user, err := c.service.FindByID(id)
	if err != nil {
		return common.NotFound("user not found")
	}
	return ctx.OK(user)
}

func (c *UsersController) Create(ctx *common.Context) error {
	var dto CreateUserDTO
	if err := ctx.Bind(&dto); err != nil {
		return common.BadRequest(err.Error())
	}
	user, err := c.service.Create(dto)
	if err != nil {
		return common.BadRequest(err.Error())
	}
	return ctx.Created(user)
}
