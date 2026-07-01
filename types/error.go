package types

type CustomError struct {
	Error      string `json:"error"`
	Message    *string `json:"message,omitempty"`
}
