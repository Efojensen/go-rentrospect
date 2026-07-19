package clients

import (
	"net/http"

	"github.com/EfoJensen/go-rentrospect/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClientHandler struct {
	store *pgxpool.Pool
}

func NewClientHandler(store *pgxpool.Pool) *ClientHandler {
	return &ClientHandler{
		store: store,
	}
}

func (c *ClientHandler) RegisterUserRoutes(h *http.ServeMux) {
	h.HandleFunc("/client/signUp", middleware.EnableCORS(c.ClientSignUp))
}