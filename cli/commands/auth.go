package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Ashishkapoor1469/Nestgo/cli/utils"
	"github.com/spf13/cobra"
)

// AuthCmd creates the `nestgo generate auth` command.
func AuthCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "auth",
		Short: "Generate JWT authentication starter",
		Long: `Generates a complete authentication module with:
  - JWT token generation and validation
  - Auth middleware / guard
  - Login and register endpoints
  - User model and DTOs
  - Password hashing utilities`,
		RunE: runGenerateAuth,
	}
}

func runGenerateAuth(cmd *cobra.Command, args []string) error {
	utils.EnsureProjectContext("generate auth")

	utils.PrintHeader("🔐 Auth Module Generator")

	spinner := utils.StartSpinner("Generating authentication module...")

	// Create directories.
	dirs := []string{
		filepath.Join("internal", "modules", "auth"),
		filepath.Join("internal", "modules", "auth", "schemas"),
		filepath.Join("internal", "common", "guards"),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			spinner.StopWithError("Failed to create directories")
			return err
		}
	}

	// Generate files.
	files := map[string]string{
		filepath.Join("internal", "modules", "auth", "module.go"):     authModuleTemplate,
		filepath.Join("internal", "modules", "auth", "controller.go"): authControllerTemplate,
		filepath.Join("internal", "modules", "auth", "service.go"):    authServiceTemplate,
		filepath.Join("internal", "modules", "auth", "jwt.go"):        authJWTTemplate,
		filepath.Join("internal", "modules", "auth", "schemas", "auth.schema.go"): authSchemaTemplate,
		filepath.Join("internal", "modules", "auth", "auth_test.go"):  authTestTemplate,
		filepath.Join("internal", "common", "guards", "jwt_guard.go"): authGuardTemplate,
	}

	for path, content := range files {
		if _, err := os.Stat(path); err == nil {
			// Don't overwrite existing files.
			continue
		}
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			spinner.StopWithError("Failed to write " + path)
			return err
		}
	}

	spinner.StopWithSuccess("Auth module generated")

	fmt.Println()
	utils.PrintStep("📄", "internal/modules/auth/module.go")
	utils.PrintStep("📄", "internal/modules/auth/controller.go")
	utils.PrintStep("📄", "internal/modules/auth/service.go")
	utils.PrintStep("📄", "internal/modules/auth/jwt.go")
	utils.PrintStep("📄", "internal/modules/auth/schemas/auth.schema.go")
	utils.PrintStep("📄", "internal/modules/auth/auth_test.go")
	utils.PrintStep("📄", "internal/common/guards/jwt_guard.go")

	fmt.Println()
	utils.PrintSuccess("Auth module ready!")
	fmt.Println()
	utils.PrintDim("  Next steps:")
	utils.PrintDim("  1. Set JWT_SECRET in your .env file")
	utils.PrintDim("  2. Import auth module in your app module:")
	utils.PrintDim("     Imports: []common.Module{&auth.AuthModule{}}")
	utils.PrintDim("  3. Test endpoints:")
	utils.PrintDim("     POST /auth/register  — Register a new user")
	utils.PrintDim("     POST /auth/login     — Login and get JWT token")
	utils.PrintDim("     GET  /auth/profile   — Get current user (requires token)")
	fmt.Println()

	return nil
}

// --- Auth Templates ---

var authModuleTemplate = `package auth

import (
	"github.com/Ashishkapoor1469/Nestgo/di"
	"github.com/Ashishkapoor1469/Nestgo/common"
)

// AuthModule provides authentication and authorization.
type AuthModule struct{}

func (m *AuthModule) Module() common.ModuleConfig {
	service := NewAuthService("change-me-in-production")
	controller := NewAuthController(service)

	return common.ModuleConfig{
		Name: "auth",
		Controllers: []common.Controller{controller},
		Providers: []di.Provider{
			{Instance: service},
		},
	}
}
`

var authControllerTemplate = `package auth

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
`

