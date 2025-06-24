package main

import (
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sollunar/kvstore-api/configs"
	_ "github.com/sollunar/kvstore-api/docs"
	"github.com/sollunar/kvstore-api/internal/kvstore"
	"github.com/sollunar/kvstore-api/pkg/middleware"
	"github.com/sollunar/kvstore-api/pkg/storage"
	"github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func App(config *configs.Config) http.Handler {
	store := storage.NewTarantoolStorage(config.TT.Host, config.TT.Port)
	logger, _ := zap.NewProduction()

	kvRepository := kvstore.NewKVRepository(store, logger)
	kvservice := kvstore.NewKVService(kvRepository)

	apiRouter := http.NewServeMux()
	kvstore.NewKVStoreHandler(apiRouter, kvservice, logger)

	rootRouter := http.NewServeMux()
	rootRouter.Handle("/swagger/", httpSwagger.WrapHandler)
	rootRouter.HandleFunc("/health", healthHandler)
	rootRouter.Handle("/api/v1/", http.StripPrefix("/api/v1", middleware.Chain(
		middleware.Logging(logger),
	)(apiRouter)))

	return rootRouter
}

func main() {
	config := configs.Load()
	app := App(config)

	server := http.Server{
		Addr:           ":" + config.Server.Port,
		Handler:        app,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		log.Println("pprof listening on :6060")
		if err := http.ListenAndServe(":6060", nil); err != nil {
			log.Fatalf("pprof server failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Server running at :%s\n", config.Server.Port)
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed with: %v", err)
		}
	}()

	<-stop
	log.Println("Server shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")
}
