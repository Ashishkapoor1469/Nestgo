// Package config provides environment-based configuration loading
// with struct binding, validation, and multi-environment support.
package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Environment represents the application environment.
type Environment string

const (
	Development Environment = "development"
	Staging     Environment = "staging"
	Production  Environment = "production"
	Testing     Environment = "testing"
)

// Env returns the current environment from APP_ENV / NESTGO_ENV.
func Env() Environment {
	env := os.Getenv("NESTGO_ENV")
	if env == "" {
		env = os.Getenv("APP_ENV")
	}
	switch strings.ToLower(env) {
	case "production", "prod":
		return Production
	case "staging", "stage":
		return Staging
	case "testing", "test":
		return Testing
	default:
		return Development
	}
}

// IsProd returns true if running in production.
func IsProd() bool {
	return Env() == Production
}

// Load loads configuration from .env files and environment variables
// into the given struct. The struct fields should use the `env` tag
// to specify the environment variable name, and `default` tag for defaults.
//
// Example:
//
//	type AppConfig struct {
//	    Port     int    `env:"PORT" default:"3000"`
//	    DBHost   string `env:"DB_HOST" default:"localhost"`
//	    Debug    bool   `env:"DEBUG" default:"false"`
//	}
//
//	cfg, err := config.Load[AppConfig](".")
func Load[T any](paths ...string) (*T, error) {
	// Load .env files in order of specificity
	envFiles := []string{".env"}
	env := Env()
	if env != "" {
		envFiles = append(envFiles, fmt.Sprintf(".env.%s", env))
	}
	envFiles = append(envFiles, ".env.local")

	for _, f := range envFiles {
		for _, basePath := range paths {
			path := basePath + "/" + f
			_ = godotenv.Load(path) // ignore missing files
		}
		_ = godotenv.Load(f) // also try current dir
	}

	cfg := new(T)
	if err := bind(cfg); err != nil {
		return nil, fmt.Errorf("config: failed to bind: %w", err)
	}

	// Run validation if the config implements Validatable.
	if v, ok := any(cfg).(Validatable); ok {
		if err := v.Validate(); err != nil {
			return nil, fmt.Errorf("config: validation failed: %w", err)
		}
	}

	return cfg, nil
}

// MustLoad loads configuration or panics.
func MustLoad[T any](paths ...string) *T {
	cfg, err := Load[T](paths...)
	if err != nil {
		panic(err)
	}
	return cfg
}

// Validatable can be implemented by config structs to add custom validation.
type Validatable interface {
	Validate() error
}

// bind maps environment variables to struct fields.
func bind(target any) error {
	v := reflect.ValueOf(target)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldVal := v.Field(i)

		if !fieldVal.CanSet() {
			continue
		}

		// Handle embedded structs recursively.
		if field.Type.Kind() == reflect.Struct && field.Anonymous {
			if err := bind(fieldVal.Addr().Interface()); err != nil {
				return err
			}
			continue
		}

		envKey := field.Tag.Get("env")
		if envKey == "" {
			continue
		}

		defaultVal := field.Tag.Get("default")
		required := field.Tag.Get("required") == "true"

		value := os.Getenv(envKey)
		if value == "" {
			value = defaultVal
		}

		if value == "" && required {
			return fmt.Errorf("config: required environment variable %s is not set", envKey)
		}

		if value == "" {
			continue
		}

		if err := setField(fieldVal, value); err != nil {
			return fmt.Errorf("config: failed to set field %s: %w", field.Name, err)
		}
	}

	return nil
}

// setField sets a reflect.Value to the string value, handling type conversion.
func setField(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(n)
	case reflect.Float32, reflect.Float64:
		n, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		field.SetFloat(n)
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(b)
	case reflect.Slice:
		if field.Type().Elem().Kind() == reflect.String {
			parts := strings.Split(value, ",")
			for i, p := range parts {
				parts[i] = strings.TrimSpace(p)
			}
			field.Set(reflect.ValueOf(parts))
		}
	default:
		return fmt.Errorf("unsupported type: %s", field.Kind())
	}
	return nil
}
