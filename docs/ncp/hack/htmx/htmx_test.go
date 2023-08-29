package htmx

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func TestHTMX(t *testing.T) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	Start(r)
	http.ListenAndServe(":3000", r)
}
