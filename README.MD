# CError - Advanced Error Handling for Go

A powerful and flexible error handling package for Go applications that provides structured error management with HTTP status code integration, user-safe messaging, and comprehensive error chaining.

## Features

- 🔍 **Unique Error Identity**: Each error has a unique ID for logical comparison regardless of message changes
- 🔗 **Error Chaining**: Full support for error wrapping with `errors.Is`, `errors.As`, and `errors.Unwrap`
- 👥 **Dual Details System**: Separate technical details for logging and user-safe details for APIs
- 🌐 **HTTP Integration**: Error codes align with HTTP status codes for seamless API development
- 🔄 **Fluent API**: Chainable methods for easy error customization
- 📝 **Safe JSON Serialization**: Automatic filtering of sensitive information in JSON output
- 🎯 **Predefined Errors**: Common error types ready to use

## Installation

```bash
go get github.com/yourusername/errors
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/yourusername/errors"
)

func main() {
    // Use predefined errors
    err := errors.ErrUserNotFound.WithUserDetail("User ID: 12345")
    fmt.Println(err.Code()) // 404
    
    // Create custom errors
    customErr := errors.New(400, "validation failed").
        WithDetail("field: email, value: invalid@").
        WithUserDetail("Please provide a valid email address")
    
    // Wrap existing errors
    if dbErr := connectDB(); dbErr != nil {
        return errors.ErrDatabaseError.WithWrap(dbErr).
            WithDetail("connection string: localhost:5432")
    }
}
```

## Core API

### Creating Errors

```go
// Create new error
err := errors.New(400, "validation failed")

// Wrap existing error
err := errors.Wrap(originalErr, 500, "database operation failed")

// Use predefined errors
err := errors.ErrUserNotFound
```

### Fluent API Methods

```go
err := errors.ErrValidationInvalid.
    WithCode(422).                              // Change status code
    WithMessage("Custom validation message").   // Change message
    WithDetail("field: email, regex: ^[a-z]+"). // Add technical detail
    WithUserDetail("Email format is invalid").  // Add user-safe detail
    WithWrap(originalErr)                       // Wrap another error
```

### Error Information

```go
err.Code()        // HTTP status code as int
err.Message()     // Error message
err.Details()     // Technical details []string
err.UserDetails() // User-safe details []string
err.Error()       // Full error string with all information
err.UserMessage() // User-friendly message
```

### Error Comparison

```go
// Check specific error type (ignores message/code changes)
if errors.Is(err, errors.ErrUserNotFound) {
    // Handle user not found
}

// Extract specific error type
var cerr *errors.CError
if errors.As(err, &cerr) {
    fmt.Println("Status code:", cerr.Code())
}
```

## Predefined Errors

The package includes comprehensive predefined errors organized by HTTP status codes:

### Validation Errors (400)
- `ErrValidationRequired` - Field is required
- `ErrValidationInvalid` - Invalid field value
- `ErrValidationFormat` - Invalid format
- `ErrValidationEmail` - Invalid email format
- `ErrValidationPassword` - Password requirements not met

### Authentication Errors (401/403)
- `ErrAuthRequired` - Authentication required
- `ErrAuthInvalid` - Invalid credentials
- `ErrAuthExpired` - Authentication expired
- `ErrAuthPermissions` - Insufficient permissions
- `ErrAuthBlocked` - Account blocked

### Not Found Errors (404)
- `ErrNotFound` - Generic resource not found
- `ErrUserNotFound` - User not found
- `ErrFileNotFound` - File not found
- `ErrRecordNotFound` - Database record not found

### Conflict Errors (409)
- `ErrConflict` - Resource conflict
- `ErrAlreadyExists` - Resource already exists
- `ErrUserExists` - User already exists
- `ErrEmailTaken` - Email already taken

### Business Logic Errors (422)
- `ErrBusinessRule` - Business rule violation
- `ErrInsufficientFunds` - Insufficient funds
- `ErrOperationNotAllowed` - Operation not allowed
- `ErrLimitExceeded` - Limit exceeded

### Rate Limiting (429)
- `ErrRateLimited` - Rate limit exceeded
- `ErrTooManyRequests` - Too many requests
- `ErrQuotaExceeded` - Quota exceeded

### Server Errors (500+)
- `ErrInternalError` - Internal server error
- `ErrDatabaseError` - Database error
- `ErrServiceUnavailable` - Service unavailable
- `ErrTimeoutError` - Operation timeout

[View all predefined errors](common_errors.go)

## Usage Examples

### HTTP API Integration

