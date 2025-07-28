package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type HTTPServer struct {
	srv *http.Server
}

func New(addr string, handler http.Handler) *HTTPServer {
	return &HTTPServer{
		srv: &http.Server{
			Addr:              addr,
			Handler:           handler,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       30 * time.Second,
		},
	}
}

func (s *HTTPServer) Start() error {
	// graceful shutdown
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = s.srv.Shutdown(ctx)
	}()

	return s.srv.ListenAndServe()
}
