package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	ctx := context.Background()
	s := &http.Server{
		Addr:    net.JoinHostPort("", "9090"),
		Handler: http.FileServer(http.Dir("./static/attack-server")),
	}
	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
	<-quit
	if err := s.Shutdown(ctx); err != nil {
		panic(err)
	}
}
