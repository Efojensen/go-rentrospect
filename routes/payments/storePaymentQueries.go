package payments

import (
	"context"
	"fmt"
	"time"

	"github.com/EfoJensen/go-rentrospect/types"
)

func (p *PaymentHandler) storeEscrowPaymentQueries(payReq types.IncomingPaymentReq,
	payDetails types.PaymentSessionRes,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)

	defer cancel()

	tx, err := p.store.Begin(ctx)

	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var walletId string
	var totalBalance int64
	var escrowBalance int64

	markForUpdateQuery := `SELECT total_balance, escrow_balance, wallet_id
		FROM wallets WHERE user_id = $1 FOR UPDATE`

	err = tx.QueryRow(ctx, markForUpdateQuery, payReq.UserId).Scan(&totalBalance, &escrowBalance, &walletId)

	if err != nil {
		return err
	}

	if totalBalance-escrowBalance < payReq.Amount {
		return fmt.Errorf("insufficient funds for escrow transaction")
	}

	updateWalletEscrowQuery := `
		UPDATE wallets SET escrow_balance = escrow_balance + $1
		WHERE user_id = $2
		RETURNING total_balance, escrow_balance
	`

	err = tx.QueryRow(ctx, updateWalletEscrowQuery, payReq.Amount, payReq.UserId).Scan(&totalBalance, &escrowBalance)

	if err != nil {
		return err
	}

	insertPaymentQuery := `
		INSERT INTO payments
		(wallet_id, aza_ref, amount, status)
		VALUES ($1, $2, $3, $4)
	`

	_, err = tx.Exec(ctx, insertPaymentQuery, walletId, payDetails.Data.Reference,
		payReq.Amount, types.Pending.String(),
	)

	if err != nil {
		return fmt.Errorf("insert payment: %w", err)
	}

	insertRentalQuery := `
		INSERT INTO rental_transactions
		(renter_id, asset_id, start_date, end_date, status, consultation_mode, price, escrow_status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING transaction_id
	`

	var rentalTransactionId string
	err = tx.QueryRow(ctx, insertRentalQuery, payReq.UserId, payReq.AssetId,
		payReq.StartDate, payReq.EndDate, types.PendingV2.String(), payReq.ConsultationMode.String(),
		payReq.Amount, types.Holding.String(),
	).Scan(&rentalTransactionId)

	if err != nil {
		return err
	}

	insertRentalLogQuery := `
		INSERT INTO wallet_transactions
		(wallet_id, type, amount, related_transaction_id, total_after, available_after, escrow_after, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	availableAfter := totalBalance - escrowBalance

	_, err = tx.Exec(ctx, insertRentalLogQuery, walletId, types.EscrowHold.String(),
		payReq.Amount, rentalTransactionId, totalBalance, availableAfter, escrowBalance,
		fmt.Sprintf("escrow payment of GHS:%d", payReq.Amount),
	)

	if err != nil {
		return fmt.Errorf("insert wallet_transactions: %w", err)
	}

	insertBookingCodeQuery := `INSERT INTO booking_codes (transaction_id, code, expires_at) VALUES ($1, $2, $3)`

	_, err = tx.Exec(ctx, insertBookingCodeQuery, rentalTransactionId, payDetails.Data.Id, payDetails.Data.ExpiresAt)

	if err != nil {
		return fmt.Errorf("insert booking_codes: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
