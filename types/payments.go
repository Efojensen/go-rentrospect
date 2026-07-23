package types

import "time"

type Split struct {
	Amount    int    `json:"amount"`
	Recipient string `json:"recipient"`
}

type IncomingPaymentReq struct {
	Amount           int                  `json:"amount"`
	UserId           string               `json:"userId"`
	EndDate          time.Time            `json:"endDate"`
	AssetId          string               `json:"assetId"`
	StartDate        time.Time            `json:"startDate"`
	ConsultationMode ConsultationModeEnum `json:"consultationMode"`
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

type PaymentStatusEnum int

const (
	Success PaymentStatusEnum = iota
	Pending
	Failed
)

func (p PaymentStatusEnum) String() string {
	return [...]string{"success", "pending", "failed"}[p]
}

type RentalStatusEnum int

const (
	PendingV2 RentalStatusEnum = iota
	Active
	Completed
	Cancelled
)

func (r RentalStatusEnum) String() string {
	return [...]string{"pending", "active", "completed", "cancelled"}[r]
}

type ConsultationModeEnum int

const (
	VideoCall ConsultationModeEnum = iota
	Meetup
	Chat
)

func (c ConsultationModeEnum) String() string {
	return [...]string{"video_call", "meetup", "chat"}[c]
}

type EscrowStatusEnum int

const (
	Holding EscrowStatusEnum = iota
	Released
	Refunded
)

func (e EscrowStatusEnum) String() string {
	return [...]string{"holding", "released", "refunded"}[e]
}

type BookingStatusEnum int

const (
	PendingV3 BookingStatusEnum = iota
	Scanned
	Expired
)

func (b BookingStatusEnum) String() string {
	return [...]string{"pending", "scanned", "expired"}[b]
}

type WalletTxnEnum int

const (
	TopUp WalletTxnEnum = iota
	Withdrawal
	EscrowHold
	EscrowRelease
	Refund
)

func (w WalletTxnEnum) String() string {
	return [...]string{"top_up", "withdrawal", "escrow_hold", "escrow_release", "refund"}[w]
}