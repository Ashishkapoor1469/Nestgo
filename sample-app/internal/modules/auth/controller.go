package auth

import (
	"net/http"

	"github.com/Ashishkapoor1469/Nestgo/common"
)

// AuthController handles authentication endpoints.
type AuthController struct {
	service *AuthService
}

// NewAuthController creates a new auth controller.
func NewAuthController(service *AuthService) *AuthController {
	return &AuthController{service: service}
}

// Prefix returns the route prefix.
func (c *AuthController) Prefix() string {
	return "/auth"
}

// Routes returns the controller's routes.
func (c *AuthController) Routes() []common.Route {
	return []common.Route{
		{Method: http.MethodPost, Path: "/register", Handler: c.Register, Summary: "Register a new user"},
		{Method: http.MethodPost, Path: "/login", Handler: c.Login, Summary: "Login and get JWT token"},
		{Method: http.MethodGet, Path: "/profile", Handler: c.Profile, Summary: "Get current user profile"},
	}
}

// Register handles user registration.
func (c *AuthController) Register(ctx *common.Context) error {
	var body RegisterRequest
	if err := ctx.Bind(&body); err != nil {
		return ctx.ValidationErrorResponse(err)
	}

	user, err := c.service.Register(body.Name, body.Email, body.Password)
	if err != nil {
		return common.Conflict(err.Error())
	}

	token, err := c.service.GenerateToken(user.ID, user.Email)
	if err != nil {
		return common.InternalError("failed to generate token")
	}

	return ctx.Created(map[string]any{
		"user":  user.ToResponse(),
		"token": token,
	})
}

// Login handles user login.
func (c *AuthController) Login(ctx *common.Context) error {
	var body LoginRequest
	if err := ctx.Bind(&body); err != nil {
		return ctx.ValidationErrorResponse(err)
	}

	user, err := c.service.Login(body.Email, body.Password)
	if err != nil {
		return common.Unauthorized(err.Error())
	}

	token, err := c.service.GenerateToken(user.ID, user.Email)
	if err != nil {
		return common.InternalError("failed to generate token")
	}

	return ctx.OK(map[string]any{
		"user":  user.ToResponse(),
		"token": token,
	})
}

// Profile returns the current user profile.
func (c *AuthController) Profile(ctx *common.Context) error {
	userID, ok := ctx.Get("user_id")
	if !ok {
		return common.Unauthorized("authentication required")
	}

	user, err := c.service.FindByID(userID.(string))
	if err != nil {
		return common.NotFound("user not found")
	}

	return ctx.OK(user.ToResponse())
}
