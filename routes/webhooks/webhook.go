package webhooks

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type WebHookHandler struct {
	store *pgxpool.Pool
}

func NewWebhookHandler(store *pgxpool.Pool) *WebHookHandler {
	return &WebHookHandler{
		store,
	}
}

func (w *WebHookHandler) RegisterWebhookRoutes(h *http.ServeMux) {
	h.HandleFunc("/aza/completedPayment", w.CompletedPayment)
}