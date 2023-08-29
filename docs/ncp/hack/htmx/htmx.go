package htmx

import (
	"bytes"
	"fmt"
	"net/http"
	"text/template"

	"github.com/go-chi/chi/v5"
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

func Start(r *chi.Mux) error {
	page := template.New("html")
	if _, err := page.New("base").Parse(baseHTML); err != nil {
		return fmt.Errorf("parse base: %w", err)
	}
	if _, err := page.New("nav").Parse(navHTML); err != nil {
		return fmt.Errorf("parse nav: %w", err)
	}

	r.Get(
		"/",
		func(w http.ResponseWriter, r *http.Request) {
			type Account struct {
				ID string
			}
			type Data struct {
				UserInfo UserInfo
				Accounts []Account
			}
			userInfo, ok := r.Context().Value(AuthContext).(UserInfo)
			if !ok {
				http.Error(w, "no auth context", http.StatusUnauthorized)
				return
			}
			data := Data{
				UserInfo: userInfo,
				Accounts: []Account{
					{
						ID: "123",
					},
				},
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
			w.Write(buf.Bytes())
		},
	)
	r.Get("/accounts/new",
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
			w.Write(buf.Bytes())
		},
	)
	r.Post("/accounts", func(w http.ResponseWriter, r *http.Request) {
		// name := r.FormValue("account-name")
		// fmt.Fprintf(w, "account name: %s", name)
		w.Header().Add("HX-Redirect", "/")
	})

	r.Post("/clicked", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	return nil
}
