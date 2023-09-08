package htmx

import (
	"bytes"
	"fmt"
	"log/slog"
	"ncp/bla"
	"net/http"
	"text/template"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/nats-io/nats.go"
)

// ContextKey is used to store context values
type ContextKey string

var AuthContext ContextKey = "AUTH_CONTEXT"

var TemplateFuncMap = template.FuncMap{
	"durationSince": func(t time.Time) string {
		return time.Since(t).Truncate(time.Second).String()
	},
}

type Server struct {
	nc *nats.Conn
	// page *template.Template
}

func renderPage(
	w http.ResponseWriter,
	tmpl *template.Template,
	text string,
	data any,
) {
	if _, err := tmpl.Parse(text); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	buf := bytes.Buffer{}
	if err := tmpl.ExecuteTemplate(&buf, "base", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(buf.Bytes())
}

func renderTemplate(tmpl string, data any) ([]byte, error) {
	t, err := template.New("tmpl").Funcs(TemplateFuncMap).Parse(tmpl)
	if err != nil {
		return nil, fmt.Errorf("parse template: %w", err)
	}
	buf := bytes.Buffer{}
	if err := t.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("execute template: %w", err)
	}
	return buf.Bytes(), nil
}

func Setup(nc *nats.Conn, r *chi.Mux) error {

	s := Server{
		nc: nc,
		// page: tmpl,
	}

	r.Get("/", s.serveHome)
	r.Get("/accounts/new", s.serveAccountsNew)
	r.Get("/accounts/{accountID}", s.serveAccount)
	r.Get("/accounts/{accountID}/users", s.serveAccountUsers)
	r.Post("/accounts/{accountID}/users", s.postAccountUsers)
	// r.Get("/accounts/{accountID}/schemas/{schemaID}", s.serveSchema)
	r.Post("/accounts", s.postAccounts)

	r.Get("/accounts/{accountID}/{actor}", s.serveActor)
	r.Get("/accounts/{accountID}/{actor}/*", s.serveActor)

	// Actor endpoints triggered by hx-get
	r.Get(
		"/actor/{accountID}/{actor}",
		func(w http.ResponseWriter, r *http.Request) {
			// Create the request to send over NATS
			userInfo, ok := r.Context().Value(AuthContext).(bla.UserInfo)
			if !ok {
				http.Error(w, "no auth context", http.StatusUnauthorized)
				return
			}
			actor := chi.URLParam(r, "actor")
			accountID := chi.URLParam(r, "accountID")
			accReply, err := bla.SendAccountGetMsg(s.nc, bla.AccountGetMsg{
				ID: accountID,
			})
			if err != nil {
				slog.Error("get actor", "error", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			req := bla.Request{
				Account:    accReply.Account,
				UserInfo:   userInfo,
				Method:     r.Method,
				RequestURI: fmt.Sprintf("/%s", chi.URLParam(r, "*")),
			}

			reply, err := bla.RequestSubject[bla.Request, bla.Reply](
				nc,
				req,
				fmt.Sprintf(bla.ActorSubjectHTTPRender, actor),
			)
			if err != nil {
				slog.Error("actor http render", "error", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(reply.Body)
		},
	)
	r.Get(
		"/actor/{accountID}/{actor}/*",
		func(w http.ResponseWriter, r *http.Request) {
			// Create the request to send over NATS
			userInfo, ok := r.Context().Value(AuthContext).(bla.UserInfo)
			if !ok {
				http.Error(w, "no auth context", http.StatusUnauthorized)
				return
			}
			actor := chi.URLParam(r, "actor")
			accountID := chi.URLParam(r, "accountID")
			accReply, err := bla.SendAccountGetMsg(s.nc, bla.AccountGetMsg{
				ID: accountID,
			})
			if err != nil {
				slog.Error("get actor", "error", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			req := bla.Request{
				Account:    accReply.Account,
				UserInfo:   userInfo,
				Method:     r.Method,
				RequestURI: fmt.Sprintf("/%s", chi.URLParam(r, "*")),
			}

			reply, err := bla.RequestSubject[bla.Request, bla.Reply](
				nc,
				req,
				fmt.Sprintf(bla.ActorSubjectHTTPRender, actor),
			)
			if err != nil {
				slog.Error("actor http render", "error", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(reply.Body)
		},
	)

	return nil
}
