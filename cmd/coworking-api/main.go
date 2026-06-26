package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/dirijable/coworking-api/internal/core/http/server"
	"github.com/dirijable/coworking-api/internal/core/pgpool"
	"github.com/dirijable/coworking-api/internal/features/room/handler"
	"github.com/dirijable/coworking-api/internal/features/room/repository"
	"github.com/dirijable/coworking-api/internal/features/room/service"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	pool, err := pgpool.NewPGPool(ctx, pgpool.MustNewConfig())
	if err != nil {
		log.Fatal(err)
	}
	roomRepository := repository.NewPostgresRepository(pool)
	roomService := service.NewService(roomRepository)
	httpHandler := handler.NewRoomHTTPHandler(roomService)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/rooms", httpHandler.Create)
	mux.HandleFunc("GET /api/v1/rooms", httpHandler.FindAll)
	mux.HandleFunc("GET /api/v1/rooms/{id}", httpHandler.FindById)
	mux.HandleFunc("DELETE /api/v1/rooms/{id}", httpHandler.DeleteById)
	httpServer := server.NewHTTPServer(server.MustNewConfig("./config.yml"), mux)
	if err := httpServer.Run(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server run error: %v", err)
	}
}
