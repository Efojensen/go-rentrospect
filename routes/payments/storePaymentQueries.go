package payments

import (
	"context"
	"log"
	"time"

	"github.com/EfoJensen/go-rentrospect/types"
)

func (p *PaymentHandler) storePaymentQueries(
	clientBal types.ClientBal,
	payReq types.IncomingPaymentReq,
	payDetails types.PaymentSession,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)

	defer cancel()

	tx, err := p.store.Begin(ctx)

	if err != nil {
		log.Println("failed to begin transaction: %w", err)
		return err
	}

	defer func() {
		tx.Rollback(ctx)
	}()

	insertWalletQuery := `
		UPDATE wallets SET total_balance = $1, escrow_balance= $2
		WHERE user_id = $3
		RETURNING wallet_id
	`

	var wallet_id string
	err = tx.QueryRow(ctx, insertWalletQuery, clientBal.AvailableBal-payReq.Amount,
		clientBal.EscrowBal+payReq.Amount, payReq.UserId).Scan(&wallet_id)

	if err != nil {
		return err
	}

	insertPaymentQuery := `
		INSERT INTO payments
		(wallet_id, aza_ref, amount, status)
		VALUES ($1, $2, $3, $4)
	`

	conn, err := tx.Exec(ctx, insertPaymentQuery, wallet_id, payDetails.Reference, payReq.Amount, types.Success.String())

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
		payReq.StartDate, payReq.EndDate, types.PendingV2, payReq.ConsultationMode.String(),
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

	conn, err = tx.Exec(ctx, insertRentalLogQuery,
		wallet_id, types.TopUp.String(), payReq.Amount, rentalTransactionId, clientBal.TotalBal,
		clientBal.TotalBal-payReq.Amount, clientBal.EscrowBal+payReq.Amount,
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
