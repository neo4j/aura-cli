package clierr

import "fmt"

// Usage Error, require feedback
func NewUsageError(msg string, a ...any) error {
	return fmt.Errorf(msg, a...)
}

// API errors, retry may solve it
func NewUpstreamError(msg string, a ...any) error {
	return fmt.Errorf(msg, a...)
}

// Fatal error, unrecoverable
func NewFatalError(msg string, a ...any) error {
	return fmt.Errorf(msg, a...)
}
