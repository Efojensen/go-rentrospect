package httpClient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/EfoJensen/go-rentrospect/types"
	"github.com/EfoJensen/go-rentrospect/utils"
)

var httpClient = &http.Client{
	Timeout: time.Second * 5,
	Transport: &http.Transport{
		MaxIdleConnsPerHost: 20,
		MaxIdleConns:        100,
		IdleConnTimeout:     30 * time.Second,
	},
}

func MakeEscrowPayment(paymentReq types.IncomingPaymentReq) (*types.PaymentSessionRes, error) {
	payload := types.PaymentSession{
		Amount:    paymentReq.Amount,
		Description: fmt.Sprintf("escrow payment of GHS:%d", paymentReq.Amount),
	}

	payloadBytes, err := utils.GenerateIdempotencyKey(payload)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, os.Getenv("PAY_ESCROW_URL"), bytes.NewBuffer(payloadBytes))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("X-Api-Key", os.Getenv("X_API_KEY"))

	res, err := httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var paymentSessionRes types.PaymentSessionRes
	if err = json.NewDecoder(res.Body).Decode(&paymentSessionRes); err != nil {
		return nil, err
	}

	return &paymentSessionRes, nil
}
