package benchmarks

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ashishkapoor1469/Nestgo/common"
	nesthttp "github.com/Ashishkapoor1469/Nestgo/http"
	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi/v5"
	"github.com/gofiber/fiber/v2"
)

// Mock controller for NestGo
type TestController struct{}

func (c *TestController) Prefix() string { return "/test" }
func (c *TestController) Routes() []common.Route {
	return []common.Route{
		{
			Method: "GET",
			Path:   "/hello",
			Handler: func(ctx *common.Context) error {
				ctx.Writer.WriteHeader(http.StatusOK)
				_, err := ctx.Writer.Write([]byte("Hello World"))
				return err
			},
		},
	}
}

func BenchmarkNestGo_Router(b *testing.B) {
	router := nesthttp.NewRouter()
	router.RegisterController(&TestController{})
	
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test/hello", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
	}
}

func BenchmarkGin_Router(b *testing.B) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.GET("/test/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test/hello", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
	}
}

func BenchmarkChi_Router(b *testing.B) {
	router := chi.NewRouter()
	router.Get("/test/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test/hello", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
	}
}


func BenchmarkFiber_Router(b *testing.B) {
	app := fiber.New()
	app.Get("/test/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})

	// Fiber has a built-in Test method for benchmarking
	req := httptest.NewRequest("GET", "/test/hello", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.Test(req)
	}
}