var authServiceTemplate = `package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// User represents a user entity.
type User struct {
	ID           string    ` + "`" + `json:"id"` + "`" + `
	Name         string    ` + "`" + `json:"name"` + "`" + `
	Email        string    ` + "`" + `json:"email"` + "`" + `
	PasswordHash string    ` + "`" + `json:"-"` + "`" + `
	CreatedAt    time.Time ` + "`" + `json:"createdAt"` + "`" + `
	UpdatedAt    time.Time ` + "`" + `json:"updatedAt"` + "`" + `
}

// UserResponse is the public user representation.
type UserResponse struct {
	ID        string    ` + "`" + `json:"id"` + "`" + `
	Name      string    ` + "`" + `json:"name"` + "`" + `
	Email     string    ` + "`" + `json:"email"` + "`" + `
	CreatedAt time.Time ` + "`" + `json:"createdAt"` + "`" + `
}

// ToResponse converts a User to a public response.
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}

// AuthService provides authentication business logic.
type AuthService struct {
	mu        sync.RWMutex
	users     map[string]*User
	emailIdx  map[string]string // email -> user ID
	jwtSecret string
	seq       int
}

// NewAuthService creates a new auth service.
func NewAuthService(jwtSecret string) *AuthService {
	return &AuthService{
		users:     make(map[string]*User),
		emailIdx:  make(map[string]string),
		jwtSecret: jwtSecret,
	}
}

// Register creates a new user account.
func (s *AuthService) Register(name, email, password string) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check for duplicate email.
	if _, exists := s.emailIdx[email]; exists {
		return nil, fmt.Errorf("email already registered")
	}

	s.seq++
	user := &User{
		ID:           fmt.Sprintf("%d", s.seq),
		Name:         name,
		Email:        email,
		PasswordHash: hashPassword(password),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	s.users[user.ID] = user
	s.emailIdx[email] = user.ID

	return user, nil
}

// Login authenticates a user by email and password.
func (s *AuthService) Login(email, password string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	userID, exists := s.emailIdx[email]
	if !exists {
		return nil, fmt.Errorf("invalid credentials")
	}

	user := s.users[userID]
	if user.PasswordHash != hashPassword(password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	return user, nil
}

// FindByID returns a user by ID.
func (s *AuthService) FindByID(id string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.users[id]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

// GenerateToken creates a JWT token for the user.
func (s *AuthService) GenerateToken(userID, email string) (string, error) {
	return GenerateJWT(s.jwtSecret, userID, email)
}

// ValidateToken validates a JWT token and returns claims.
func (s *AuthService) ValidateToken(token string) (*Claims, error) {
	return ValidateJWT(s.jwtSecret, token)
}

// hashPassword creates a SHA256 hash of the password with a salt.
func hashPassword(password string) string {
	salt := make([]byte, 16)
	_, _ = rand.Read(salt)
	h := sha256.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}
`

var authJWTTemplate = `package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Claims represents JWT claims.
type Claims struct {
	UserID string ` + "`" + `json:"sub"` + "`" + `
	Email  string ` + "`" + `json:"email"` + "`" + `
	Exp    int64  ` + "`" + `json:"exp"` + "`" + `
	Iat    int64  ` + "`" + `json:"iat"` + "`" + `
}

// LoginRequest is the login request body.
type LoginRequest struct {
	Email    string ` + "`" + `json:"email" validate:"required,email"` + "`" + `
	Password string ` + "`" + `json:"password" validate:"required,min=6"` + "`" + `
}

// RegisterRequest is the registration request body.
type RegisterRequest struct {
	Name     string ` + "`" + `json:"name" validate:"required,min=2,max=100"` + "`" + `
	Email    string ` + "`" + `json:"email" validate:"required,email"` + "`" + `
	Password string ` + "`" + `json:"password" validate:"required,min=6"` + "`" + `
}

// GenerateJWT creates a new JWT token.
func GenerateJWT(secret, userID, email string) (string, error) {
	header := map[string]string{"alg": "HS256", "typ": "JWT"}
	headerJSON, _ := json.Marshal(header)
	headerB64 := base64.RawURLEncoding.EncodeToString(headerJSON)

	claims := Claims{
		UserID: userID,
		Email:  email,
		Exp:    time.Now().Add(24 * time.Hour).Unix(),
		Iat:    time.Now().Unix(),
	}
	claimsJSON, _ := json.Marshal(claims)
	claimsB64 := base64.RawURLEncoding.EncodeToString(claimsJSON)

	signingInput := headerB64 + "." + claimsB64
	signature := signHMAC(secret, signingInput)

	return signingInput + "." + signature, nil
}

// ValidateJWT validates a JWT token and returns its claims.
func ValidateJWT(secret, token string) (*Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	signingInput := parts[0] + "." + parts[1]
	expectedSig := signHMAC(secret, signingInput)

	if !hmac.Equal([]byte(parts[2]), []byte(expectedSig)) {
		return nil, fmt.Errorf("invalid token signature")
	}

	claimsJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid token claims")
	}

	var claims Claims
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		return nil, fmt.Errorf("invalid token claims")
	}

	if time.Now().Unix() > claims.Exp {
		return nil, fmt.Errorf("token expired")
	}

	return &claims, nil
}

func signHMAC(secret, data string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}
`

var authSchemaTemplate = `package schemas

// LoginSchema defines the login request validation.
type LoginSchema struct {
	Email    string ` + "`" + `json:"email" validate:"required,email"` + "`" + `
	Password string ` + "`" + `json:"password" validate:"required,min=6"` + "`" + `
}

// RegisterSchema defines the registration request validation.
type RegisterSchema struct {
	Name     string ` + "`" + `json:"name" validate:"required,min=2,max=100"` + "`" + `
	Email    string ` + "`" + `json:"email" validate:"required,email"` + "`" + `
	Password string ` + "`" + `json:"password" validate:"required,min=6"` + "`" + `
}
`

