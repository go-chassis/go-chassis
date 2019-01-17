package client

import (
	"fmt"
)

// RegistryException structure contains message and error information for the exception caused by service-center
type RegistryException struct {
	Title   string
	OrglErr error
	Message string
}

// Error gets the Error message from the Error
func (e *RegistryException) Error() string {
	if e.OrglErr == nil {
		return fmt.Sprintf("%s(%s)", e.Title, e.Message)
	}
	return fmt.Sprintf("%s(%s), %s", e.Title, e.OrglErr.Error(), e.Message)
}

func formatMessage(args []interface{}) string {
	if len(args) == 0 {
		return ""
	}
	format, ok := args[0].(string)
	if !ok {
		return fmt.Sprintf("%v", args)
	}
	return fmt.Sprintf(format, args[1:]...)
}

func newException(t string, e error, message string) *RegistryException {
	return &RegistryException{
		Title:   t,
		OrglErr: e,
		Message: message,
	}
}

// NewCommonException creates a generic exception
func NewCommonException(format string, args ...interface{}) error {
	return newException("Common exception", nil, fmt.Sprintf(format, args...))
}

// NewJSONException creates a JSON exception
func NewJSONException(e error, args ...interface{}) error {
	return newException("JSON exception", e, formatMessage(args))
}

// NewIOException create and IO exception
func NewIOException(e error, args ...interface{}) error {
	return newException("IO exception", e, formatMessage(args))
}
