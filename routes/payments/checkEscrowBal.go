package payments

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/EfoJensen/go-rentrospect/types"
	"github.com/EfoJensen/go-rentrospect/utils"
	"github.com/jackc/pgx/v5"
)

func (p *PaymentHandler) InitiatePayment(w http.ResponseWriter, r *http.Request) {
	var paymentReq types.PaymentReq

	if err := json.NewDecoder(r.Body).Decode(&paymentReq); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	defer r.Body.Close()
	clientBal, err := p.checkAvailableBalQuery(paymentReq.UserId)

	if err != nil {
		if err == pgx.ErrNoRows {
			utils.WriteErrorResponse(w, http.StatusNotFound, err)
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	if paymentReq.Amount > clientBal.AvailableBal {
		utils.WriteErrorResponse(w, http.StatusForbidden, errors.New("insufficient escrow funds"))
		return
	}

	utils.WriteResponse(w, http.StatusAccepted, clientBal)
}