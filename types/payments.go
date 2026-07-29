package types

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type IncomingPaymentReq struct {
	Amount           int64                `json:"amount"`
	UserId           string               `json:"userId"`
	EndDate          time.Time            `json:"endDate"`
	AssetId          string               `json:"assetId"`
	StartDate        time.Time            `json:"startDate"`
	ConsultationMode ConsultationModeEnum `json:"consultationMode"`
}

type PaymentSession struct {
	Amount int64 `json:"amount"`
	// Arbitrary JSON string, returned unchanged on the session and in webhooks
	Metadata string `json:"metadata"`
	// Our own reference
	Reference string `json:"reference"`
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
	Completed_
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

type PaymentSessionRes struct {
	Success bool `json:"success"`
	Data    struct {
		Id          string            `json:"id"`
		Amount      decimal.Decimal   `json:"amount"`
		Status      SessionStatusEnum `json:"status"`
		Currency    string            `json:"currency"`
		Reference   string            `json:"reference"`
		ExpiresAt   string            `json:"expiresAt"`
		CreatedAt   string            `json:"createdAt"`
		CheckoutUrl string            `json:"checkoutUrl"`
	} `json:"data"`
}

type SessionStatusEnum int

const (
	Pending_ SessionStatusEnum = iota
	Completed_V
	Expired_
	Cancelled_
	Refunded_
)

var sessionStatusMap = map[string]SessionStatusEnum{
	"PENDING":   Pending_,
	"EXPIRED":   Expired_,
	"REFUNDED":  Refunded_,
	"CANCELLED": Cancelled_,
	"COMPLETED": Completed_V,
}

func (ss *SessionStatusEnum) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	value, ok := sessionStatusMap[s]
	if !ok {
		return fmt.Errorf("invalid session status: %q", s)
	}

	*ss = value
	return nil
}

func (s SessionStatusEnum) String() string {
	return [...]string{"PENDING", "EXPIRED", "REFUNDED", "CANCELLED", "COMPLETED"}[s]
}
