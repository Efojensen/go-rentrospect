package webhooks

import (
	"context"
	"fmt"
	"time"

	"github.com/EfoJensen/go-rentrospect/types"
)

func (h *WebHookHandler) recordCompletedPayment(code string, totalAmount int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)

	defer cancel()

	tx, err := h.store.Begin(ctx)

	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var total_balance int64
	var escrow_balance int64

	markForModification := `
		SELECT w.total_balance, w.escrow_balance FROM wallets AS w
		INNER JOIN payments AS p ON w.wallet_id = p.wallet_id
		INNER JOIN booking_codes AS b ON p.transaction_id = b.transaction_id
		WHERE b.code = $1 FOR UPDATE
	`
	err = tx.QueryRow(ctx, markForModification, code).Scan(&total_balance, &escrow_balance)

	if err != nil {
		return err
	}

	var transaction_id string
	updateBooking := `
		UPDATE booking_codes SET scanned_at = NOW(), status = $1
		WHERE code = $2 AND scanned_at IS NULL
		RETURNING transaction_id
	`

	err = tx.QueryRow(ctx, updateBooking, types.Scanned.String(), code).Scan(&transaction_id)

	if err != nil {
		return err
	}

	updateRentalTx := `
		UPDATE rental_transactions SET status = $1, escrow_status = $2, updated_at = NOW()
		WHERE transaction_id = $3
	`

	_, err = tx.Exec(ctx, updateRentalTx, types.Active.String(), types.Released.String(), transaction_id)

	if err != nil {
		return err
	}

	var wallet_id string
	updatePaymentQuery := `
		UPDATE payments SET status = $1, completed_at = NOW() WHERE transaction_id = $2 RETURNING wallet_id;
	`

	err = tx.QueryRow(ctx, updatePaymentQuery, types.Success.String(), transaction_id).Scan(&wallet_id)

	if err != nil {
		return err
	}

	updateWalletQuery := `UPDATE wallets SET escrow_balance = escrow_balance - $1,
		total_balance = total_balance - $1, updated_at = NOW()
		WHERE wallet_id = $2`

	_, err = tx.Exec(ctx, updateWalletQuery, totalAmount, wallet_id)

	if err != nil {
		return err
	}

	totalAfter := total_balance - totalAmount
	escrowAfter := escrow_balance - totalAmount
	availableAfter := totalAfter - escrowAfter

	updateWalletTxQuery := `INSERT INTO wallet_transactions (wallet_id, type, amount, related_transaction_id,
		total_after, available_after, escrow_after, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err = tx.Exec(ctx, updateWalletTxQuery, wallet_id, types.EscrowRelease.String(), totalAmount, transaction_id,
		totalAfter, availableAfter, escrowAfter, fmt.Sprintf("Escrow release of GH₵%d", totalAmount),
	)

	if err != nil {
		return err
	}

	err = tx.Commit(ctx)

	if err != nil {
		return err
	}

	return nil
}
