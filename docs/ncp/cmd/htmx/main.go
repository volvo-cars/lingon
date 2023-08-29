package main

import (
	"context"
	"log/slog"
	"ncp/hack/htmx"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, htmx.AuthContext, htmx.UserInfo{
				Sub:    "123",
				Iss:    "http://localhost:9998/",
				Name:   "John Doe",
				Email:  "local@localhost",
				Groups: []string{"admin"},
			})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	if err := htmx.Start(r); err != nil {
		slog.Error("start htmx", "error", err)
		os.Exit(1)
	}
	http.ListenAndServe(":3000", r)
}
