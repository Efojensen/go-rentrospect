package clients

import (
	"log"
	"net/http"

	"github.com/EfoJensen/go-rentrospect/cache"
	"github.com/EfoJensen/go-rentrospect/middleware"
	"github.com/EfoJensen/go-rentrospect/upload"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClientHandler struct {
	store       *pgxpool.Pool
	cache       *cache.UrlCache
	objectStore *upload.Storage
}

func NewClientHandler(store *pgxpool.Pool, objectStore *upload.Storage) *ClientHandler {
	return &ClientHandler{
		store: store,
		objectStore: objectStore,
	}
}

func (c *ClientHandler) RegisterClientRoutes(h *http.ServeMux) {
	cache, err := cache.NewUrlCache(c.objectStore)

	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("cache initialization status: ✅")
	}

	c.cache = cache

	h.HandleFunc("/client/signUp", middleware.EnableCORS(c.ClientSignUp))
	h.HandleFunc("/client/getAssets", middleware.EnableCORS(c.getAssets))
	h.HandleFunc("/client/getCategory", middleware.EnableCORS(c.getCategorizedAssets))
}
