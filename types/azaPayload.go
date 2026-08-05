package types

import "time"

type Split struct {
	Amount    float64 `json:"amount"`
	Status    string  `json:"status"`
	Recipient string  `json:"recipient"`
}

type CheckoutCompletedEvent struct {
	Event       string    `json:"event"`
	Amount      float64   `json:"amount"`
	LiveMode    bool      `json:"livemode"`
	Metadata    string    `json:"metadata"`
	Currency    string    `json:"currency"`
	NetAmount   float64   `json:"netAmount"`
	Reference   string    `json:"reference"`
	OccurredAt  time.Time `json:"occurredAt"`
	SessionID   string    `json:"sessionId"`
	MerchantID  string    `json:"merchantId"`
	PlatformFee float64   `json:"platformFee"`
	Description string    `json:"description"`
	CompletedAt time.Time `json:"completedAt"`
	Splits      []Split   `json:"splits,omitempty"`
}
