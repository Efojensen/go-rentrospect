package httpClient

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/EfoJensen/go-rentrospect/types"
)

var httpClient = &http.Client{
	Timeout: time.Second * 5,
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 20,
		IdleConnTimeout:     30 * time.Second,
	},
}

func SendPayment(vendorMail string, paymentReq types.IncomingPaymentReq) (*types.PaymentSession, error) {
	payload := types.PaymentSession{
		Amount:    paymentReq.Amount,
		Recipient: vendorMail,
	}

	payloadBytes, err := json.Marshal(&payload)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, os.Getenv("PAY_VENDOR_URL"), bytes.NewBuffer(payloadBytes))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("X-Api-Key", os.Getenv("X-API-KEY"))

	res, err := httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var paymentSession types.PaymentSession
	if err = json.NewDecoder(res.Body).Decode(&paymentSession); err != nil {
		return nil, err
	}

	return &paymentSession, nil
}
