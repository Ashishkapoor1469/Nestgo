package common

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidationError holds validation errors for a request body.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Tag     string `json:"tag,omitempty"`
}

// ValidationErrors is a collection of validation errors.
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

func (e *ValidationErrors) Error() string {
	var msgs []string
	for _, err := range e.Errors {
		msgs = append(msgs, fmt.Sprintf("%s: %s", err.Field, err.Message))
	}
	return "validation failed: " + strings.Join(msgs, "; ")
}

// Add adds a validation error.
func (e *ValidationErrors) Add(field, message string) {
	e.Errors = append(e.Errors, ValidationError{Field: field, Message: message})
}

// HasErrors returns true if there are any errors.
func (e *ValidationErrors) HasErrors() bool {
	return len(e.Errors) > 0
}

// Validator provides fluent validation for DTOs.
type Validator struct {
	errors *ValidationErrors
}

// NewValidator creates a new validator.
func NewValidator() *Validator {
	return &Validator{
		errors: &ValidationErrors{},
	}
}

// Required validates that a string is not empty.
func (v *Validator) Required(field, value string) *Validator {
	if strings.TrimSpace(value) == "" {
		v.errors.Add(field, "is required")
	}
	return v
}

// MinLength validates minimum string length.
func (v *Validator) MinLength(field, value string, min int) *Validator {
	if len(value) < min {
		v.errors.Add(field, fmt.Sprintf("must be at least %d characters", min))
	}
	return v
}

// MaxLength validates maximum string length.
func (v *Validator) MaxLength(field, value string, max int) *Validator {
	if len(value) > max {
		v.errors.Add(field, fmt.Sprintf("must be at most %d characters", max))
	}
	return v
}

// Min validates minimum integer value.
func (v *Validator) Min(field string, value, min int) *Validator {
	if value < min {
		v.errors.Add(field, fmt.Sprintf("must be at least %d", min))
	}
	return v
}

// Max validates maximum integer value.
func (v *Validator) Max(field string, value, max int) *Validator {
	if value > max {
		v.errors.Add(field, fmt.Sprintf("must be at most %d", max))
	}
	return v
}

// Email validates email format.
func (v *Validator) Email(field, value string) *Validator {
	if value == "" {
		return v
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(value) {
		v.errors.Add(field, "must be a valid email address")
	}
	return v
}

// URL validates URL format.
func (v *Validator) URL(field, value string) *Validator {
	if value == "" {
		return v
	}
	urlRegex := regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
	if !urlRegex.MatchString(value) {
		v.errors.Add(field, "must be a valid URL")
	}
	return v
}

// OneOf validates that a value is one of the allowed values.
func (v *Validator) OneOf(field, value string, allowed []string) *Validator {
	for _, a := range allowed {
		if value == a {
			return v
		}
	}
	v.errors.Add(field, fmt.Sprintf("must be one of: %s", strings.Join(allowed, ", ")))
	return v
}

// Pattern validates a value against a regex pattern.
func (v *Validator) Pattern(field, value, pattern string) *Validator {
	if value == "" {
		return v
	}
	if matched, _ := regexp.MatchString(pattern, value); !matched {
		v.errors.Add(field, fmt.Sprintf("must match pattern: %s", pattern))
	}
	return v
}

// Custom adds a custom validation rule.
func (v *Validator) Custom(field string, valid bool, message string) *Validator {
	if !valid {
		v.errors.Add(field, message)
	}
	return v
}

// Validate returns an error if there are validation failures.
func (v *Validator) Validate() error {
	if v.errors.HasErrors() {
		return v.errors
	}
	return nil
}
