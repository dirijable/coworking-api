package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type HTTPServer struct {
	config Config
	mux    *http.ServeMux
}

func NewHTTPServer(config Config, mux *http.ServeMux) *HTTPServer {
	return &HTTPServer{
		config: config,
		mux:    mux,
	}
}

func (s *HTTPServer) Run(ctx context.Context) error {
	server := &http.Server{
		Addr:         s.config.Addr,
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
		IdleTimeout:  s.config.IdleTimeout,
		Handler:      s.mux,
	}
	errorsChan := make(chan error, 1)
	go func() {
		defer close(errorsChan)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errorsChan <- err
		}
	}()

	select {
	case err := <-errorsChan:
		return fmt.Errorf("server listen and serve: %w", err)
	case <-ctx.Done():
		shutdownCtx, stop := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
		defer stop()
		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()
			return fmt.Errorf("shutdown HTTP server: %w", err)
		}
	}
	return nil
}
