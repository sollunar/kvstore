package main

import (
	"log"
	"net/http"

	"github.com/sollunar/kvstore-api/configs"
	_ "github.com/sollunar/kvstore-api/docs"
	"github.com/sollunar/kvstore-api/internal/kvstore"
	"github.com/sollunar/kvstore-api/pkg/storage"
	"github.com/swaggo/http-swagger"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func App(config *configs.Config) http.Handler {
	store := storage.NewTarantoolStorage(config.TT.Host, config.TT.Port)
	router := http.NewServeMux()

	kvRepository := kvstore.NewKVRepository(store)
	kvservice := kvstore.NewKVService(kvRepository)
	kvstore.NewKVStoreHandler(router, kvservice)

	router.HandleFunc("/health", healthHandler)

	router.Handle("/swagger/", httpSwagger.WrapHandler)

	return router
}

func main() {
	config := configs.Load()
	app := App(config)

	server := http.Server{
		Addr:    ":" + config.Server.Port,
		Handler: app,
	}

	log.Printf("Server running at :%s\n", config.Server.Port)
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
	log.Println("Server stopped gracefully")
}
