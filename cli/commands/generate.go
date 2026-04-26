package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
	"unicode"

	"github.com/Ashishkapoor1469/Nestgo/cli/utils"
	"github.com/spf13/cobra"
)

// GenerateCmd creates the `nestgo generate` command group.
func GenerateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate [type] [name]",
		Aliases: []string{"g"},
		Short:   "Generate NestGo components",
		Long:    "Generate modules, controllers, services, guards, middleware, interceptors, schemas, DTOs, tests, and full resources.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return runInteractiveGenerate()
			}
			return cmd.Help()
		},
	}

	cmd.AddCommand(
		generateModuleCmd(),
		generateControllerCmd(),
		generateServiceCmd(),
		generateMiddlewareCmd(),
		generateGuardCmd(),
		generateInterceptorCmd(),
		generateResourceCmd(),
		generateSchemaCmd(),
		generateDTOCmd(),
		generateTestCmd(),
		AuthCmd(),
	)

	return cmd
}

func generateModuleCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "module [name]",
		Short: "Generate a new module",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return generateComponent("module", args[0])
		},
	}
}

func generateControllerCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "controller [name]",
		Short: "Generate a new controller",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return generateComponent("controller", args[0])
		},
	}
}

func generateServiceCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "service [name]",
		Short: "Generate a new service",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return generateComponent("service", args[0])
		},
	}
}

func generateMiddlewareCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "middleware [name]",
		Short: "Generate a new middleware",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return generateComponent("middleware", args[0])
		},
	}
}

func generateGuardCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "guard [name]",
		Short: "Generate a new guard",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return generateComponent("guard", args[0])
		},
	}
}

func generateInterceptorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "interceptor [name]",
		Short: "Generate a new interceptor",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return generateComponent("interceptor", args[0])
		},
	}
}

func generateResourceCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "resource [name]",
		Short: "Generate a full CRUD resource (module + controller + service + DTOs + tests)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return generateResource(args[0])
		},
	}
}

