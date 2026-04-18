package common

import "fmt"

// HttpException is an HTTP error with a status code.
type HttpException struct {
	Status  int    `json:"statusCode"`
	Message string `json:"message"`
}

func (e *HttpException) Error() string {
	return fmt.Sprintf("%d: %s", e.Status, e.Message)
}

// NewHttpException creates a new HttpException.
func NewHttpException(status int, message string) *HttpException {
	return &HttpException{Status: status, Message: message}
}

// Common HTTP exceptions.
func BadRequest(msg string) *HttpException          { return NewHttpException(400, msg) }
func Unauthorized(msg string) *HttpException        { return NewHttpException(401, msg) }
func Forbidden(msg string) *HttpException           { return NewHttpException(403, msg) }
func NotFound(msg string) *HttpException            { return NewHttpException(404, msg) }
func Conflict(msg string) *HttpException            { return NewHttpException(409, msg) }
func UnprocessableEntity(msg string) *HttpException { return NewHttpException(422, msg) }
func InternalError(msg string) *HttpException       { return NewHttpException(500, msg) }
