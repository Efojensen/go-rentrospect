package api

import (
	"log"
	"net/http"
	"os"

	"github.com/EfoJensen/go-rentrospect/routes/assets"
	"github.com/EfoJensen/go-rentrospect/routes/clients"
	"github.com/EfoJensen/go-rentrospect/routes/payments"
	"github.com/EfoJensen/go-rentrospect/routes/vendors"
	"github.com/EfoJensen/go-rentrospect/upload"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/identity"
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

	bucket := os.Getenv("OCI_BUCKET")

	config := common.DefaultConfigProvider()

	_, err := identity.NewIdentityClientWithConfigurationProvider(config)
	if err != nil {
		log.Fatal("Config error: ", err)
	} else {
		log.Println("config load status: ✅")
	}

	storage, err := upload.NewStorage(bucket, config)
	if err != nil {
		log.Fatal("Failed to initialize storage: ", err)
	} else {
		log.Println("object storage init status: ✅")
	}

	clientHandler := clients.NewClientHandler(server.Db)
	clientHandler.RegisterClientRoutes(mux)

	vendorHandler := vendors.NewVendorHandler(server.Db)
	vendorHandler.RegisterVendorRoutes(mux)

	paymentHandler := payments.NewPaymentHandler(server.Db)
	paymentHandler.RegisterPaymentRoutes(mux)

	assetHandler := assets.NewAssetHandler(server.Db, storage)
	assetHandler.RegisterAssetRoutes(mux)

	log.Println("server running on http://localhost" + server.Port + " ✅")
	return http.ListenAndServe(server.Port, mux)
}