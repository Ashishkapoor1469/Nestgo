package common

import (
	"fmt"
	"net"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

// ─── Tag-Based Struct Validation ────────────────────────────────────────────
//
// Provides automatic validation of structs using `validate` struct tags.
// This complements the existing fluent Validator for cases where developers
// prefer a declarative, tag-based approach.
//
// Supported tags (comma-separated):
//   - required       — field must be non-zero/non-empty
//   - min=N          — minimum length (string) or value (numeric)
//   - max=N          — maximum length (string) or value (numeric)
//   - email          — valid email format
//   - url            — valid URL format
//   - numeric        — string must contain only digits
//   - alpha          — string must contain only letters
//   - alphanum       — string must contain only letters and digits
//   - oneof=a|b|c    — value must be one of the listed options
//   - pattern=regex  — value must match the regex pattern
//   - ip             — valid IP address
//   - uuid           — valid UUID format
//
// Example:
//
//	type CreateUserSchema struct {
//	    Name     string `json:"name" validate:"required,min=3,max=100"`
//	    Email    string `json:"email" validate:"required,email"`
//	    Password string `json:"password" validate:"required,min=6"`
//	    Role     string `json:"role" validate:"oneof=admin|user|mod"`
//	}

// ─── Custom Validator Registry ──────────────────────────────────────────────

// ValidatorFunc is a custom validation function.
// It receives the field value and returns true if valid.
type ValidatorFunc func(value any) bool

var (
	customValidators   = make(map[string]ValidatorFunc)
	customValidatorsMu sync.RWMutex
)

// RegisterValidator registers a custom validator by name for use in validate tags.
//
// Example:
//
//	common.RegisterValidator("strong_password", func(v any) bool {
//	    s, ok := v.(string)
//	    if !ok { return false }
//	    return len(s) >= 8 && hasUppercase(s) && hasDigit(s)
//	})
//
// Usage in struct tag: `validate:"required,strong_password"`
func RegisterValidator(name string, fn ValidatorFunc) {
	customValidatorsMu.Lock()
	defer customValidatorsMu.Unlock()
	customValidators[name] = fn
}

// ─── Type Cache ─────────────────────────────────────────────────────────────

// Cached parsed validation rules per struct type to avoid re-parsing tags.
var typeCache sync.Map // map[reflect.Type][]fieldRules

type fieldRules struct {
	index int
	name  string // json name or field name
	rules []rule
}

type rule struct {
	name  string // e.g., "required", "min", "email"
	param string // e.g., "3" for min=3
}

// ─── Core Validation ────────────────────────────────────────────────────────

// ValidateStruct validates a struct using its `validate` tags.
// Returns nil if all validations pass, or a *ValidationErrors with field-level errors.
func ValidateStruct(v any) error {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return fmt.Errorf("validate: nil pointer")
		}
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("validate: expected struct, got %s", val.Kind())
	}

	rules := getRules(val.Type())
	if len(rules) == 0 {
		return nil
	}

	errs := &ValidationErrors{}
	for _, fr := range rules {
		fieldVal := val.Field(fr.index)
		for _, r := range fr.rules {
			if msg := applyRule(r, fieldVal, fr.name); msg != "" {
				errs.Add(fr.name, msg)
				break // one error per field
			}
		}
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}

// HasValidateTags checks if a struct type has any validate tags.
func HasValidateTags(v any) bool {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return false
	}
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Tag.Get("validate") != "" {
			return true
		}
	}
	return false
}

// getRules returns cached or freshly-parsed validation rules for a type.
func getRules(t reflect.Type) []fieldRules {
	if cached, ok := typeCache.Load(t); ok {
		return cached.([]fieldRules)
	}

	var parsed []fieldRules
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("validate")
		if tag == "" || tag == "-" {
			continue
		}

		// Determine the field name (prefer json tag).
		name := field.Name
		if jsonTag := field.Tag.Get("json"); jsonTag != "" {
			parts := strings.Split(jsonTag, ",")
			if parts[0] != "" && parts[0] != "-" {
				name = parts[0]
			}
		}

		fr := fieldRules{
			index: i,
			name:  name,
			rules: parseRules(tag),
		}
		parsed = append(parsed, fr)
	}

	typeCache.Store(t, parsed)
	return parsed
}

// parseRules parses a validate tag string into individual rules.
func parseRules(tag string) []rule {
	parts := strings.Split(tag, ",")
	var rules []rule
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		r := rule{name: p}
		if idx := strings.Index(p, "="); idx > 0 {
			r.name = p[:idx]
			r.param = p[idx+1:]
		}
		rules = append(rules, r)
	}
	return rules
}

// ─── Rule Application ───────────────────────────────────────────────────────

// applyRule applies a single validation rule and returns an error message or "".
func applyRule(r rule, val reflect.Value, fieldName string) string {
	switch r.name {
	case "required":
		return validateRequired(val)
	case "min":
		return validateMin(val, r.param)
	case "max":
		return validateMax(val, r.param)
	case "email":
		return validateEmail(val)
	case "url":
		return validateURL(val)
	case "numeric":
		return validateNumeric(val)
	case "alpha":
		return validateAlpha(val)
	case "alphanum":
		return validateAlphanum(val)
	case "oneof":
		return validateOneOf(val, r.param)
	case "pattern":
		return validatePattern(val, r.param)
	case "ip":
		return validateIP(val)
	case "uuid":
		return validateUUID(val)
	default:
		// Check custom validators.
		customValidatorsMu.RLock()
		fn, ok := customValidators[r.name]
		customValidatorsMu.RUnlock()
		if ok {
			if !fn(val.Interface()) {
				return fmt.Sprintf("failed %s validation", r.name)
			}
			return ""
		}
		return ""
	}
}

