package assets

import (
	"net/http"

	"github.com/EfoJensen/go-rentrospect/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AssetHandler struct {
	store *pgxpool.Pool
}

func NewAssetHandler(store *pgxpool.Pool) *AssetHandler {
	return &AssetHandler{
		store: store,
	}
}

func (a *AssetHandler) RegisterAssetRoutes(h *http.ServeMux) {
	h.HandleFunc("/assets", middleware.EnableCORS(a.UploadAsset))
}