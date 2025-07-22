package errors

// Validation errors (400 Bad Request)
var (
	ErrValidationRequired    = New(400, "field is required")
	ErrValidationInvalid     = New(400, "field value is invalid")
	ErrValidationFormat      = New(400, "field format is invalid")
	ErrValidationLength      = New(400, "field length is invalid")
	ErrValidationRange       = New(400, "field value out of range")
	ErrValidationEmail       = New(400, "invalid email format")
	ErrValidationPhone       = New(400, "invalid phone format")
	ErrValidationURL         = New(400, "invalid url format")
	ErrValidationPassword    = New(400, "password does not meet requirements")
	ErrValidationConfirm     = New(400, "confirmation does not match")
)

// Authentication errors
var (
	ErrAuthRequired          = New(401, "authentication required")
	ErrAuthInvalid           = New(401, "invalid credentials")
	ErrAuthExpired           = New(401, "authentication expired")
	ErrAuthTokenInvalid      = New(401, "invalid token")
	ErrAuthTokenExpired      = New(401, "token expired")
	ErrAuthPermissions       = New(403, "insufficient permissions")
	ErrAuthBlocked           = New(403, "account blocked")
	ErrAuthSuspended         = New(403, "account suspended")
)

// Not found errors (404 Not Found)
var (
	ErrNotFound             = New(404, "resource not found")
	ErrUserNotFound         = New(404, "user not found")
	ErrFileNotFound         = New(404, "file not found")
	ErrPageNotFound         = New(404, "page not found")
	ErrRecordNotFound       = New(404, "record not found")
	ErrEndpointNotFound     = New(404, "endpoint not found")
)

// Conflict errors (409 Conflict)
var (
	ErrConflict             = New(409, "resource conflict")
	ErrAlreadyExists        = New(409, "resource already exists")
	ErrUserExists           = New(409, "user already exists")
	ErrEmailTaken           = New(409, "email already taken")
	ErrDuplicateEntry       = New(409, "duplicate entry")
	ErrVersionConflict      = New(409, "version conflict")
)

// Business logic errors (422 Unprocessable Entity)
var (
	ErrUnprocessableEntity  = New(422, "unprocessable entity")
	ErrBusinessRule         = New(422, "business rule violation")
	ErrInsufficientFunds    = New(422, "insufficient funds")
	ErrOperationNotAllowed  = New(422, "operation not allowed")
	ErrLimitExceeded        = New(422, "limit exceeded")
	ErrExpiredResource      = New(422, "resource expired")
	ErrWorkflowError        = New(422, "workflow error")
)

// Rate limit errors (429 Too Many Requests)
var (
	ErrRateLimited          = New(429, "rate limit exceeded")
	ErrTooManyRequests      = New(429, "too many requests")
	ErrQuotaExceeded        = New(429, "quota exceeded")
	ErrAPILimitReached      = New(429, "api limit reached")
)

// File operation errors
var (
	ErrFileTooLarge         = New(413, "file too large")
	ErrFileFormatInvalid    = New(415, "invalid file format")
	ErrFileUploadFailed     = New(400, "file upload failed")
	ErrFileProcessing       = New(422, "file processing error")
	ErrStorageFull          = New(507, "storage full")
)

// Server errors
var (
	ErrInternalError        = New(500, "internal server error")
	ErrDatabaseError        = New(500, "database error")
	ErrTimeoutError         = New(504, "operation timeout")
	ErrServiceUnavailable   = New(503, "service unavailable")
	ErrMaintenanceMode      = New(503, "service under maintenance")
	ErrExternalService      = New(502, "external service error")
)

// Network errors
var (
	ErrNetworkError         = New(500, "network error")
	ErrConnectionRefused    = New(503, "connection refused")
	ErrDNSError             = New(502, "dns resolution error")
	ErrSSLError             = New(502, "ssl/tls error")
	ErrProxyError           = New(502, "proxy error")
)

// Parsing errors (400 Bad Request)
var (
	ErrParseError           = New(400, "parsing error")
	ErrJSONInvalid          = New(400, "invalid json")
	ErrXMLInvalid           = New(400, "invalid xml")
	ErrFormatUnsupported    = New(415, "unsupported format")
	ErrEncodingError        = New(400, "encoding error")
)
