package commands

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/Ashishkapoor1469/Nestgo/cli/utils"
	"github.com/spf13/cobra"
)

// RoutesCmd creates the `nestgo routes` command.
func RoutesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "routes",
		Short: "Display all registered API routes",
		Long:  "Scans the project to discover and display all registered HTTP routes.",
		RunE:  runRoutes,
	}
	cmd.Flags().BoolP("verbose", "v", false, "Show controller and handler info")
	return cmd
}

func runRoutes(cmd *cobra.Command, args []string) error {
	utils.EnsureProjectContext("routes")
	verbose, _ := cmd.Flags().GetBool("verbose")

	utils.PrintHeader("🌐 API Route Explorer")

	routes := scanRoutes()

	if len(routes) == 0 {
		utils.PrintWarning("No routes found.")
		utils.PrintDim("Make sure your controllers implement Prefix() and Routes().")
		return nil
	}

	if verbose {
		headers := []string{"METHOD", "PATH", "CONTROLLER", "HANDLER", "SUMMARY"}
		var rows [][]string
		for _, r := range routes {
			methodStyled := colorMethod(r.method)
			rows = append(rows, []string{methodStyled, r.path, r.controller, r.handler, r.summary})
		}
		utils.PrintTable(headers, rows)
	} else {
		headers := []string{"METHOD", "PATH", "SUMMARY"}
		var rows [][]string
		for _, r := range routes {
			methodStyled := colorMethod(r.method)
			rows = append(rows, []string{methodStyled, r.path, r.summary})
		}
		utils.PrintTable(headers, rows)
	}

	fmt.Println()
	utils.PrintDim(fmt.Sprintf("  Total: %d routes", len(routes)))
	fmt.Println()

	return nil
}

type scannedRoute struct {
	method     string
	path       string
	controller string
	handler    string
	summary    string
}

func colorMethod(method string) string {
	switch strings.ToUpper(method) {
	case "GET":
		return utils.StyleSuccess.Render("GET")
	case "POST":
		return utils.StyleInfo.Render("POST")
	case "PUT":
		return utils.StyleWarning.Render("PUT")
	case "PATCH":
		return utils.StyleWarning.Render("PATCH")
	case "DELETE":
		return utils.StyleError.Render("DELETE")
	default:
		return method
	}
}

func scanRoutes() []scannedRoute {
	var routes []scannedRoute

	searchDir := filepath.Join("internal", "modules")
	if _, err := os.Stat(searchDir); os.IsNotExist(err) {
		searchDir = "."
	}

	_ = filepath.Walk(searchDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || filepath.Ext(path) != ".go" {
			return nil
		}

		// Skip vendor and node_modules blocks
		if strings.Contains(path, "vendor"+string(os.PathSeparator)) || strings.Contains(path, "node_modules"+string(os.PathSeparator)) {
			return nil
		}

		fset := token.NewFileSet()
		node, parseErr := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if parseErr != nil {
			return nil
		}

		// Detect the controller struct and prefix.
		controllerName := ""
		prefix := ""

		for _, decl := range node.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv == nil {
				continue
			}

			// Find Prefix() method.
			if fn.Name.Name == "Prefix" {
				controllerName = extractReceiverType(fn)
				// Extract return value.
				prefix = extractStringReturn(fn)
			}
		}

		if controllerName == "" || prefix == "" {
			return nil
		}

		// Find Routes() method.
		for _, decl := range node.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv == nil || fn.Name.Name != "Routes" {
				continue
			}

			// Parse route definitions from the return statement.
			parsedRoutes := extractRouteDefinitions(fn, prefix, controllerName)
			routes = append(routes, parsedRoutes...)
		}

		return nil
	})

	return routes
}

func extractReceiverType(fn *ast.FuncDecl) string {
	if fn.Recv == nil || len(fn.Recv.List) == 0 {
		return ""
	}
	t := fn.Recv.List[0].Type
	// Handle *Type
	if star, ok := t.(*ast.StarExpr); ok {
		if ident, ok := star.X.(*ast.Ident); ok {
			return ident.Name
		}
	}
	if ident, ok := t.(*ast.Ident); ok {
		return ident.Name
	}
	return ""
}

func extractStringReturn(fn *ast.FuncDecl) string {
	if fn.Body == nil {
		return ""
	}
	for _, stmt := range fn.Body.List {
		ret, ok := stmt.(*ast.ReturnStmt)
		if !ok || len(ret.Results) == 0 {
			continue
		}
		lit, ok := ret.Results[0].(*ast.BasicLit)
		if ok {
			return strings.Trim(lit.Value, `"`)
		}
	}
	return ""
}

