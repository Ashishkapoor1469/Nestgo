package common

import (
	"encoding/json"
	"fmt"
)

// ─── Model Transformers ─────────────────────────────────────────────────────
//
// Lightweight helpers to convert between internal models and response schemas.
// Uses JSON round-trip for simplicity and reliability — avoids heavy reflection
// while supporting all JSON-tagged fields.
//
// Example:
//
//	type User struct {
//	    ID       int    `json:"id"`
//	    Name     string `json:"name"`
//	    Email    string `json:"email"`
//	    Password string `json:"password"`
//	}
//
//	type UserResponse struct {
//	    ID    int    `json:"id"`
//	    Name  string `json:"name"`
//	    Email string `json:"email"`
//	}
//
//	response, err := common.Transform[User, UserResponse](user)

// Transform converts a source value to a destination type by mapping
// matching JSON-tagged fields. Fields present in the source but absent
// in the destination are silently dropped.
func Transform[From any, To any](src From) (To, error) {
	var dst To
	data, err := json.Marshal(src)
	if err != nil {
		return dst, fmt.Errorf("transform: marshal failed: %w", err)
	}
	if err := json.Unmarshal(data, &dst); err != nil {
		return dst, fmt.Errorf("transform: unmarshal failed: %w", err)
	}
	return dst, nil
}

// TransformSlice converts a slice of source values to a slice of destination type.
func TransformSlice[From any, To any](src []From) ([]To, error) {
	result := make([]To, 0, len(src))
	for _, item := range src {
		transformed, err := Transform[From, To](item)
		if err != nil {
			return nil, err
		}
		result = append(result, transformed)
	}
	return result, nil
}

// MustTransform is like Transform but panics on error.
func MustTransform[From any, To any](src From) To {
	dst, err := Transform[From, To](src)
	if err != nil {
		panic(err)
	}
	return dst
}

// MustTransformSlice is like TransformSlice but panics on error.
func MustTransformSlice[From any, To any](src []From) []To {
	dst, err := TransformSlice[From, To](src)
	if err != nil {
		panic(err)
	}
	return dst
}
