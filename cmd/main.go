package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"jamascrorpJS/gwatch/internal"
	"jamascrorpJS/gwatch/internal/interactors"
	"jamascrorpJS/gwatch/internal/transport"
	cache2 "jamascrorpJS/gwatch/pkg/cache"
	config2 "jamascrorpJS/gwatch/pkg/config"
	redis2 "jamascrorpJS/gwatch/pkg/redis"
)

func main() {
	mux := http.NewServeMux()
	cache := cache2.New()
	config := config2.NewConfig()
	redis := redis2.New(context.Background(), config)
	coordinator := interactors.New(cache, redis)

	r := transport.NewRoutes(mux, coordinator)

	mux.HandleFunc("POST /api/v1/location", r.Save)
	mux.HandleFunc("GET /api/v1/location/{deviceId}", r.Get)

	s := internal.New(mux, config)
	go func() {
		if err := s.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("not start: %v", err)
		}
	}()
	log.Printf("started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown err: %v", err)
	}
}
