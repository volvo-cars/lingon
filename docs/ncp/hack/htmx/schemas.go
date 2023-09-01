package htmx

import (
	_ "embed"
	"ncp/bla"
	"net/http"

	"github.com/go-chi/chi/v5"
)

//go:embed schemas.html
var schemasHTML string

func (s *Server) serveSchemas(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("LOLOLOL"))
	// return
	userInfo, ok := r.Context().Value(AuthContext).(UserInfo)
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
	schemaListReply, err := bla.SendSchemaListForAccountMsg(
		s.nc,
		accountID,
		bla.SchemaListMsg{},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	type Data struct {
		UserInfo UserInfo
		Account  bla.Account
		Schemas  []bla.Schema
	}
	data := Data{
		UserInfo: userInfo,
		Account:  reply.Account,
		Schemas:  schemaListReply.Schemas,
	}
	s.renderTemplate(w, schemasHTML, data)
}