// generateComponent generates a single component file.
func generateComponent(componentType, name string) error {
	utils.EnsureProjectContext("generate " + componentType)

	name = strings.ToLower(name)
	pascal := toPascalCase(name)

	dir := filepath.Join("internal", "modules", name)
	if componentType == "middleware" {
		dir = filepath.Join("internal", "common", "middleware")
	} else if componentType == "guard" {
		dir = filepath.Join("internal", "common", "guards")
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data := map[string]string{
		"Name":       name,
		"PascalName": pascal,
		"Package":    name,
	}

	if componentType == "middleware" || componentType == "guard" {
		data["Package"] = componentType
		if componentType == "guard" {
			data["Package"] = "guards"
		}
	}

	var tmplStr string
	var fileName string

	switch componentType {
	case "module":
		tmplStr = moduleTemplate
		fileName = "module.go"
	case "controller":
		tmplStr = controllerTemplate
		fileName = "controller.go"
	case "service":
		tmplStr = serviceTemplate
		fileName = "service.go"
	case "middleware":
		tmplStr = middlewareGenTemplate
		fileName = name + "_middleware.go"
	case "guard":
		tmplStr = guardGenTemplate
		fileName = name + "_guard.go"
	case "interceptor":
		tmplStr = interceptorGenTemplate
		fileName = name + "_interceptor.go"
	default:
		return fmt.Errorf("unknown component type: %s", componentType)
	}

	filePath := filepath.Join(dir, fileName)
	if err := writeGenTemplate(filePath, tmplStr, data); err != nil {
		return err
	}

	fmt.Printf("✅ Generated %s: %s\n", componentType, filePath)
	return nil
}

func runInteractiveGenerate() error {
	utils.PrintHeader("🛠️ Generate Components")
	fmt.Println("  What would you like to generate?")
	fmt.Println("    1. resource     (Full CRUD component)")
	fmt.Println("    2. module       (Module definition)")
	fmt.Println("    3. controller   (HTTP Controller)")
	fmt.Println("    4. service      (Business logic)")
	fmt.Println("    5. auth         (Authentication starter)")
	fmt.Println("    6. schema       (Validation schema)")
	fmt.Println("    7. dto          (Data Transfer Object)")
	fmt.Println("    8. test         (Test scaffolding)")
	fmt.Println()
	
	fmt.Print(utils.StyleAccent.Render("  > Choose an option (1-8): "))
	
	var option int
	_, err := fmt.Scanf("%d", &option)
	if err != nil {
		return fmt.Errorf("invalid input")
	}

	var name string
	if option != 5 { // Auth doesn't need a name
		fmt.Print(utils.StyleAccent.Render("  > Enter component name: "))
		_, err = fmt.Scanf("%s", &name)
		if err != nil || name == "" {
			return fmt.Errorf("name is required")
		}
	}

	switch option {
	case 1:
		return generateResource(name)
	case 2:
		return generateComponent("module", name)
	case 3:
		return generateComponent("controller", name)
	case 4:
		return generateComponent("service", name)
	case 5:
		return runGenerateAuth(nil, nil)
	case 6:
		return generateSchema(name)
	case 7:
		return generateDTO(name)
	case 8:
		return generateTest(name)
	default:
		return fmt.Errorf("invalid option selected")
	}
}

// generateResource generates a complete CRUD resource.
func generateResource(name string) error {
	utils.EnsureProjectContext("generate resource")

	name = strings.ToLower(name)
	pascal := toPascalCase(name)
	singular := strings.TrimSuffix(name, "s")
	singularPascal := toPascalCase(singular)

	dir := filepath.Join("internal", "modules", name)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	fmt.Printf("\n  ✦ Scaffolding resource: %s\n\n", name)

	data := map[string]string{
		"Name":           name,
		"PascalName":     pascal,
		"Singular":       singular,
		"SingularPascal": singularPascal,
		"Package":        name,
	}

	moduleFiles := map[string]string{
		"module.go":     resourceModuleTemplate,
		"controller.go": resourceControllerTemplate,
		"service.go":    resourceServiceTemplate,
		"dto.go":        resourceDTOTemplate,
		"entity.go":     resourceEntityTemplate,
		name + "_test.go": resourceTestTemplate,
	}

	for fileName, tmplStr := range moduleFiles {
		filePath := filepath.Join(dir, fileName)
		if err := writeGenTemplate(filePath, tmplStr, data); err != nil {
			return fmt.Errorf("failed to generate %s: %w", fileName, err)
		}
		fmt.Printf("  ├── %s\n", filePath)
	}

	// Generate migration file
	if err := os.MkdirAll("migrations", 0755); err == nil {
		stamp := fmt.Sprintf("%d", time.Now().Unix())
		migFile := filepath.Join("migrations", stamp+"_create_"+name+".sql")
		if err := writeGenTemplate(migFile, resourceMigrationTemplate, data); err == nil {
			fmt.Printf("  ├── %s\n", migFile)
		}
	}

	fmt.Printf("\n  ✅ Resource '%s' generated successfully!\n", name)
	fmt.Println()
	fmt.Println("  Next steps:")
	fmt.Printf("    1. Register in app_module.go: Imports: []common.Module{&%s.%sModule{}}\n", name, pascal)
	fmt.Printf("    2. Run migrations:            nestgo migration:run\n")
	fmt.Printf("    3. Test your endpoints:       curl http://localhost:3000/api/%s\n", name)
	fmt.Println()
	return nil
}

// generateSchema generates a schema file for a module.
func generateSchema(name string) error {
	utils.EnsureProjectContext("generate schema")

	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, "-", "_")
	pascal := toPascalCase(name)

	// Determine module name (strip create_/update_ prefix if present).
	moduleName := name
	for _, prefix := range []string{"create_", "update_", "delete_"} {
		moduleName = strings.TrimPrefix(moduleName, prefix)
	}

	dir := filepath.Join("internal", "modules", moduleName, "schemas")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data := map[string]string{
		"Name":       name,
		"PascalName": pascal,
		"Module":     moduleName,
		"Package":    "schemas",
	}

	fileName := name + ".schema.go"
	filePath := filepath.Join(dir, fileName)

	if err := writeGenTemplate(filePath, schemaGenTemplate, data); err != nil {
		return err
	}

	utils.PrintSuccess(fmt.Sprintf("Schema generated: %s", filePath))
	return nil
}

