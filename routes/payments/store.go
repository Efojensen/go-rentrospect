package payments

import (
	"context"
	"time"

	"github.com/EfoJensen/go-rentrospect/types"
)

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
