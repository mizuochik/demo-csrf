package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync/atomic"
	"syscall"
)

const authenticatedUser = "authenticated-user"

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	resource := &atomic.Value{}
	resource.Store("null")

	s := &http.Server{
		Addr: net.JoinHostPort("", "8080"),
		Handler: http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			switch req.URL.Path {
			case "/login":
				http.SetCookie(rw, &http.Cookie{
					Name:  "user",
					Value: authenticatedUser,
				})
				rw.WriteHeader(http.StatusOK)
				fmt.Fprintln(rw, "Logged in")
			case "/resource":
				c, err := req.Cookie("user")
				if err != nil {
					log.Printf("error: %s", err)
					rw.WriteHeader(http.StatusUnauthorized)
					fmt.Fprintln(rw, "Unauthorized")
					return
				}
				if c.Value != authenticatedUser {
					rw.WriteHeader(http.StatusUnauthorized)
					fmt.Fprintln(rw, "Unauthorized")
					return
				}
				switch req.Method {
				case http.MethodGet:
					rw.WriteHeader(http.StatusOK)
					fmt.Fprintln(rw, resource.Load().(string))
				case http.MethodPost, http.MethodPut:
					rw.WriteHeader(http.StatusOK)
					b := &strings.Builder{}
					io.Copy(b, req.Body)
					resource.Store(b.String())
					fmt.Fprintln(rw, resource.Load().(string))
				}
			}
		}),
	}
	go func() {
		log.Printf("listening on %s", s.Addr)
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
	ctx := context.Background()
	<-quit
	if err := s.Shutdown(ctx); err != nil {
		panic(err)
	}
}
