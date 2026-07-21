package payments

import (
	"encoding/json"
	"net/http"

	"github.com/EfoJensen/go-rentrospect/types"
	"github.com/EfoJensen/go-rentrospect/utils"
)

func (p *PaymentHandler) ReceiveClientPay(w http.ResponseWriter, r *http.Request) {
	var payment types.PaymentReq

	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}
}