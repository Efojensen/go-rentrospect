package renters

import (
	"net/http"

	"github.com/EfoJensen/go-rentrospect/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RenterHandler struct {
	store *pgxpool.Pool
}

func NewRenterHandler(store *pgxpool.Pool) *RenterHandler {
	return &RenterHandler{
		store: store,
	}
}

func (rr *RenterHandler) RegisterUserRoutes(h *http.ServeMux) {
	h.HandleFunc("/renter/signUp", middleware.EnableCORS(rr.RenterSignUp))
}