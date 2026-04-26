package nesthttp

import (
	"encoding/json"
	"fmt"
	nethttp "net/http"
	"reflect"
	"strings"
)

// OpenAPISpec represents an OpenAPI 3.0 specification.
type OpenAPISpec struct {
	OpenAPI string              `json:"openapi"`
	Info    OpenAPIInfo         `json:"info"`
	Paths   map[string]PathItem `json:"paths"`
	Servers []OpenAPIServer     `json:"servers,omitempty"`
}

// OpenAPIInfo holds API metadata.
type OpenAPIInfo struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Version     string `json:"version"`
}

// OpenAPIServer represents a server.
type OpenAPIServer struct {
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
}

// PathItem represents operations on a single path.
type PathItem map[string]Operation

// Operation represents a single API operation.
type Operation struct {
	Summary     string                 `json:"summary,omitempty"`
	Description string                 `json:"description,omitempty"`
	OperationID string                 `json:"operationId,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	Parameters  []Parameter            `json:"parameters,omitempty"`
	RequestBody *RequestBody           `json:"requestBody,omitempty"`
	Responses   map[string]APIResponse `json:"responses"`
}

// Parameter represents an API parameter.
type Parameter struct {
	Name        string `json:"name"`
	In          string `json:"in"` // "path", "query", "header"
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required,omitempty"`
	Schema      Schema `json:"schema"`
}

// RequestBody represents a request body.
type RequestBody struct {
	Description string               `json:"description,omitempty"`
	Required    bool                 `json:"required,omitempty"`
	Content     map[string]MediaType `json:"content"`
}

// MediaType defines the schema for a media type.
type MediaType struct {
	Schema Schema `json:"schema"`
}

// APIResponse represents an API response in OpenAPI spec.
type APIResponse struct {
	Description string               `json:"description"`
	Content     map[string]MediaType `json:"content,omitempty"`
}

// Schema defines a data schema.
type Schema struct {
	Type       string            `json:"type,omitempty"`
	Format     string            `json:"format,omitempty"`
	Properties map[string]Schema `json:"properties,omitempty"`
	Items      *Schema           `json:"items,omitempty"`
	Required   []string          `json:"required,omitempty"`
	Example    any               `json:"example,omitempty"`
}

// OpenAPIGenerator generates OpenAPI specs from routes.
type OpenAPIGenerator struct {
	spec OpenAPISpec
}

// NewOpenAPIGenerator creates a new OpenAPI generator.
func NewOpenAPIGenerator(title, version string) *OpenAPIGenerator {
	return &OpenAPIGenerator{
		spec: OpenAPISpec{
			OpenAPI: "3.0.3",
			Info: OpenAPIInfo{
				Title:   title,
				Version: version,
			},
			Paths: make(map[string]PathItem),
		},
	}
}

// AddServer adds a server to the spec.
func (g *OpenAPIGenerator) AddServer(url, description string) {
	g.spec.Servers = append(g.spec.Servers, OpenAPIServer{
		URL:         url,
		Description: description,
	})
}

// AddRoute adds a route to the spec.
func (g *OpenAPIGenerator) AddRoute(method, path, summary string, tags []string) {
	if g.spec.Paths[path] == nil {
		g.spec.Paths[path] = make(PathItem)
	}

	op := Operation{
		Summary: summary,
		Tags:    tags,
		Responses: map[string]APIResponse{
			"200": {Description: "Successful response"},
		},
	}

	// Extract path parameters.
	parts := strings.Split(path, "/")
	for _, part := range parts {
		if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			paramName := strings.TrimSuffix(strings.TrimPrefix(part, "{"), "}")
			op.Parameters = append(op.Parameters, Parameter{
				Name:     paramName,
				In:       "path",
				Required: true,
				Schema:   Schema{Type: "string"},
			})
		}
	}

	g.spec.Paths[path][strings.ToLower(method)] = op
}

// AddRouteWithBody adds a route with request body schema.
func (g *OpenAPIGenerator) AddRouteWithBody(method, path, summary string, tags []string, bodyType any) {
	g.AddRoute(method, path, summary, tags)

	if bodyType != nil {
		pathItem := g.spec.Paths[path]
		op := pathItem[strings.ToLower(method)]
		op.RequestBody = &RequestBody{
			Required: true,
			Content: map[string]MediaType{
				"application/json": {
					Schema: reflectSchema(reflect.TypeOf(bodyType)),
				},
			},
		}
		pathItem[strings.ToLower(method)] = op
	}
}

// FromRoutes auto-generates spec from registered routes.
func (g *OpenAPIGenerator) FromRoutes(routes []RouteInfo) {
	for _, r := range routes {
		// Convert chi path params to OpenAPI format.
		path := r.Path
		path = strings.ReplaceAll(path, ":", "")

		tags := []string{}
		if r.Controller != "" {
			// Extract controller name for tags.
			parts := strings.Split(r.Controller, ".")
			tag := parts[len(parts)-1]
			tag = strings.TrimPrefix(tag, "*")
			tags = []string{tag}
		}

		g.AddRoute(r.Method, path, r.Summary, tags)
	}
}

// JSON returns the spec as JSON.
func (g *OpenAPIGenerator) JSON() ([]byte, error) {
	return json.MarshalIndent(g.spec, "", "  ")
}

// Spec returns the spec struct.
func (g *OpenAPIGenerator) Spec() OpenAPISpec {
	return g.spec
}

// reflectSchema creates a Schema from a Go type.
func reflectSchema(t reflect.Type) Schema {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	schema := Schema{
		Properties: make(map[string]Schema),
	}

	switch t.Kind() {
	case reflect.Struct:
		schema.Type = "object"
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			jsonTag := field.Tag.Get("json")
			if jsonTag == "-" {
				continue
			}
			name := strings.Split(jsonTag, ",")[0]
			if name == "" {
				name = field.Name
			}

			fieldSchema := typeToSchema(field.Type)
			schema.Properties[name] = fieldSchema

			if field.Tag.Get("validate") == "required" {
				schema.Required = append(schema.Required, name)
			}
		}
	default:
		return typeToSchema(t)
	}

	return schema
}

func typeToSchema(t reflect.Type) Schema {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.String:
		return Schema{Type: "string"}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return Schema{Type: "integer"}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return Schema{Type: "integer"}
	case reflect.Float32, reflect.Float64:
		return Schema{Type: "number"}
	case reflect.Bool:
		return Schema{Type: "boolean"}
	case reflect.Slice, reflect.Array:
		itemSchema := typeToSchema(t.Elem())
		return Schema{Type: "array", Items: &itemSchema}
	case reflect.Struct:
		return reflectSchema(t)
	default:
		return Schema{Type: "string"}
	}
}

// SwaggerUIHandler returns an HTTP handler that serves the Swagger UI.
func SwaggerUIHandler(specPath string) nethttp.HandlerFunc {
	return func(w nethttp.ResponseWriter, r *nethttp.Request) {
		html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>API Documentation</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
    SwaggerUIBundle({
        url: "%s",
        dom_id: '#swagger-ui',
        presets: [SwaggerUIBundle.presets.apis, SwaggerUIBundle.SwaggerUIStandalonePreset],
        layout: "BaseLayout"
    });
    </script>
</body>
</html>`, specPath)
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte(html))
	}
}
