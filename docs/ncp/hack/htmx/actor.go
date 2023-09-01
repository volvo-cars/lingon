package htmx

import (
	_ "embed"
	"ncp/bla"
	"net/http"

	"github.com/go-chi/chi/v5"
)

//go:embed actor.html
var actorHTML string

func (s *Server) serveActor(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := r.Context().Value(AuthContext).(UserInfo)
	if !ok {
		http.Error(w, "no auth context", http.StatusUnauthorized)
		return
	}
	accountID := chi.URLParam(r, "accountID")
	actor := chi.URLParam(r, "actor")
	subPath := chi.URLParam(r, "*")

	reply, err := bla.SendAccountGetMsg(s.nc, bla.AccountGetMsg{
		ID: accountID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	type Data struct {
		UserInfo UserInfo
		Account  bla.Account
		Actors   []ActorPage
		Actor    string
		Path     string
	}
	data := Data{
		UserInfo: userInfo,
		Account:  reply.Account,
		Actor:    actor,
		Path:     subPath,
		Actors: []ActorPage{
			{
				Name: "Schemas",
				Link: "schema",
			},
		},
	}
	s.renderPage(w, actorHTML, data)
}
