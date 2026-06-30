package types

type Errors struct {
	Error      string `json:"error"`
	Message    *string `json:"message,omitempty"`
}
