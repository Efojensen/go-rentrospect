package payments

import (
	"context"
	"fmt"
	"time"

	"github.com/EfoJensen/go-rentrospect/types"
	"github.com/jackc/pgx/v5"
)

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
		TotalBal:     totalBalance,
		EscrowBal:    escrowBalance,
		AvailableBal: totalBalance - escrowBalance,
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
