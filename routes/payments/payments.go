package payments

import (
	"net/http"

	"github.com/EfoJensen/go-rentrospect/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentHandler struct {
	store *pgxpool.Pool
}

func NewPaymentHandler(store *pgxpool.Pool) *PaymentHandler {
	return &PaymentHandler{
		store,
	}
}

func (p *PaymentHandler) RegisterPaymentRoutes(h *http.ServeMux) {
	h.HandleFunc("/payment", middleware.EnableCORS(p.ReceiveClientPay))
}