package errors

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync/atomic"
)

var errorCounter int32

// CError represents a custom error with code, message, and optional details.
// It supports error wrapping and provides two levels of details:
// - details: technical details for logging and debugging
// - userdetails: user-friendly details safe to expose to end users
//
// CError is designed for flexible error handling:
// - Create base errors once with New() or Wrap()
// - Customize message and code later with WithMessage() and WithCode()
// - Errors with the same ID are considered equal even with different messages/codes
// - This allows reusing the same logical error in different contexts
type CError struct {
	wrapped     error
	details     []string
	userdetails []string
	msg         string
	id          int32
	code        uint16
}

// New creates a new CError with the given code and message.
func New(code uint16, msg string) *CError {
	id := atomic.AddInt32(&errorCounter, 1)
	return &CError{
		id:   id,
		code: code,
		msg:  msg,
	}
}

// Wrap creates a new CError that wraps the given error with additional code and message.
func Wrap(err error, code uint16, msg string) *CError {
	id := atomic.AddInt32(&errorCounter, 1)
	return &CError{
		id:      id,
		wrapped: err,
		code:    code,
		msg:     msg,
	}
}

func (e *CError) Error() string {
	var errstr strings.Builder

	errstr.WriteString(fmt.Sprintf("message: %s, code: %d", e.msg, e.code))
	if e.details != nil {
		errstr.WriteString(", details: [")
		for i, detail := range e.details {
			if i > 0 {
				errstr.WriteString(", ")
			}
			errstr.WriteString(detail)
		}
		errstr.WriteString("]")
	}

	if e.wrapped != nil {
		errstr.WriteString(fmt.Sprintf(", wrapped: %s", e.wrapped.Error()))
	}

	return errstr.String()
}

// Code returns the error code as int.
func (e *CError) Code() int {
	return int(e.code)
}

// Message returns the error message.
func (c *CError) Message() string {
	return c.msg
}

// UserDetails returns user-safe details that can be exposed to end users.
func (c *CError) UserDetails() []string {
	return c.userdetails
}

// Details returns technical details intended for logging and debugging.
func (c *CError) Details() []string {
	return c.details
}

// Is reports whether any error in err's chain matches target.
// Two CErrors match if they have the same unique ID, regardless of their
// current message or code values. This allows comparing logical error identity.
func (e *CError) Is(target error) bool {
	if cerr, ok := target.(*CError); ok {
		return e.id == cerr.id
	}
	return false
}

// As finds the first error in err's chain that matches target.
func (e *CError) As(target interface{}) bool {
	if cerr, ok := target.(**CError); ok {
		*cerr = e
		return true
	}
	return false
}

// Unwrap returns the wrapped error, if any.
func (c *CError) Unwrap() error {
	return c.wrapped
}

// UnwrapAll returns all errors in the wrapping chain.
func (c *CError) UnwrapAll() []error {
	if c.wrapped == nil {
		return nil
	}

	var wrapped []error
	current := c.wrapped

	for current != nil {
		wrapped = append(wrapped, current)

		switch err := current.(type) {
		case *CError:
			current = err.wrapped
		case interface{ Unwrap() error }:
			current = err.Unwrap()
		default:
			current = nil
		}
	}

	return wrapped
}

// WithCode sets the error code and returns a new error copy for chaining.
func (e *CError) WithCode(code uint16) *CError {
	newErr := *e
	newErr.code = code
	return &newErr
}

// WithMessage sets the error message and returns a new error copy for chaining.
func (e *CError) WithMessage(msg string) *CError {
	newErr := *e
	newErr.msg = msg
	return &newErr
}

// WithWrap wraps another error and returns a new error copy for chaining.
func (e *CError) WithWrap(err error) *CError {
	newErr := *e
	newErr.wrapped = err
	return &newErr
}

// WithDetail adds a technical detail to the error and returns a new error copy for chaining.
func (e *CError) WithDetail(detail string) *CError {
	newErr := *e
	newErr.details = make([]string, len(e.details), len(e.details)+1)
	copy(newErr.details, e.details)
	newErr.details = append(newErr.details, detail)
	return &newErr
}

// WithUserDetail adds a user-safe detail to the error and returns a new error copy for chaining.
func (e *CError) WithUserDetail(detail string) *CError {
	newErr := *e
	newErr.userdetails = make([]string, len(e.userdetails), len(e.userdetails)+1)
	copy(newErr.userdetails, e.userdetails)
	newErr.userdetails = append(newErr.userdetails, detail)
	return &newErr
}

// UserMessage returns a user-friendly error message.
func (e *CError) UserMessage() string {
	if e.msg != "" {
		return e.msg
	}
	return "unknown error occurred"
}

// IsCode checks if the error has the specified code.
// Returns false if the provided code doesn't fit in uint16 range.
func (e *CError) IsCode(code int) bool {
	if code < 0 || code > 65535 {
		return false
	}
	return e.code == uint16(code)
}

// MarshalJSON serializes the error to JSON with minimal user-safe information.
// Only includes message, code, and user-safe details.
// Technical details, wrapped errors, and internal IDs are excluded for security.
func (e *CError) MarshalJSON() ([]byte, error) {
	type errorJSON struct {
		Message string   `json:"message"`
		Code    uint16   `json:"code"`
		Details []string `json:"details,omitempty"`
	}

	data := errorJSON{
		Message: e.msg,
		Code:    e.code,
		Details: e.userdetails,
	}

	return json.Marshal(data)
}

func (e *CError) UnmarshalJSON(data []byte) error {
	type errorJSON struct {
		Message string   `json:"message"`
		Code    uint16   `json:"code"`
		Details []string `json:"details,omitempty"`
	}

	var tmp errorJSON
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	e.msg = tmp.Message
	e.code = tmp.Code
	e.userdetails = tmp.Details

	return nil
}

// FullData returns a structured object containing complete error information
// including technical details, wrapped error chains, and internal IDs.
// This method expands all wrapped errors recursively and should be used
// for logging, debugging, and monitoring purposes only.
func (e *CError) FullData() map[string]interface{} {
	data := map[string]interface{}{
		"message": e.msg,
		"code":    e.code,
	}

	if len(e.details) > 0 {
		data["details"] = e.details
	}

	if len(e.userdetails) > 0 {
		data["user_details"] = e.userdetails
	}

	if e.wrapped != nil {
		wrappedChain := e.expandAllWrapped(e.wrapped)
		if len(wrappedChain) > 0 {
			data["wrapped"] = wrappedChain
		}
	}

	return data
}

func (e *CError) expandAllWrapped(err error) []map[string]interface{} {
	var chain []map[string]interface{}
	current := err

	for current != nil {
		wrappedData := make(map[string]interface{})

		if cerr, ok := current.(*CError); ok {
			wrappedData["message"] = cerr.msg
			wrappedData["code"] = cerr.code

			if len(cerr.details) > 0 {
				wrappedData["details"] = cerr.details
			}
			if len(cerr.userdetails) > 0 {
				wrappedData["user_details"] = cerr.userdetails
			}

			current = cerr.wrapped
		} else {
			wrappedData["message"] = current.Error()

			if unwrapper, ok := current.(interface{ Unwrap() error }); ok {
				current = unwrapper.Unwrap()
			} else {
				current = nil
			}
		}

		chain = append(chain, wrappedData)
	}

	return chain
}