// generateDTO generates a DTO file for a module.
func generateDTO(name string) error {
	utils.EnsureProjectContext("generate dto")

	name = strings.ToLower(name)
	pascal := toPascalCase(name)
	singular := strings.TrimSuffix(name, "s")
	singularPascal := toPascalCase(singular)

	dir := filepath.Join("internal", "modules", name)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data := map[string]string{
		"Name":           name,
		"PascalName":     pascal,
		"Singular":       singular,
		"SingularPascal": singularPascal,
		"Package":        name,
	}

	filePath := filepath.Join(dir, "dto.go")
	if err := writeGenTemplate(filePath, dtoGenTemplate, data); err != nil {
		return err
	}

	utils.PrintSuccess(fmt.Sprintf("DTO generated: %s", filePath))
	return nil
}

// generateTest generates a test file for a module.
func generateTest(name string) error {
	utils.EnsureProjectContext("generate test")

	name = strings.ToLower(name)
	pascal := toPascalCase(name)
	singular := strings.TrimSuffix(name, "s")
	singularPascal := toPascalCase(singular)

	dir := filepath.Join("internal", "modules", name)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data := map[string]string{
		"Name":           name,
		"PascalName":     pascal,
		"Singular":       singular,
		"SingularPascal": singularPascal,
		"Package":        name,
	}

	filePath := filepath.Join(dir, name+"_test.go")
	if err := writeGenTemplate(filePath, testGenTemplate, data); err != nil {
		return err
	}

	utils.PrintSuccess(fmt.Sprintf("Test file generated: %s", filePath))
	return nil
}

func generateSchemaCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "schema [name]",
		Short: "Generate a validation schema",
		Long: `Generate a validation schema with struct tags for a module.

Example:
  nestgo generate schema create-user
  nestgo generate schema update-product`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return generateSchema(args[0])
		},
	}
}

func generateDTOCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "dto [name]",
		Short: "Generate DTOs (Create + Update)",
		Long: `Generate Data Transfer Object structs with validation for a module.

Example:
  nestgo generate dto users
  nestgo generate dto products`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return generateDTO(args[0])
		},
	}
}

func generateTestCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "test [name]",
		Short: "Generate a test file",
		Long: `Generate a table-driven test file for a module.

Example:
  nestgo generate test users`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return generateTest(args[0])
		},
	}
}

func writeGenTemplate(path, tmplStr string, data any) error {
	funcMap := template.FuncMap{
		"lower": strings.ToLower,
		"upper": strings.ToUpper,
	}

	tmpl, err := template.New("gen").Funcs(funcMap).Parse(tmplStr)
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, data)
}

func toPascalCase(s string) string {
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || r == ' '
	})
	var result strings.Builder
	for _, part := range parts {
		if len(part) > 0 {
			runes := []rune(part)
			runes[0] = unicode.ToUpper(runes[0])
			result.WriteString(string(runes))
		}
	}
	if result.Len() == 0 {
		runes := []rune(s)
		if len(runes) > 0 {
			runes[0] = unicode.ToUpper(runes[0])
		}
		return string(runes)
	}
	return result.String()
}

// --- Component Templates ---

var moduleTemplate = `package {{.Package}}

import (
	"github.com/Ashishkapoor1469/Nestgo/di"
	"github.com/Ashishkapoor1469/Nestgo/common"
)

// {{.PascalName}}Module defines the {{.Name}} feature module.
type {{.PascalName}}Module struct{}

func (m *{{.PascalName}}Module) Module() common.ModuleConfig {
	return common.ModuleConfig{
		Name:        "{{.Name}}",
		Imports:     []common.Module{},
		Controllers: []common.Controller{},
		Providers:   []di.Provider{},
	}
}
`

