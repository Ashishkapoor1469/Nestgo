// Package testing provides test utilities for NestGo applications.
package testing

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/Ashishkapoor1469/Nestgo/common"
	"github.com/Ashishkapoor1469/Nestgo/core"
	"github.com/Ashishkapoor1469/Nestgo/di"
)

// TestApp provides a test harness for NestGo applications.
type TestApp struct {
	App       *core.NestGoApp
	Container *di.Container
	Server    *httptest.Server
}

// NewTestApp creates a new test application.
func NewTestApp(mods ...common.Module) *TestApp {
	app := core.New(core.WithAddress(":0"))

	for _, m := range mods {
		app.RegisterModule(m)
	}

	return &TestApp{
		App:       app,
		Container: app.Container(),
	}
}

// Start starts the test server.
func (t *TestApp) Start() *TestApp {
	t.Server = httptest.NewServer(t.App.Router().Handler())
	return t
}

// Close shuts down the test server.
func (t *TestApp) Close() {
	if t.Server != nil {
		t.Server.Close()
	}
}

// URL returns the test server URL.
func (t *TestApp) URL() string {
	return t.Server.URL
}

// OverrideProvider replaces a provider in the DI container (for mocking).
func (t *TestApp) OverrideProvider(value any) *TestApp {
	_ = t.Container.ProvideValue(value)
	return t
}

// --- HTTP Test Helpers ---

// TestRequest builds and sends test HTTP requests.
type TestRequest struct {
	method  string
	path    string
	body    any
	headers map[string]string
	baseURL string
}

// GET creates a GET test request.
func GET(baseURL, path string) *TestRequest {
	return &TestRequest{
		method:  http.MethodGet,
		path:    path,
		headers: make(map[string]string),
		baseURL: baseURL,
	}
}

// POST creates a POST test request.
func POST(baseURL, path string, body any) *TestRequest {
	return &TestRequest{
		method:  http.MethodPost,
		path:    path,
		body:    body,
		headers: make(map[string]string),
		baseURL: baseURL,
	}
}

// PUT creates a PUT test request.
func PUT(baseURL, path string, body any) *TestRequest {
	return &TestRequest{
		method:  http.MethodPut,
		path:    path,
		body:    body,
		headers: make(map[string]string),
		baseURL: baseURL,
	}
}

// DELETE creates a DELETE test request.
func DELETE(baseURL, path string) *TestRequest {
	return &TestRequest{
		method:  http.MethodDelete,
		path:    path,
		headers: make(map[string]string),
		baseURL: baseURL,
	}
}

// WithHeader adds a header.
func (r *TestRequest) WithHeader(key, value string) *TestRequest {
	r.headers[key] = value
	return r
}

// WithAuth adds a bearer token.
func (r *TestRequest) WithAuth(token string) *TestRequest {
	r.headers["Authorization"] = "Bearer " + token
	return r
}

// Send sends the request and returns the response.
func (r *TestRequest) Send() (*TestResponse, error) {
	var bodyReader io.Reader
	if r.body != nil {
		data, err := json.Marshal(r.body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(r.method, r.baseURL+r.path, bodyReader)
	if err != nil {
		return nil, err
	}

	if r.body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	for k, v := range r.headers {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return &TestResponse{Response: resp}, nil
}

// TestResponse wraps an HTTP response with convenience methods.
type TestResponse struct {
	*http.Response
}

// JSON decodes the response body as JSON.
func (r *TestResponse) JSON(v any) error {
	defer func() { _ = r.Body.Close() }()
	return json.NewDecoder(r.Body).Decode(v)
}

// BodyString returns the response body as a string.
func (r *TestResponse) BodyString() string {
	defer func() { _ = r.Body.Close() }()
	data, _ := io.ReadAll(r.Body)
	return string(data)
}

// IsOK returns true if status is 200.
func (r *TestResponse) IsOK() bool {
	return r.StatusCode == http.StatusOK
}

// IsCreated returns true if status is 201.
func (r *TestResponse) IsCreated() bool {
	return r.StatusCode == http.StatusCreated
}

// IsNotFound returns true if status is 404.
func (r *TestResponse) IsNotFound() bool {
	return r.StatusCode == http.StatusNotFound
}

// --- Mock Utilities ---

// MockProvider creates a mock provider for testing.
func MockProvider[T any](mock T) di.Provider {
	return di.Provider{
		Instance: mock,
	}
}
