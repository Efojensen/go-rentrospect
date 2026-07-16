package vendors

import (
	"net/http"

	"github.com/EfoJensen/go-rentrospect/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VendorHandler struct {
	store *pgxpool.Pool
}

func NewVendorHandler(store *pgxpool.Pool) *VendorHandler {
	return &VendorHandler{
		store,
	}
}

func (v *VendorHandler) RegisterVendorRoutes(h *http.ServeMux) {
	h.HandleFunc("/vendor/signUp", middleware.EnableCORS(v.VendorSignUp))
}