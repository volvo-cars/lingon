package htmx

import (
	_ "embed"
	"ncp/bla"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/nats-io/jsm.go/api"
)

//go:embed schema.html
var schemaHTML string

func (s *Server) serveSchema(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := r.Context().Value(AuthContext).(UserInfo)
	if !ok {
		http.Error(w, "no auth context", http.StatusUnauthorized)
		return
	}
	accountID := chi.URLParam(r, "accountID")
	schemaID := chi.URLParam(r, "schemaID")
	reply, err := bla.SendAccountGetMsg(s.nc, bla.AccountGetMsg{
		ID: accountID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	schemaReply, err := bla.SendSchemaGetForAccountMsg(s.nc, accountID, bla.SchemaGetMsg{
		Key: schemaID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	schema := schemaReply.Schema
	consumerInfo, err := bla.ConsumerState(
		s.nc,
		accountID,
		schema.Spec.Name,
		strings.ReplaceAll(schema.Spec.Version, ".", "_"),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	type Data struct {
		UserInfo     UserInfo
		Account      bla.Account
		Schema       bla.Schema
		ConsumerInfo *api.ConsumerInfo
	}
	data := Data{
		UserInfo:     userInfo,
		Account:      reply.Account,
		Schema:       schemaReply.Schema,
		ConsumerInfo: consumerInfo,
	}
	s.renderHTML(w, schemaHTML, data)
}
