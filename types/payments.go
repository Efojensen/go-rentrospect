package types

type Split struct {
	Amount    int    `json:"amount"`
	Recipient string `json:"recipient"`
}

type PaymentReq struct {
	Amount  int    `json:"amount"`
	UserId  string `json:"userId"`
	AssetId string `json:"assetId"`
}

type PaymentSession struct {
	Amount    int    `json:"amount"`
	Recipient string `json:"recipient"`
	// Our own reference
	Reference string `json:"reference"`
	// Arbitrary JSON string, returned unchanged on the session and in webhooks
	Metadata string `json:"metadata"`
	// Description: Note shown to customer
	Description    string `json:"description"`
	IdempotencyKey string `json:"idempotencyKey"`
}
