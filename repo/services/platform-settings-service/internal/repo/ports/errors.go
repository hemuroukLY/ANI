package ports

import "errors"

// ErrNotImplemented marks a port adapter method that is declared but not yet implemented.
var ErrNotImplemented = errors.New("NOT_IMPLEMENTED")