func extractRouteDefinitions(fn *ast.FuncDecl, prefix, controller string) []scannedRoute {
	var routes []scannedRoute

	if fn.Body == nil {
		return routes
	}

	// Walk the AST looking for composite literals that look like Route structs.
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		comp, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}

		// Check if it has Method, Path, Handler fields.
		method := ""
		path := ""
		handler := ""
		summary := ""

		for _, elt := range comp.Elts {
			kv, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}

			key, ok := kv.Key.(*ast.Ident)
			if !ok {
				continue
			}

			switch key.Name {
			case "Method":
				method = extractSelectorOrLiteral(kv.Value)
			case "Path":
				if lit, ok := kv.Value.(*ast.BasicLit); ok {
					path = strings.Trim(lit.Value, `"`)
				}
			case "Handler":
				if sel, ok := kv.Value.(*ast.SelectorExpr); ok {
					handler = sel.Sel.Name
				}
			case "Summary":
				if lit, ok := kv.Value.(*ast.BasicLit); ok {
					summary = strings.Trim(lit.Value, `"`)
				}
			}
		}

		if method != "" && path != "" {
			fullPath := prefix
			if path != "/" {
				fullPath = prefix + path
			}

			routes = append(routes, scannedRoute{
				method:     method,
				path:       fullPath,
				controller: controller,
				handler:    handler,
				summary:    summary,
			})
		}

		return true
	})

	return routes
}

func extractSelectorOrLiteral(expr ast.Expr) string {
	// Handle http.MethodGet style
	if sel, ok := expr.(*ast.SelectorExpr); ok {
		switch sel.Sel.Name {
		case "MethodGet":
			return "GET"
		case "MethodPost":
			return "POST"
		case "MethodPut":
			return "PUT"
		case "MethodPatch":
			return "PATCH"
		case "MethodDelete":
			return "DELETE"
		case "MethodOptions":
			return "OPTIONS"
		case "MethodHead":
			return "HEAD"
		default:
			return sel.Sel.Name
		}
	}
	// Handle string literal
	if lit, ok := expr.(*ast.BasicLit); ok {
		return strings.Trim(lit.Value, `"`)
	}
	return ""
}

// LintArchCmd creates the `nestgo lint-arch` command.
func LintArchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "lint-arch",
		Short: "Check module isolation and architecture rules",
		Long:  "Enforces clean architecture boundaries by detecting illegal cross-module imports.",
		RunE:  runLintArch,
	}
}

func runLintArch(cmd *cobra.Command, args []string) error {
	utils.EnsureProjectContext("lint-arch")

	utils.PrintHeader("🏗️  Architecture Linter")

	modulesDir := filepath.Join("internal", "modules")
	entries, err := os.ReadDir(modulesDir)
	if err != nil {
		utils.PrintWarning("No modules directory found.")
		return nil
	}

	violations := 0
	checked := 0

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		modName := entry.Name()
		modPath := filepath.Join(modulesDir, modName)

		_ = filepath.Walk(modPath, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() || filepath.Ext(path) != ".go" {
				return nil
			}

			checked++
			fset := token.NewFileSet()
			node, parseErr := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
			if parseErr != nil {
				return nil
			}

			for _, imp := range node.Imports {
				importPath := strings.Trim(imp.Path.Value, `"`)

				// Check for cross-module imports (modules importing from other modules).
				if strings.Contains(importPath, "internal/modules/") {
					importedModule := extractModuleName(importPath)
					if importedModule != "" && importedModule != modName {
						violations++
						utils.PrintError(fmt.Sprintf(
							"%s imports from module '%s' (cross-module violation)",
							path, importedModule,
						))
						utils.PrintDim(fmt.Sprintf("    Import: %s", importPath))
						utils.PrintDim("    Fix: Use dependency injection or shared interfaces in internal/common/")
					}
				}
			}

			return nil
		})
	}

	fmt.Println()
	if violations == 0 {
		utils.PrintSuccess(fmt.Sprintf("Architecture check passed! (%d files checked)", checked))
	} else {
		utils.PrintError(fmt.Sprintf("%d violation(s) found in %d files", violations, checked))
		utils.PrintDim("  Modules should communicate through DI providers, not direct imports.")
	}
	fmt.Println()

	return nil
}

func extractModuleName(importPath string) string {
	idx := strings.Index(importPath, "internal/modules/")
	if idx < 0 {
		return ""
	}
	rest := importPath[idx+len("internal/modules/"):]
	parts := strings.SplitN(rest, "/", 2)
	return parts[0]
}
