package auth

import (
	"net/http"

	"github.com/Ashishkapoor1469/Nestgo/common"
)

// AuthController handles incoming HTTP requests for the /auth routes.
type AuthController struct {
	authService *AuthService
}

// NewAuthController acts as a provider factory function.
// NestGo's DI container will automatically inject the required *AuthService.
func NewAuthController(authService *AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// Prefix returns the route prefix for all endpoints in this controller.
func (c *AuthController) Prefix() string {
	return "/api/auth"
}

// Routes connects HTTP methods and paths to standard handler functions,
// along with any specific middleware or guards for those endpoints.
func (c *AuthController) Routes() []common.Route {
	return []common.Route{
		{
			Method:  http.MethodPost,
			Path:    "/login",
			Handler: c.login,
			Summary: "Login and get token",
		},
		{
			Method:  http.MethodGet,
			Path:    "/profile",
			Handler: c.profile,
			Summary: "Get protected profile data",
			// Apply the custom AuthGuard specifically to this route.
			Guards: []common.Guard{NewAuthGuard()},
		},
	}
}

// login is a handler that processes a login request.
func (c *AuthController) login(ctx *common.Context) error {
	var dto LoginDTO

	// Bind the incoming JSON body to our DTO struct.
	// We use the framework's built-in ctx.Bind method.
	if err := ctx.Bind(&dto); err != nil {
		return common.BadRequest("invalid payload")
	}

	// Call the injected service to handle the logic.
	token, err := c.authService.Login(dto.Username, dto.Password)
	if err != nil {
		return common.Unauthorized(err.Error())
	}

	// Send a standard JSON JSON response.
	// ctx.JSON lets you send any structure as JSON.
	return ctx.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}

// profile is a protected handler that only runs if the Guard passes.
func (c *AuthController) profile(ctx *common.Context) error {
	// The Guard saved the user_id into the request context.
	// We can retrieve it easily.
	userIDVal, _ := ctx.Get("user_id")
	userID, _ := userIDVal.(string)

	// Return some mock user profile data using the standard Context Success response envelope.
	return ctx.Success(map[string]any{
		"id":       userID,
		"username": "admin",
		"role":     "superuser",
		"status":   "active",
	})
}
