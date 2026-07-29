package payments

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/EfoJensen/go-rentrospect/types"
)

func (p *PaymentHandler) storePaymentQueries(
	clientBal types.ClientBal,
	payReq types.IncomingPaymentReq,
	payDetails types.PaymentSessionRes,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), 14*time.Second)

	defer cancel()

	tx, err := p.store.Begin(ctx)

	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		tx.Rollback(ctx)
	}()

	insertWalletQuery := `
		UPDATE wallets SET escrow_balance = escrow_balance - $1
		WHERE user_id = $2
		RETURNING wallet_id
	`

	var wallet_id string
	err = tx.QueryRow(ctx, insertWalletQuery, clientBal.EscrowBal+payReq.Amount,
		payReq.UserId).Scan(&wallet_id)

	if err != nil {
		return err
	}

	insertPaymentQuery := `
		INSERT INTO payments
		(wallet_id, aza_ref, amount, status)
		VALUES ($1, $2, $3, $4)
	`

	conn, err := tx.Exec(ctx, insertPaymentQuery, wallet_id, payDetails.Data.Reference,
		payReq.Amount, types.Pending.String(),
	)

	if err != nil {
		log.Println(conn.String())
		return err
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

	newEscrow := clientBal.EscrowBal + payReq.Amount
	availableAfter := clientBal.TotalBal - newEscrow

	conn, err = tx.Exec(ctx, insertRentalLogQuery,
		wallet_id, types.Holding.String(), payReq.Amount, rentalTransactionId,
		clientBal.TotalBal, availableAfter, newEscrow,
		fmt.Sprintf("escrow payment of GHS:%d", payReq.Amount),
	)

	if err != nil {
		log.Println(conn.String())
		return err
	}

	insertBookingCodeQuery := `INSERT INTO booking_codes (transaction_id) VALUES ($1)`

	conn, err = tx.Exec(ctx, insertBookingCodeQuery, rentalTransactionId)

	if err != nil {
		log.Println(conn.String())
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