var controllerTemplate = `package {{.Package}}

import (
	"net/http"

	"github.com/Ashishkapoor1469/Nestgo/common"
)

// {{.PascalName}}Controller handles HTTP requests for {{.Name}}.
type {{.PascalName}}Controller struct {
	service *{{.PascalName}}Service
}

// New{{.PascalName}}Controller creates a new controller.
func New{{.PascalName}}Controller(service *{{.PascalName}}Service) *{{.PascalName}}Controller {
	return &{{.PascalName}}Controller{service: service}
}

// Prefix returns the route prefix.
func (c *{{.PascalName}}Controller) Prefix() string {
	return "/{{.Name}}"
}

// Routes returns the controller's routes.
func (c *{{.PascalName}}Controller) Routes() []common.Route {
	return []common.Route{
		{Method: http.MethodGet, Path: "/", Handler: c.FindAll, Summary: "List all {{.Name}}"},
		{Method: http.MethodGet, Path: "/{id}", Handler: c.FindOne, Summary: "Get {{.Name}} by ID"},
		{Method: http.MethodPost, Path: "/", Handler: c.Create, Summary: "Create {{.Name}}"},
		{Method: http.MethodPut, Path: "/{id}", Handler: c.Update, Summary: "Update {{.Name}}"},
		{Method: http.MethodDelete, Path: "/{id}", Handler: c.Delete, Summary: "Delete {{.Name}}"},
	}
}

func (c *{{.PascalName}}Controller) FindAll(ctx *common.Context) error {
	// TODO: Implement
	return ctx.OK(map[string]string{"message": "list {{.Name}}"})
}

func (c *{{.PascalName}}Controller) FindOne(ctx *common.Context) error {
	id := ctx.Param("id")
	return ctx.OK(map[string]string{"id": id})
}

func (c *{{.PascalName}}Controller) Create(ctx *common.Context) error {
	// TODO: Implement
	return ctx.Created(map[string]string{"message": "created"})
}

func (c *{{.PascalName}}Controller) Update(ctx *common.Context) error {
	id := ctx.Param("id")
	return ctx.OK(map[string]string{"id": id, "message": "updated"})
}

func (c *{{.PascalName}}Controller) Delete(ctx *common.Context) error {
	// TODO: Implement
	return ctx.NoContent()
}
`

var serviceTemplate = `package {{.Package}}

// {{.PascalName}}Service provides business logic for {{.Name}}.
type {{.PascalName}}Service struct{}

// New{{.PascalName}}Service creates a new service.
func New{{.PascalName}}Service() *{{.PascalName}}Service {
	return &{{.PascalName}}Service{}
}
`

var middlewareGenTemplate = `package middleware

import (
	"net/http"
)

// {{.PascalName}}Middleware is a custom middleware.
func {{.PascalName}}Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Add middleware logic here
		next.ServeHTTP(w, r)
	})
}
`

var guardGenTemplate = `package guards

import (
	"github.com/Ashishkapoor1469/Nestgo/common"
)

// {{.PascalName}}Guard implements authorization.
type {{.PascalName}}Guard struct{}

// New{{.PascalName}}Guard creates a new guard.
func New{{.PascalName}}Guard() *{{.PascalName}}Guard {
	return &{{.PascalName}}Guard{}
}

// CanActivate checks if the request is authorized.
func (g *{{.PascalName}}Guard) CanActivate(ctx *common.Context) (bool, error) {
	// TODO: Implement authorization logic
	return true, nil
}
`

var interceptorGenTemplate = `package {{.Package}}

import (
	"github.com/Ashishkapoor1469/Nestgo/common"
)

// {{.PascalName}}Interceptor is a custom interceptor.
type {{.PascalName}}Interceptor struct{}

// New{{.PascalName}}Interceptor creates a new interceptor.
func New{{.PascalName}}Interceptor() *{{.PascalName}}Interceptor {
	return &{{.PascalName}}Interceptor{}
}

// Intercept processes the request/response.
func (i *{{.PascalName}}Interceptor) Intercept(ctx *common.Context, next common.HandlerFunc) error {
	// Before handler
	err := next(ctx)
	// After handler
	return err
}
`

// --- Resource Templates ---

var resourceModuleTemplate = `package {{.Package}}

import (
	"github.com/Ashishkapoor1469/Nestgo/di"
	"github.com/Ashishkapoor1469/Nestgo/common"
)

// {{.PascalName}}Module defines the {{.Name}} feature module.
type {{.PascalName}}Module struct{}

func (m *{{.PascalName}}Module) Module() common.ModuleConfig {
	service := New{{.PascalName}}Service()
	controller := New{{.PascalName}}Controller(service)

	return common.ModuleConfig{
		Name: "{{.Name}}",
		Controllers: []common.Controller{controller},
		Providers: []di.Provider{
			{Instance: service},
		},
	}
}
`

