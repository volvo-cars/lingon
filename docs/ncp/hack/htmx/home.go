package htmx

import (
	_ "embed"
	"log/slog"
	"ncp/bla"
	"net/http"
)

//go:embed home.html
var homeHTML string

func (s *Server) serveHome(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := r.Context().Value(AuthContext).(UserInfo)
	if !ok {
		http.Error(w, "no auth context", http.StatusUnauthorized)
		return
	}
	reply, err := bla.SendAccountListMsg(s.nc, bla.AccountListMsg{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type Data struct {
		UserInfo UserInfo
		Accounts []bla.Account
	}
	data := Data{
		UserInfo: userInfo,
		Accounts: reply.Accounts,
	}
	s.renderHTML(w, homeHTML, data)
}

func (s *Server) postAccounts(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("account-name")
	reply, err := bla.SendAccountCreateMsg(s.nc, bla.AccountCreateMsg{
		Name: name,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	slog.Info("account created", "name", name, "reply", reply)
	w.Header().Add("HX-Redirect", "/")
}
