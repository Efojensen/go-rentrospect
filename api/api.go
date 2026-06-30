package api

import (
	"fmt"
	"net/http"

	"github.com/EfoJensen/go-rentrospect/routes/assets"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ApiServer struct {
	Port string
	Db *pgxpool.Pool
}

func NewApiServer(port string, db *pgxpool.Pool) *ApiServer {
	return &ApiServer{
		Port: port,
		Db: db,
	}
}

func (server *ApiServer) Run() error {
	mux := http.NewServeMux()

	assetHandler := assets.NewAssetHandler(server.Db)
	assetHandler.RegisterAssetRoutes(mux)

	fmt.Println("Server running on http://localhost" + server.Port)
	return http.ListenAndServe(server.Port, mux)
}