```go
func getUserHandler(w http.ResponseWriter, r *http.Request) {
    user, err := getUserByID(userID)
    if err != nil {
        var cerr *errors.CError
        if errors.As(err, &cerr) {
            http.Error(w, cerr.UserMessage(), cerr.Code())
            return
        }
        http.Error(w, "Internal server error", 500)
        return
    }
    
    json.NewEncoder(w).Encode(user)
}

func getUserByID(id string) (*User, error) {
    if id == "" {
        return nil, errors.ErrValidationRequired.
            WithUserDetail("User ID is required")
    }
    
    user, err := db.FindUser(id)
    if err != nil {
        if isNotFoundError(err) {
            return nil, errors.ErrUserNotFound.
                WithDetail(fmt.Sprintf("user_id: %s", id))
        }
        return nil, errors.ErrDatabaseError.
            WithWrap(err).
            WithDetail(fmt.Sprintf("query: SELECT * FROM users WHERE id = %s", id))
    }
    
    return user, nil
}
```

### JSON Serialization

```go
err := errors.ErrValidationInvalid.
    WithDetail("technical info for logs").
    WithUserDetail("Please check your input")

// Safe JSON output (excludes technical details)
jsonData, _ := json.Marshal(err)
fmt.Println(string(jsonData))
// Output: {"message":"field value is invalid","code":400,"details":["Please check your input"]}

// Full data for logging
logData := err.FullData()
logger.Error("Validation failed", "error", logData)
```

### Error Chaining and Unwrapping

```go
func processFile(filename string) error {
    data, err := os.ReadFile(filename)
    if err != nil {
        return errors.ErrFileNotFound.
            WithWrap(err).
            WithDetail(fmt.Sprintf("filename: %s", filename)).
            WithUserDetail("The requested file could not be found")
    }
    
    if err := validateData(data); err != nil {
        return errors.ErrFileProcessing.WithWrap(err)
    }
    
    return nil
}

// Later in code
err := processFile("config.json")
if err != nil {
    // Check for specific error types in chain
    if errors.Is(err, errors.ErrFileNotFound) {
        // Handle file not found
    }
    
    // Get all errors in chain
    var cerr *errors.CError
    if errors.As(err, &cerr) {
        allWrapped := cerr.UnwrapAll()
        for _, wrappedErr := range allWrapped {
            fmt.Println("Wrapped:", wrappedErr.Error())
        }
    }
}
```

### Business Logic Validation

```go
func transferFunds(from, to string, amount decimal.Decimal) error {
    if amount.LessThanOrEqual(decimal.Zero) {
        return errors.ErrValidationRange.
            WithUserDetail("Transfer amount must be greater than zero")
    }
    
    balance, err := getBalance(from)
    if err != nil {
        return errors.ErrDatabaseError.WithWrap(err)
    }
    
    if balance.LessThan(amount) {
        return errors.ErrInsufficientFunds.
            WithDetail(fmt.Sprintf("balance: %s, requested: %s", balance, amount)).
            WithUserDetail("Insufficient funds for this transfer")
    }
    
    // Perform transfer...
    return nil
}
```

## Error Identity and Reusability

One of the key features of CError is that errors with the same ID are considered equal regardless of their current message or code. This allows for flexible error reuse:

```go
// Base error defined once
var ErrUserNotFound = errors.New(404, "user not found")

// Used in different contexts with different details
func getUserByEmail(email string) (*User, error) {
    // Same logical error, different context
    return nil, ErrUserNotFound.
        WithDetail(fmt.Sprintf("email: %s", email)).
        WithUserDetail("No user found with this email address")
}

func getUserByID(id string) (*User, error) {
    // Same logical error, different context
    return nil, ErrUserNotFound.
        WithDetail(fmt.Sprintf("user_id: %s", id)).
        WithUserDetail("User not found")
}

// Both errors are considered equal
err1 := getUserByEmail("test@example.com")
err2 := getUserByID("123")
fmt.Println(errors.Is(err1, err2)) // true - same logical error
```

## Advanced Features

### Full Error Data for Logging

```go
err := errors.ErrDatabaseError.
    WithWrap(originalErr).
    WithDetail("connection_pool: primary").
    WithDetail("query_timeout: 30s")

// Get complete error information including wrapped chain
fullData := err.FullData()
logger.Error("Database operation failed", "error_data", fullData)
```

### Custom Error Codes

```go
// Create errors with custom codes
err := errors.New(4001, "custom validation error")

// Change codes on existing errors
customErr := errors.ErrValidationInvalid.WithCode(4001)
```

## Best Practices

1. **Use predefined errors** when possible for consistency
2. **Add context with details** - technical details for logging, user details for APIs
3. **Wrap underlying errors** to maintain error chains
4. **Use consistent error codes** that align with HTTP status codes
5. **Separate technical and user-facing information**
6. **Log full error data** but return only user-safe information to clients

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.
