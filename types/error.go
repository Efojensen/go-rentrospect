package types

type CustomError struct {
	Code    string  `json:"code"`
	Message *string `json:"message,omitempty"`
}

type ErrorCodeEnum int

const (
	NotFound            ErrorCodeEnum = 404
	Forbidden           ErrorCodeEnum = 403
	BadRequest          ErrorCodeEnum = 400
	BadGateway          ErrorCodeEnum = 502
	Unauthorized        ErrorCodeEnum = 401
	GatewayTimeout      ErrorCodeEnum = 504
	TooManyRequests     ErrorCodeEnum = 429
	MethodNotAllowed    ErrorCodeEnum = 405
	ServiceUnavailable  ErrorCodeEnum = 503
	InternalServerError ErrorCodeEnum = 500
)

func (e ErrorCodeEnum) ToString() string {
	return [...]string{
		"BAD_REQUEST", "UNAUTHORIZED", "FORBIDDEN", "NOT_FOUND", "METHOD_NOT_ALLOWED", "TOO_MANY_REQUESTS",
		"INTERNAL_SERVER_ERROR", "BAD_GATEWAY", "SERVICE_UNAVAILABLE", "GATEWAY_TIMEOUT",
	}[e]
}
