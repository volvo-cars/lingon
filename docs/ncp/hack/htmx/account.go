package htmx

import (
	_ "embed"
	"ncp/bla"
	"net/http"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/nats-io/jwt/v2"
)

//go:embed account.html
var accountHTML string

// const accountTmplName = "account"

var accountTmpl = template.Must(template.Must(baseTmpl.Clone()).Parse(accountHTML))

type ActorPage struct {
	Name string
	Link string
}

func (s *Server) serveAccount(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := r.Context().Value(AuthContext).(bla.UserInfo)
	if !ok {
		http.Error(w, "no auth context", http.StatusUnauthorized)
		return
	}
	accountID := chi.URLParam(r, "accountID")
	reply, err := bla.SendAccountGetMsg(s.nc, bla.AccountGetMsg{
		ID: accountID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	type Data struct {
		UserInfo bla.UserInfo
		Account  bla.Account
		Actors   []ActorPage
	}
	data := Data{
		UserInfo: userInfo,
		Account:  reply.Account,
		Actors: []ActorPage{
			{
				Name: "Schemas",
				Link: "schema",
			},
		},
	}
	renderPage(w, accountTmpl, accountDashboardHTML, data)
}

func (s *Server) serveAccountUsers(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := r.Context().Value(AuthContext).(bla.UserInfo)
	if !ok {
		http.Error(w, "no auth context", http.StatusUnauthorized)
		return
	}
	accountID := chi.URLParam(r, "accountID")
	reply, err := bla.SendAccountGetMsg(s.nc, bla.AccountGetMsg{
		ID: accountID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	type Data struct {
		UserInfo bla.UserInfo
		Account  bla.Account
		Actors   []ActorPage
	}
	data := Data{
		UserInfo: userInfo,
		Account:  reply.Account,
		Actors: []ActorPage{
			{
				Name: "Schemas",
				Link: "schema",
			},
		},
	}
	renderPage(w, accountTmpl, accountUsersHTML, data)
}

func (s *Server) postAccountUsers(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "accountID")
	name := r.FormValue("user-name")
	reply, err := bla.SendUserCreateForAccountMsg(s.nc, accountID, bla.UserCreateMsg{
		Name: name,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	creds, err := jwt.FormatUserConfig(reply.JWT, reply.NKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(creds)
}
