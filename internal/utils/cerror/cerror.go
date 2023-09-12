package cerror

import (
	"fmt"
)

// customError is a type that wraps an error message.
type customError struct {
	detail string
}

// NewCustomError creates a new custom error with the given detail.
func NewCustomError(detail string) customError {
	return customError{detail: detail}
}

// Error returns common error message format.
func (e customError) Error() string {
	return fmt.Sprintf("error: %s", e.detail)
}
