package htmx

import (
	"bytes"
	"fmt"
	"log/slog"
	"ncp/bla"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/nats-io/jsm.go/api"
	"github.com/nats-io/jwt/v2"
	"github.com/nats-io/nats.go"
)

// ContextKey is used to store context values
type ContextKey string

var AuthContext ContextKey = "AUTH_CONTEXT"

type UserInfo struct {
	Sub     string   `json:"sub"`
	Iss     string   `json:"iss"`
	Name    string   `json:"name"`
	Email   string   `json:"email"`
	Groups  []string `json:"groups"`
	Picture string   `json:"picture"`
}

func Setup(nc *nats.Conn, r *chi.Mux) error {
	page := template.New("html").Funcs(template.FuncMap{
		"durationSince": func(t time.Time) string {
			return time.Since(t).Truncate(time.Second).String()
		},
	})
	if _, err := page.New("base").Parse(baseHTML); err != nil {
		return fmt.Errorf("parse base: %w", err)
	}
	if _, err := page.New("nav").Parse(navHTML); err != nil {
		return fmt.Errorf("parse nav: %w", err)
	}

	r.Get(
		"/",
		func(w http.ResponseWriter, r *http.Request) {
			userInfo, ok := r.Context().Value(AuthContext).(UserInfo)
			if !ok {
				http.Error(w, "no auth context", http.StatusUnauthorized)
				return
			}
			reply, err := bla.SendAccountListMsg(nc, bla.AccountListMsg{})
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
			if _, err := page.New("home").Parse(homeHTML); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			buf := bytes.Buffer{}
			if err := page.ExecuteTemplate(&buf, "base", data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(buf.Bytes())
		},
	)
	r.Get(
		"/accounts/new",
		func(w http.ResponseWriter, r *http.Request) {
			if _, err := page.New("accounts_new").Parse(accountsNewHTML); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			buf := bytes.Buffer{}
			if err := page.ExecuteTemplate(&buf, "base", nil); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(buf.Bytes())
		},
	)
	r.Get(
		"/accounts/{accountID}",
		func(w http.ResponseWriter, r *http.Request) {
			userInfo, ok := r.Context().Value(AuthContext).(UserInfo)
			if !ok {
				http.Error(w, "no auth context", http.StatusUnauthorized)
				return
			}
			accountID := chi.URLParam(r, "accountID")
			reply, err := bla.SendAccountGetMsg(nc, bla.AccountGetMsg{
				ID: accountID,
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			schemaListReply, err := bla.SendSchemaListForAccountMsg(
				nc,
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
			// slog.Info("get account", "account", reply.Account, "data", data)
			if _, err := page.New("account").Parse(accountHTML); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			buf := bytes.Buffer{}
			if err := page.ExecuteTemplate(&buf, "base", data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(buf.Bytes())
		},
	)
	r.Get(
		"/accounts/{accountID}/schemas/{schemaID}",
		func(w http.ResponseWriter, r *http.Request) {
			userInfo, ok := r.Context().Value(AuthContext).(UserInfo)
			if !ok {
				http.Error(w, "no auth context", http.StatusUnauthorized)
				return
			}
			accountID := chi.URLParam(r, "accountID")
			schemaID := chi.URLParam(r, "schemaID")
			reply, err := bla.SendAccountGetMsg(nc, bla.AccountGetMsg{
				ID: accountID,
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			schemaReply, err := bla.SendSchemaGetForAccountMsg(nc, accountID, bla.SchemaGetMsg{
				Key: schemaID,
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			schema := schemaReply.Schema
			consumerInfo, err := bla.ConsumerState(
				nc,
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
			if _, err := page.New("schema").Parse(schemaHTML); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			buf := bytes.Buffer{}
			if err := page.ExecuteTemplate(&buf, "base", data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(buf.Bytes())
		},
	)
	r.Post("/accounts", func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("account-name")
		reply, err := bla.SendAccountCreateMsg(nc, bla.AccountCreateMsg{
			Name: name,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		slog.Info("account created", "name", name, "reply", reply)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Add("HX-Redirect", "/")
	})
	r.Post("/accounts/{accountID}/users", func(w http.ResponseWriter, r *http.Request) {
		accountID := chi.URLParam(r, "accountID")
		name := r.FormValue("user-name")
		reply, err := bla.SendUserCreateForAccountMsg(nc, accountID, bla.UserCreateMsg{
			Name: name,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		slog.Info("user created", "name", name, "reply", reply)
		creds, err := jwt.FormatUserConfig(reply.JWT, reply.NKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(creds)
	})

	r.Post("/clicked", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	return nil
}
