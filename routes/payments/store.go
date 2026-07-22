package payments

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/EfoJensen/go-rentrospect/types"
	"github.com/jackc/pgx/v5"
)

func (p *PaymentHandler) storePaymentQuery(
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

	var wallet_id string;
	err = tx.QueryRow(ctx, insertWalletQuery, clientBal.AvailableBal - payReq.Amount, clientBal.EscrowBal + payReq.Amount,payReq.UserId).Scan(&wallet_id)

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
	`

	conn, err = tx.Exec(ctx, insertRentalQuery, payReq.UserId, payReq.AssetId,
		payReq.StartDate, payReq.EndDate, types.PendingV2, payReq.ConsultationMode.String(),
		payReq.Amount, types.Holding.String(),
	)

	if err != nil {
		log.Println(conn.String())
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (p *PaymentHandler) checkAvailableBalQuery(userId string) (*types.ClientBal, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	checkBalQuery := `
	SELECT u.name, w.total_balance, w.escrow_balance FROM users AS u
		INNER JOIN wallets AS w ON u.user_id = w.user_id
		WHERE u.user_id = $1
	`

	var name string
	var totalBalance int
	var escrowBalance int
	err := p.store.QueryRow(ctx, checkBalQuery, userId).Scan(&name, &totalBalance, &escrowBalance)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("no user found with this id")
		}
		return nil, err
	}

	var clientBal types.ClientBal = types.ClientBal{
		Name:         name,
		AvailableBal: totalBalance - escrowBalance,
		EscrowBal: escrowBalance,
	}

	return &clientBal, nil
}

func (p *PaymentHandler) getVendorEmailFromAsset(paymentReq types.IncomingPaymentReq) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	getVendorQuery := `
		SELECT u.email FROM users AS u
		INNER JOIN assets AS a ON a.vendor_id = u.user_id
		WHERE a.vendor_id = $1
	`
	var vendorEmail string
	err := p.store.QueryRow(ctx, getVendorQuery, paymentReq.AssetId).Scan(&vendorEmail)

	if err != nil {
		return "", err
	}

	return vendorEmail, nil
}