var resourceControllerTemplate = `package {{.Package}}

import (
	"net/http"

	"github.com/Ashishkapoor1469/Nestgo/common"
)

// {{.PascalName}}Controller handles HTTP requests for {{.Name}}.
type {{.PascalName}}Controller struct {
	service *{{.PascalName}}Service
}

// New{{.PascalName}}Controller creates a new controller.
func New{{.PascalName}}Controller(service *{{.PascalName}}Service) *{{.PascalName}}Controller {
	return &{{.PascalName}}Controller{service: service}
}

// Prefix returns the route prefix.
func (c *{{.PascalName}}Controller) Prefix() string {
	return "/{{.Name}}"
}

// Routes returns the controller's route definitions.
func (c *{{.PascalName}}Controller) Routes() []common.Route {
	return []common.Route{
		{Method: http.MethodGet, Path: "/", Handler: c.FindAll, Summary: "List all {{.Name}}"},
		{Method: http.MethodGet, Path: "/{id}", Handler: c.FindOne, Summary: "Get a {{.Singular}} by ID"},
		{Method: http.MethodPost, Path: "/", Handler: c.Create, Summary: "Create a new {{.Singular}}"},
		{Method: http.MethodPut, Path: "/{id}", Handler: c.Update, Summary: "Update a {{.Singular}}"},
		{Method: http.MethodDelete, Path: "/{id}", Handler: c.Delete, Summary: "Delete a {{.Singular}}"},
	}
}

// FindAll returns all {{.Name}}.
func (c *{{.PascalName}}Controller) FindAll(ctx *common.Context) error {
	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 20)

	items, total, err := c.service.FindAll(page, limit)
	if err != nil {
		return common.InternalError(err.Error())
	}

	return ctx.Paginated(items, total, page, limit)
}

// FindOne returns a single {{.Singular}}.
func (c *{{.PascalName}}Controller) FindOne(ctx *common.Context) error {
	id := ctx.Param("id")

	item, err := c.service.FindByID(id)
	if err != nil {
		return common.NotFound("{{.Singular}} not found")
	}

	return ctx.OK(item)
}

// Create creates a new {{.Singular}}.
func (c *{{.PascalName}}Controller) Create(ctx *common.Context) error {
	var dto Create{{.SingularPascal}}DTO
	if err := ctx.Bind(&dto); err != nil {
		return common.BadRequest(err.Error())
	}

	item, err := c.service.Create(dto)
	if err != nil {
		return common.InternalError(err.Error())
	}

	return ctx.Created(item)
}

// Update updates a {{.Singular}}.
func (c *{{.PascalName}}Controller) Update(ctx *common.Context) error {
	id := ctx.Param("id")

	var dto Update{{.SingularPascal}}DTO
	if err := ctx.Bind(&dto); err != nil {
		return common.BadRequest(err.Error())
	}

	item, err := c.service.Update(id, dto)
	if err != nil {
		return common.NotFound("{{.Singular}} not found")
	}

	return ctx.OK(item)
}

// Delete removes a {{.Singular}}.
func (c *{{.PascalName}}Controller) Delete(ctx *common.Context) error {
	id := ctx.Param("id")

	if err := c.service.Delete(id); err != nil {
		return common.NotFound("{{.Singular}} not found")
	}

	return ctx.NoContent()
}
`

var resourceServiceTemplate = `package {{.Package}}

import (
	"fmt"
	"sync"
	"time"
)

// {{.PascalName}}Service provides business logic for {{.Name}}.
type {{.PascalName}}Service struct {
	mu    sync.RWMutex
	items map[string]*{{.SingularPascal}}
	seq   int
}

// New{{.PascalName}}Service creates a new service.
func New{{.PascalName}}Service() *{{.PascalName}}Service {
	return &{{.PascalName}}Service{
		items: make(map[string]*{{.SingularPascal}}),
	}
}

// FindAll returns paginated {{.Name}}.
func (s *{{.PascalName}}Service) FindAll(page, limit int) ([]*{{.SingularPascal}}, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	all := make([]*{{.SingularPascal}}, 0, len(s.items))
	for _, item := range s.items {
		all = append(all, item)
	}

	total := len(all)
	start := (page - 1) * limit
	if start >= total {
		return []*{{.SingularPascal}}{}, total, nil
	}

	end := start + limit
	if end > total {
		end = total
	}

	return all[start:end], total, nil
}

// FindByID returns a {{.Singular}} by ID.
func (s *{{.PascalName}}Service) FindByID(id string) (*{{.SingularPascal}}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.items[id]
	if !ok {
		return nil, fmt.Errorf("{{.Singular}} not found: %s", id)
	}
	return item, nil
}

// Create creates a new {{.Singular}}.
func (s *{{.PascalName}}Service) Create(dto Create{{.SingularPascal}}DTO) (*{{.SingularPascal}}, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.seq++
	item := &{{.SingularPascal}}{
		ID:        fmt.Sprintf("%d", s.seq),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.items[item.ID] = item
	return item, nil
}

// Update updates a {{.Singular}}.
func (s *{{.PascalName}}Service) Update(id string, dto Update{{.SingularPascal}}DTO) (*{{.SingularPascal}}, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.items[id]
	if !ok {
		return nil, fmt.Errorf("{{.Singular}} not found: %s", id)
	}

	item.UpdatedAt = time.Now()
	return item, nil
}

// Delete removes a {{.Singular}}.
func (s *{{.PascalName}}Service) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.items[id]; !ok {
		return fmt.Errorf("{{.Singular}} not found: %s", id)
	}

	delete(s.items, id)
	return nil
}
`

