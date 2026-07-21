package payments

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/EfoJensen/go-rentrospect/types"
	"github.com/jackc/pgx/v5"
)

func (p *PaymentHandler) storePaymentQuery(payDetails types.PaymentSession) error {
	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)

	defer cancel()

	tx, err := p.store.Begin(ctx)

	if err != nil {
		log.Println("failed to begin transaction: %w", err)
		return err
	}

	defer func(){
		tx.Rollback(ctx)
	}()

	// insertPaymentQuery := `
	// 	INSERT INTO payments
	// 	()
	// `

	return nil
}

func (p *PaymentHandler) checkAvailableBalQuery(userId string) (*types.ClientBal, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)

	defer cancel()

	checkBalQuery := `
		SELECT u.name, w.total_balance, w.escrow_balance FROM users AS u
		INNER JOIN wallets AS w ON u.user_id = w.user_id
		WHERE u.user_id = $1
	`

	var name string;
	var totalBalance int;
	var escrowBalance int;
	err := p.store.QueryRow(ctx, checkBalQuery, userId).Scan(&name, &totalBalance, &escrowBalance)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("no user found with this id")
		}
		return nil, err
	}

	var clientBal types.ClientBal = types.ClientBal{
		Name: name,
		AvailableBal: totalBalance - escrowBalance,
	}

	

	return &clientBal, nil
}