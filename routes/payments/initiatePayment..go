package payments

import (
	"encoding/json"
	"net/http"

	"github.com/EfoJensen/go-rentrospect/httpClient"
	"github.com/EfoJensen/go-rentrospect/types"
	"github.com/EfoJensen/go-rentrospect/utils"
)

func (p *PaymentHandler) InitiatePayment(w http.ResponseWriter, r *http.Request) {
	var paymentReq types.IncomingPaymentReq

	if err := json.NewDecoder(r.Body).Decode(&paymentReq); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	defer r.Body.Close()

	paymentSession, err := httpClient.MakeEscrowPayment(paymentReq)

	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	err = p.storeEscrowPaymentQueries(paymentReq, *paymentSession)

	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteResponse(w, http.StatusAccepted, paymentSession)
}