var resourceDTOTemplate = `package {{.Package}}

import (
	"github.com/Ashishkapoor1469/Nestgo/common"
)

// Create{{.SingularPascal}}DTO is the request body for creating a {{.Singular}}.
type Create{{.SingularPascal}}DTO struct {
	Name        string ` + "`" + `json:"name"` + "`" + `
	Description string ` + "`" + `json:"description,omitempty"` + "`" + `
}

// Validate validates the DTO.
func (d *Create{{.SingularPascal}}DTO) Validate() error {
	return common.NewValidator().
		Required("name", d.Name).
		MinLength("name", d.Name, 2).
		MaxLength("name", d.Name, 100).
		Validate()
}

// Update{{.SingularPascal}}DTO is the request body for updating a {{.Singular}}.
type Update{{.SingularPascal}}DTO struct {
	Name        *string ` + "`" + `json:"name,omitempty"` + "`" + `
	Description *string ` + "`" + `json:"description,omitempty"` + "`" + `
}

// Validate validates the DTO.
func (d *Update{{.SingularPascal}}DTO) Validate() error {
	v := common.NewValidator()
	if d.Name != nil {
		v.MinLength("name", *d.Name, 2).
			MaxLength("name", *d.Name, 100)
	}
	return v.Validate()
}
`

var resourceEntityTemplate = `package {{.Package}}

import "time"

// {{.SingularPascal}} is the {{.Singular}} entity.
type {{.SingularPascal}} struct {
	ID          string    ` + "`" + `json:"id"` + "`" + `
	Name        string    ` + "`" + `json:"name"` + "`" + `
	Description string    ` + "`" + `json:"description,omitempty"` + "`" + `
	CreatedAt   time.Time ` + "`" + `json:"createdAt"` + "`" + `
	UpdatedAt   time.Time ` + "`" + `json:"updatedAt"` + "`" + `
}
`

var resourceTestTemplate = `package {{.Package}}

import (
	"testing"
)

func TestNew{{.PascalName}}Service(t *testing.T) {
	service := New{{.PascalName}}Service()
	if service == nil {
		t.Fatal("expected service to be non-nil")
	}
}

func Test{{.PascalName}}Service_Create(t *testing.T) {
	service := New{{.PascalName}}Service()

	dto := Create{{.SingularPascal}}DTO{
		Name: "Test Item",
	}

	item, err := service.Create(dto)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if item.ID == "" {
		t.Fatal("expected ID to be set")
	}
}

func Test{{.PascalName}}Service_FindAll(t *testing.T) {
	service := New{{.PascalName}}Service()

	// Create test data.
	for i := 0; i < 5; i++ {
		_, _ = service.Create(Create{{.SingularPascal}}DTO{Name: "Item"})
	}

	items, total, err := service.FindAll(1, 10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 5 {
		t.Fatalf("expected 5 total, got %d", total)
	}
	if len(items) != 5 {
		t.Fatalf("expected 5 items, got %d", len(items))
	}
}

func Test{{.PascalName}}Service_FindByID(t *testing.T) {
	service := New{{.PascalName}}Service()

	created, _ := service.Create(Create{{.SingularPascal}}DTO{Name: "Test"})

	found, err := service.FindByID(created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if found.ID != created.ID {
		t.Fatalf("expected ID %s, got %s", created.ID, found.ID)
	}
}

func Test{{.PascalName}}Service_Delete(t *testing.T) {
	service := New{{.PascalName}}Service()

	created, _ := service.Create(Create{{.SingularPascal}}DTO{Name: "Test"})
	err := service.Delete(created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = service.FindByID(created.ID)
	if err == nil {
		t.Fatal("expected error after deletion")
	}
}
`

