package htmx

import (
	_ "embed"
	"net/http"
)

//go:embed accounts_new.html
var accountsNewHTML string

func (s *Server) serveAccountsNew(w http.ResponseWriter, r *http.Request) {
	s.renderPage(w, accountsNewHTML, nil)
}
