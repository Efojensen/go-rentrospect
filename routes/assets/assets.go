package assets

import (
	"net/http"

	"github.com/EfoJensen/go-rentrospect/middleware"
	"github.com/EfoJensen/go-rentrospect/upload"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AssetHandler struct {
	store       *pgxpool.Pool
	objectStore *upload.Storage
}

func NewAssetHandler(store *pgxpool.Pool, objStore *upload.Storage) *AssetHandler {
	return &AssetHandler{
		store: store,
		objectStore: objStore,
	}
}

func (a *AssetHandler) RegisterAssetRoutes(h *http.ServeMux) {
	h.HandleFunc("/assets", middleware.EnableCORS(a.UploadAsset))
}