// --- Schema Template ---

var schemaGenTemplate = `package {{.Package}}

// {{.PascalName}}Schema defines the validation schema for {{.Name}}.
// Uses struct tags for automatic validation via common.ValidateStruct().
type {{.PascalName}}Schema struct {
	Name        string ` + "`" + `json:"name" validate:"required,min=2,max=100"` + "`" + `
	Email       string ` + "`" + `json:"email" validate:"required,email"` + "`" + `
	Description string ` + "`" + `json:"description,omitempty" validate:"max=500"` + "`" + `
}
`

// --- DTO Template ---

var dtoGenTemplate = `package {{.Package}}

// Create{{.SingularPascal}}DTO is the request body for creating a {{.Singular}}.
// Uses validate struct tags for automatic validation.
type Create{{.SingularPascal}}DTO struct {
	Name        string ` + "`" + `json:"name" validate:"required,min=2,max=100"` + "`" + `
	Description string ` + "`" + `json:"description,omitempty" validate:"max=500"` + "`" + `
}

// Update{{.SingularPascal}}DTO is the request body for updating a {{.Singular}}.
type Update{{.SingularPascal}}DTO struct {
	Name        *string ` + "`" + `json:"name,omitempty" validate:"min=2,max=100"` + "`" + `
	Description *string ` + "`" + `json:"description,omitempty" validate:"max=500"` + "`" + `
}
`

// --- Test Template ---

var testGenTemplate = `package {{.Package}}

import (
	"testing"
)

func TestNew{{.PascalName}}Service(t *testing.T) {
	service := New{{.PascalName}}Service()
	if service == nil {
		t.Fatal("expected service to be non-nil")
	}
}

func Test{{.PascalName}}Service_CRUD(t *testing.T) {
	service := New{{.PascalName}}Service()

	tests := []struct {
		name    string
		action  string
		wantErr bool
	}{
		{"create valid", "create", false},
		{"find all", "findAll", false},
		{"find by id", "findById", false},
		{"delete", "delete", false},
		{"find deleted", "findDeleted", true},
	}

	var createdID string

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.action {
			case "create":
				dto := Create{{.SingularPascal}}DTO{Name: "Test"}
				item, err := service.Create(dto)
				if (err != nil) != tt.wantErr {
					t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				}
				if item != nil {
					createdID = item.ID
				}
			case "findAll":
				items, total, err := service.FindAll(1, 10)
				if (err != nil) != tt.wantErr {
					t.Errorf("FindAll() error = %v, wantErr %v", err, tt.wantErr)
				}
				if total == 0 || len(items) == 0 {
					t.Error("expected items")
				}
			case "findById":
				_, err := service.FindByID(createdID)
				if (err != nil) != tt.wantErr {
					t.Errorf("FindByID() error = %v, wantErr %v", err, tt.wantErr)
				}
			case "delete":
				err := service.Delete(createdID)
				if (err != nil) != tt.wantErr {
					t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				}
			case "findDeleted":
				_, err := service.FindByID(createdID)
				if (err != nil) != tt.wantErr {
					t.Errorf("FindByID() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}
`

var resourceMigrationTemplate = `-- Migration: create_{{.Name}}
-- Generated by nestgo
-- Run: nestgo migration:run

CREATE TABLE IF NOT EXISTS {{.Name}} (
  id          VARCHAR(36)  PRIMARY KEY DEFAULT gen_random_uuid(),
  name        VARCHAR(255) NOT NULL,
  created_at  TIMESTAMP    NOT NULL DEFAULT NOW(),
  updated_at  TIMESTAMP    NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_{{.Name}}_created_at ON {{.Name}} (created_at DESC);

-- Rollback:
-- DROP TABLE IF EXISTS {{.Name}};
`
