package bla_test

import (
	"context"
	"ncp/bla"
	"testing"
)

func TestOIDC(t *testing.T) {
	ctx := context.Background()
	ts := bla.StartTestServer(t)
	conn, err := ts.ActorUserConn()
	if err != nil {
		t.Fatal("getting actor user connection: ", err)
	}

	actor := bla.OIDCActor{
		Conn: conn,
	}
	if err := bla.RegisterOIDCActor(ctx, &actor); err != nil {
		t.Fatal("starting account actor: ", err)
	}
	defer actor.Close()

	ts.NS.WaitForShutdown()
}

// func TestOIDC(t *testing.T) {
// 	ctx := context.TODO()
// 	provider, err := oidc.NewProvider(ctx, "http://localhost:9998/")
// 	if err != nil {
// 		t.Fatal("getting provider: ", err)
// 	}

// 	// Configure an OpenID Connect aware OAuth2 client.
// 	oauth2Config := oauth2.Config{
// 		ClientID:     "web",
// 		ClientSecret: "secret",
// 		RedirectURL:  "http://localhost:9999/auth/callback",

// 		// Discovery returns the OAuth2 endpoints.
// 		Endpoint: provider.Endpoint(),

// 		// "openid" is a required scope for OpenID Connect flows.
// 		Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
// 	}
// 	oidcConfig := &oidc.Config{
// 		ClientID: oauth2Config.ClientID,
// 	}
// 	verifier := provider.Verifier(oidcConfig)

// 	r := chi.NewRouter()
// 	r.Use(middleware.Logger)
// 	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
// 		w.Write([]byte("Hello World!"))
// 	})
// 	r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
// 		handleLogin(&oauth2Config, w, r)
// 	})
// 	r.Get("/logout", func(w http.ResponseWriter, r *http.Request) {
// 		handleLogout(w, r)
// 	})
// 	r.Get("/auth/callback", func(w http.ResponseWriter, r *http.Request) {
// 		handleAuthCallback(ctx, &oauth2Config, verifier, w, r)
// 	})
// 	http.ListenAndServe(":9999", r)
// }

// func handleLogin(config *oauth2.Config, w http.ResponseWriter, r *http.Request) {
// 	state := uuid.New().String()
// 	nonce := uuid.New().String()
// 	// if p.devMode {
// 	// 	sessionID, err := p.store.NewSession(devUser)
// 	// 	if err != nil {
// 	// 		http.Error(w, "Creating new session: "+err.Error(), http.StatusInternalServerError)
// 	// 	}
// 	// 	writeSessionCookieHeader(w, sessionID)
// 	// 	http.Redirect(w, r, p.redirectURI, http.StatusSeeOther)
// 	// 	return
// 	// }
// 	// setCallbackCookie(w, r, "nonce", nonce)
// 	// setCallbackCookie(w, r, "state", state)
// 	http.Redirect(w, r, config.AuthCodeURL(state, oidc.Nonce(nonce)), http.StatusFound)
// }

// func handleLogout(w http.ResponseWriter, r *http.Request) {
// 	// sessionID, err := getSessionCookie(r)
// 	// if err != nil {
// 	// 	http.Error(w, "Invalid session: "+err.Error(), http.StatusUnauthorized)
// 	// 	return
// 	// }
// 	// if err := p.store.EndSession(sessionID); err != nil {
// 	// 	http.Error(w, "Cannot logout: "+err.Error(), http.StatusBadRequest)
// 	// 	return
// 	// }
// 	http.Redirect(w, r, "http://localhost:9999", http.StatusSeeOther)
// }

// func handleAuthCallback(
// 	ctx context.Context,
// 	config *oauth2.Config,
// 	verifier *oidc.IDTokenVerifier,
// 	w http.ResponseWriter,
// 	r *http.Request,
// ) {
// 	// Verify the state
// 	state := r.FormValue("state")
// 	if state == "" {
// 		http.Error(w, "State not found", http.StatusBadRequest)
// 		return
// 	}
// 	expectedState := r.URL.Query().Get("state")
// 	if expectedState != state {
// 		http.Error(w, "Invalid state", http.StatusBadRequest)
// 		return
// 	}

// 	oauth2Token, err := config.Exchange(ctx, r.URL.Query().Get("code"))
// 	if err != nil {
// 		http.Error(w, "exchange: "+err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	// Extract the ID Token from OAuth2 token.
// 	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
// 	if !ok {
// 		http.Error(w, "missing token", http.StatusBadRequest)
// 		return
// 	}

// 	// Parse and verify ID Token payload.
// 	idToken, err := verifier.Verify(ctx, rawIDToken)
// 	if err != nil {
// 		http.Error(w, "token verification failed: "+err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	// Extract custom claims
// 	var claims struct {
// 		Email    string `json:"email"`
// 		Verified bool   `json:"email_verified"`
// 	}
// 	if err := idToken.Claims(&claims); err != nil {
// 		// handle error
// 		slog.Error("claims: ", err)
// 	}
// 	slog.Info("claims", "claims", claims)
// 	http.Redirect(w, r, "http://localhost:9999", http.StatusSeeOther)
// }