func validateRequired(val reflect.Value) string {
	switch val.Kind() {
	case reflect.String:
		if strings.TrimSpace(val.String()) == "" {
			return "is required"
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if val.Int() == 0 {
			return "is required"
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if val.Uint() == 0 {
			return "is required"
		}
	case reflect.Float32, reflect.Float64:
		if val.Float() == 0 {
			return "is required"
		}
	case reflect.Bool:
		// bools are always "set"
	case reflect.Slice, reflect.Map:
		if val.IsNil() || val.Len() == 0 {
			return "is required"
		}
	case reflect.Ptr, reflect.Interface:
		if val.IsNil() {
			return "is required"
		}
	default:
		if val.IsZero() {
			return "is required"
		}
	}
	return ""
}

func validateMin(val reflect.Value, param string) string {
	n, _ := strconv.Atoi(param)
	switch val.Kind() {
	case reflect.String:
		if len(val.String()) < n {
			return fmt.Sprintf("must be at least %d characters", n)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if val.Int() < int64(n) {
			return fmt.Sprintf("must be at least %d", n)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if val.Uint() < uint64(n) {
			return fmt.Sprintf("must be at least %d", n)
		}
	case reflect.Float32, reflect.Float64:
		if val.Float() < float64(n) {
			return fmt.Sprintf("must be at least %d", n)
		}
	case reflect.Slice, reflect.Map:
		if val.Len() < n {
			return fmt.Sprintf("must have at least %d items", n)
		}
	}
	return ""
}

func validateMax(val reflect.Value, param string) string {
	n, _ := strconv.Atoi(param)
	switch val.Kind() {
	case reflect.String:
		if len(val.String()) > n {
			return fmt.Sprintf("must be at most %d characters", n)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if val.Int() > int64(n) {
			return fmt.Sprintf("must be at most %d", n)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if val.Uint() > uint64(n) {
			return fmt.Sprintf("must be at most %d", n)
		}
	case reflect.Float32, reflect.Float64:
		if val.Float() > float64(n) {
			return fmt.Sprintf("must be at most %d", n)
		}
	case reflect.Slice, reflect.Map:
		if val.Len() > n {
			return fmt.Sprintf("must have at most %d items", n)
		}
	}
	return ""
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func validateEmail(val reflect.Value) string {
	if val.Kind() != reflect.String {
		return ""
	}
	s := val.String()
	if s == "" {
		return "" // skip empty (use required for non-empty)
	}
	if !emailRegex.MatchString(s) {
		return "must be a valid email address"
	}
	return ""
}

var urlRegex = regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)

func validateURL(val reflect.Value) string {
	if val.Kind() != reflect.String {
		return ""
	}
	s := val.String()
	if s == "" {
		return ""
	}
	if !urlRegex.MatchString(s) {
		return "must be a valid URL"
	}
	return ""
}

func validateNumeric(val reflect.Value) string {
	if val.Kind() != reflect.String {
		return ""
	}
	s := val.String()
	if s == "" {
		return ""
	}
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return "must contain only digits"
		}
	}
	return ""
}

func validateAlpha(val reflect.Value) string {
	if val.Kind() != reflect.String {
		return ""
	}
	s := val.String()
	if s == "" {
		return ""
	}
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return "must contain only letters"
		}
	}
	return ""
}

func validateAlphanum(val reflect.Value) string {
	if val.Kind() != reflect.String {
		return ""
	}
	s := val.String()
	if s == "" {
		return ""
	}
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return "must contain only letters and digits"
		}
	}
	return ""
}

func validateOneOf(val reflect.Value, param string) string {
	if val.Kind() != reflect.String {
		return ""
	}
	s := val.String()
	if s == "" {
		return ""
	}
	options := strings.Split(param, "|")
	for _, opt := range options {
		if s == strings.TrimSpace(opt) {
			return ""
		}
	}
	return fmt.Sprintf("must be one of: %s", strings.Join(options, ", "))
}

func validatePattern(val reflect.Value, param string) string {
	if val.Kind() != reflect.String {
		return ""
	}
	s := val.String()
	if s == "" {
		return ""
	}
	matched, err := regexp.MatchString(param, s)
	if err != nil {
		return fmt.Sprintf("invalid pattern: %s", param)
	}
	if !matched {
		return fmt.Sprintf("must match pattern: %s", param)
	}
	return ""
}

func validateIP(val reflect.Value) string {
	if val.Kind() != reflect.String {
		return ""
	}
	s := val.String()
	if s == "" {
		return ""
	}
	if net.ParseIP(s) == nil {
		return "must be a valid IP address"
	}
	return ""
}

var uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

func validateUUID(val reflect.Value) string {
	if val.Kind() != reflect.String {
		return ""
	}
	s := val.String()
	if s == "" {
		return ""
	}
	if !uuidRegex.MatchString(s) {
		return "must be a valid UUID"
	}
	return ""
}
