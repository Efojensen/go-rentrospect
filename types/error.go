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

var errorCodes = map[ErrorCodeEnum]string{
    Forbidden:           "FORBIDDEN",
    NotFound:            "NOT_FOUND",
    BadRequest:          "BAD_REQUEST",
    BadGateway:          "BAD_GATEWAY",
    Unauthorized:        "UNAUTHORIZED",
    GatewayTimeout:      "GATEWAY_TIMEOUT",
    TooManyRequests:     "TOO_MANY_REQUESTS",
    MethodNotAllowed:    "METHOD_NOT_ALLOWED",
    ServiceUnavailable:  "SERVICE_UNAVAILABLE",
    InternalServerError: "INTERNAL_SERVER_ERROR",
}

func (e ErrorCodeEnum) ToString() string {
    if s, ok := errorCodes[e]; ok {
        return s
    }
    return "UNKNOWN_ERROR"
}
