package webhooks

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/EfoJensen/go-rentrospect/types"
	"github.com/EfoJensen/go-rentrospect/utils"
	"github.com/jackc/pgx/v5"
)

func (h *WebHookHandler) CompletedPayment(w http.ResponseWriter, r *http.Request) {
	var payload types.CheckoutCompletedEvent

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	err := h.recordCompletedPayment(payload.SessionID, int64(payload.Amount))

	if err != nil {
		if err == pgx.ErrNoRows {
			log.Println("payment already processed")
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
}