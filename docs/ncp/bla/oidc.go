package bla

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"golang.org/x/oauth2"
)

type OIDCActor struct {
	// Conn is a nats connection to the actor account.
	Conn *nats.Conn

	subs     []*nats.Subscription
	listener net.Listener
}

func (oa *OIDCActor) Close() {
	oa.listener.Close()
	for _, sub := range oa.subs {
		sub.Unsubscribe()
	}
}

func RegisterOIDCActor(ctx context.Context, actor *OIDCActor) error {
	name := "oidc"
	js, err := actor.Conn.JetStream()
	if err != nil {
		return fmt.Errorf("jetstream: %w", err)
	}
	bucket, err := js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket: name,
	})
	if err != nil {
		return fmt.Errorf("create key value: %w", err)
	}
	// Setup OIDC server
	provider, err := oidc.NewProvider(ctx, "http://localhost:9998/")
	if err != nil {
		return fmt.Errorf("new oidc provider: %w", err)
	}
	// Configure an OpenID Connect aware OAuth2 client.
	oauth2Config := oauth2.Config{
		ClientID:     "web",
		ClientSecret: "secret",
		RedirectURL:  "http://localhost:9999/auth/callback",

		// Discovery returns the OAuth2 endpoints.
		Endpoint: provider.Endpoint(),

		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
	}
	oidcConfig := &oidc.Config{
		ClientID: oauth2Config.ClientID,
	}
	verifier := provider.Verifier(oidcConfig)

	or := &oidcRoutes{
		ctx:      ctx,
		config:   &oauth2Config,
		verifier: verifier,
		bucket:   bucket,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Route("/", func(r chi.Router) {
		r.Use(or.authMiddleware)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello World!"))
		})
	})

	r.Get("/login", or.login)
	r.Get("/logout", or.logout)
	r.Get("/auth/callback", or.authCallback)
	l, err := net.Listen("tcp", ":9999")
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	actor.listener = l

	go func() {
		if err := http.Serve(l, r); err != nil {
			slog.Error("serve", "error", err)
		}
	}()

	return nil
}

type oidcRoutes struct {
	ctx      context.Context
	config   *oauth2.Config
	verifier *oidc.IDTokenVerifier
	bucket   nats.KeyValue
}

func (or *oidcRoutes) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := getSessionCookie(r)
		if err != nil {
			http.Error(w, "getting session cookie: "+err.Error(), http.StatusUnauthorized)
			return
		}
		kve, err := or.bucket.Get(sessionID.String())
		if err != nil {
			http.Error(w, "invalid session: "+err.Error(), http.StatusUnauthorized)
			return
		}
		slog.Info("MIDDLEWARE", "kve", string(kve.Value()))

		// TODO: add JWT to authorization header
		w.Header().Set("Authorization", "Bearer: "+"TODO")
		next.ServeHTTP(w, r)
	})
}

func (or *oidcRoutes) login(w http.ResponseWriter, req *http.Request) {
	state := uuid.New().String()
	nonce := uuid.New().String()
	// if p.devMode {
	// 	sessionID, err := p.store.NewSession(devUser)
	// 	if err != nil {
	// 		http.Error(w, "Creating new session: "+err.Error(), http.StatusInternalServerError)
	// 	}
	// 	writeSessionCookieHeader(w, sessionID)
	// 	http.Redirect(w, r, p.redirectURI, http.StatusSeeOther)
	// 	return
	// }
	// setCallbackCookie(w, r, "nonce", nonce)
	// setCallbackCookie(w, r, "state", state)
	http.Redirect(w, req, or.config.AuthCodeURL(state, oidc.Nonce(nonce)), http.StatusFound)
}

func (or *oidcRoutes) logout(w http.ResponseWriter, req *http.Request) {
	// sessionID, err := getSessionCookie(r)
	// if err != nil {
	// 	http.Error(w, "Invalid session: "+err.Error(), http.StatusUnauthorized)
	// 	return
	// }
	// if err := p.store.EndSession(sessionID); err != nil {
	// 	http.Error(w, "Cannot logout: "+err.Error(), http.StatusBadRequest)
	// 	return
	// }
	http.Redirect(w, req, "http://localhost:9999", http.StatusSeeOther)
}

func (or *oidcRoutes) authCallback(w http.ResponseWriter, req *http.Request) {
	// Verify the state
	state := req.FormValue("state")
	if state == "" {
		http.Error(w, "State not found", http.StatusBadRequest)
		return
	}
	expectedState := req.URL.Query().Get("state")
	if expectedState != state {
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}

	oauth2Token, err := or.config.Exchange(or.ctx, req.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, "exchange: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "missing token", http.StatusBadRequest)
		return
	}

	// Parse and verify ID Token payload.
	idToken, err := or.verifier.Verify(or.ctx, rawIDToken)
	if err != nil {
		http.Error(w, "token verification failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Extract custom claims
	var claims struct {
		Email    string `json:"email"`
		Verified bool   `json:"email_verified"`
	}
	if err := idToken.Claims(&claims); err != nil {
		// handle error
		slog.Error("claims: ", err)
	}
	slog.Info("claims", "claims", claims)
	sessionID := uuid.New()
	if _, err := or.bucket.Put(sessionID.String(), []byte("I am safe")); err != nil {
		http.Error(w, "cannot create session: "+err.Error(), http.StatusInternalServerError)
		return
	}
	writeSessionCookieHeader(w, sessionID)
	http.Redirect(w, req, "http://localhost:9999", http.StatusSeeOther)
}

const cookieName = "ncp_session"

func getSessionCookie(r *http.Request) (uuid.UUID, error) {
	sessionCookie, err := r.Cookie(cookieName)
	if err != nil {
		return uuid.UUID{}, err
	}
	return uuid.Parse(sessionCookie.Value)
}

func writeSessionCookieHeader(w http.ResponseWriter, sessionID uuid.UUID) {
	// Create cookie
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    sessionID.String(),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(time.Hour.Seconds() * 8),
		Path:     "/",
	})
}