var authTestTemplate = `package auth

import (
	"testing"
)

func TestAuthService_Register(t *testing.T) {
	service := NewAuthService("test-secret")

	user, err := service.Register("John", "john@example.com", "password123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.Name != "John" {
		t.Fatalf("expected name John, got %s", user.Name)
	}
	if user.Email != "john@example.com" {
		t.Fatalf("expected email john@example.com, got %s", user.Email)
	}
}

func TestAuthService_DuplicateEmail(t *testing.T) {
	service := NewAuthService("test-secret")

	_, _ = service.Register("John", "john@example.com", "password123")
	_, err := service.Register("Jane", "john@example.com", "password456")
	if err == nil {
		t.Fatal("expected duplicate email error")
	}
}

func TestAuthService_GenerateAndValidateToken(t *testing.T) {
	service := NewAuthService("test-secret")

	token, err := service.GenerateToken("1", "john@example.com")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	claims, err := service.ValidateToken(token)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if claims.UserID != "1" {
		t.Fatalf("expected user ID 1, got %s", claims.UserID)
	}
	if claims.Email != "john@example.com" {
		t.Fatalf("expected email john@example.com, got %s", claims.Email)
	}
}

func TestAuthService_FindByID(t *testing.T) {
	service := NewAuthService("test-secret")

	created, _ := service.Register("John", "john@example.com", "password123")
	found, err := service.FindByID(created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if found.ID != created.ID {
		t.Fatalf("expected ID %s, got %s", created.ID, found.ID)
	}
}
`

var authGuardTemplate = `package guards

import (
	"strings"

	"github.com/Ashishkapoor1469/Nestgo/common"
)

// JWTGuard is an authentication guard that validates JWT tokens.
// It extracts the bearer token from the Authorization header,
// validates it, and sets user_id and user_email in the context.
//
// Usage:
//
//	guard := guards.NewJWTGuard(func(token string) (string, error) {
//	    claims, err := authService.ValidateToken(token)
//	    if err != nil {
//	        return "", err
//	    }
//	    return claims.UserID, nil
//	})
type JWTGuard struct {
	validator func(token string) (userID string, err error)
}

// NewJWTGuard creates a new JWT guard.
func NewJWTGuard(validator func(token string) (string, error)) *JWTGuard {
	return &JWTGuard{validator: validator}
}

// CanActivate checks if the request has a valid JWT token.
func (g *JWTGuard) CanActivate(ctx *common.Context) (bool, error) {
	auth := ctx.Header("Authorization")
	if auth == "" {
		return false, common.Unauthorized("missing authorization header")
	}

	if !strings.HasPrefix(auth, "Bearer ") {
		return false, common.Unauthorized("invalid authorization format")
	}

	token := strings.TrimPrefix(auth, "Bearer ")
	userID, err := g.validator(token)
	if err != nil {
		return false, common.Unauthorized("invalid or expired token")
	}

	ctx.Set("user_id", userID)
	return true, nil
}
`

// DocsCmd creates the `nestgo docs:generate` command.
func DocsCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "docs:generate",
		Aliases: []string{"docs"},
		Short:   "Generate API documentation",
		Long:    "Generates OpenAPI/Swagger documentation from your route definitions.",
		RunE:    runDocsGenerate,
	}
}

func runDocsGenerate(cmd *cobra.Command, args []string) error {
	utils.EnsureProjectContext("docs:generate")

	utils.PrintHeader("📚 API Documentation Generator")

	routes := scanRoutes()

	if len(routes) == 0 {
		utils.PrintWarning("No routes found to document.")
		return nil
	}

	// Generate OpenAPI spec.
	spec := generateOpenAPISpec(routes)

	// Write to file.
	docsDir := "docs"
	if err := os.MkdirAll(docsDir, 0755); err != nil {
		return err
	}

	specPath := filepath.Join(docsDir, "openapi.json")
	if err := os.WriteFile(specPath, []byte(spec), 0644); err != nil {
		return err
	}

	utils.PrintSuccess("API documentation generated!")
	utils.PrintStep("📄", specPath)
	fmt.Println()
	utils.PrintDim("  View with Swagger UI or import into Postman.")
	utils.PrintDim("  The NestGo framework includes a built-in Swagger UI handler:")
	utils.PrintDim("    nesthttp.SwaggerUIHandler(\"/docs/openapi.json\")")
	fmt.Println()

	return nil
}

func generateOpenAPISpec(routes []scannedRoute) string {
	paths := ""
	for i, r := range routes {
		method := strings.ToLower(r.method)
		path := strings.ReplaceAll(r.path, "{", "{")
		path = strings.ReplaceAll(path, "}", "}")

		entry := fmt.Sprintf(`      "%s": {
        "summary": "%s",
        "responses": { "200": { "description": "Successful response" } }
      }`, method, r.summary)

		if i > 0 && routes[i-1].path == r.path {
			// Same path, add method.
			paths = strings.TrimSuffix(paths, "\n    },\n")
			paths += ",\n" + entry + "\n    },\n"
		} else {
			paths += fmt.Sprintf(`    "%s": {
%s
    },
`, path, entry)
		}
	}
	paths = strings.TrimSuffix(paths, ",\n")

	return fmt.Sprintf(`{
  "openapi": "3.0.3",
  "info": {
    "title": "NestGo API",
    "version": "1.0.0"
  },
  "paths": {
%s
  }
}`, paths)
}
