package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"ncp/bla"
	"ncp/hack/htmx"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	workdir := filepath.Join(os.TempDir(), "htmx")
	serverJWTAuthFile := filepath.Join(workdir, "server-jwt-auth.json")
	slog.Info("config", "workdir", workdir, "serverJWTAuthFile", serverJWTAuthFile)

	jwtAuth, err := loadOrBootstrapJWTAuth(serverJWTAuthFile)
	if err != nil {
		slog.Error("loading or bootstrapping jwt auth", "error", err)
		os.Exit(1)
	}

	ts, err := bla.NewTestServer(
		bla.WithDir(workdir),
		bla.WithJWTAuth(jwtAuth),
		// WithConfigureLogger(true),
	)
	if err != nil {
		slog.Error("new test server", "error", err)
		os.Exit(1)
	}
	if err := ts.StartUntilReady(); err != nil {
		slog.Error("start test server", "error", err)
		os.Exit(1)
	}
	slog.Info("nats server started", "url", ts.NS.ClientURL())

	if err := ts.PublishActorAccount(); err != nil {
		slog.Error("publish actor account", "error", err)
		os.Exit(1)
	}
	defer ts.NS.Shutdown()

	nc, err := ts.ActorUserConn()
	if err != nil {
		slog.Error("actor user connection", "error", err)
		os.Exit(1)
	}

	sysConn, err := ts.SysUserConn()
	if err != nil {
		slog.Error("sys user connection", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()
	// TODO: register actors
	aa := bla.AccountActor{
		OperatorNKey:       ts.Auth.OperatorNKey,
		SysAccountConn:     sysConn,
		ActorAccountConn:   nc,
		ActorAccountPubKey: ts.Auth.ActorAccountPublicKey,
	}
	if err := bla.RegisterAccountActor(ctx, &aa); err != nil {
		slog.Error("starting account actor", "error", err)
		os.Exit(1)
	}
	defer aa.Close()

	sa := bla.SchemaActor{
		Conn: nc,
	}
	if err := bla.RegisterSchemaActor(ctx, &sa); err != nil {
		slog.Error("starting schema actor", "error", err)
		os.Exit(1)
	}
	defer sa.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Register a dummy auth handler that cuts out the OIDC stuff and just creates
	// a user and adds it to the context (same as the OIDC handler would).
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, htmx.AuthContext, htmx.UserInfo{
				Sub:    "123",
				Iss:    "http://localhost:9998/",
				Name:   "John Doe",
				Email:  "local@localhost",
				Groups: []string{"admin"},
			})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	if err := htmx.Setup(nc, r); err != nil {
		slog.Error("start htmx", "error", err)
		os.Exit(1)
	}

	l, err := net.Listen("tcp", ":9999")
	if err != nil {
		slog.Error("listen", "error", err)
		os.Exit(1)
	}
	defer l.Close()

	slog.Info("starting server", "url", "http://localhost:9999/")

	if err := http.Serve(l, r); err != nil {
		slog.Error("serve", "error", err)
	}
}

func loadOrBootstrapJWTAuth(file string) (bla.ServerJWTAuth, error) {
	f, err := os.Open(file)
	if err != nil {
		if !os.IsNotExist(err) {
			return bla.ServerJWTAuth{}, fmt.Errorf("open: %w", err)
		}
		// If it doesn't exist, bootstrap the auth and save it
		jwtAuth, err := bla.BootstrapServerJWTAuth()
		if err != nil {
			return bla.ServerJWTAuth{}, fmt.Errorf("bootstrap: %w", err)
		}
		// Ensure the workdir directory exists
		if err := os.MkdirAll(filepath.Dir(file), os.ModePerm); err != nil {
			return bla.ServerJWTAuth{}, fmt.Errorf("mkdir: %w", err)
		}
		// Save the auth
		b, err := json.Marshal(jwtAuth)
		if err != nil {
			return bla.ServerJWTAuth{}, fmt.Errorf("marshal: %w", err)
		}
		if err := os.WriteFile(file, b, os.ModePerm); err != nil {
			return bla.ServerJWTAuth{}, fmt.Errorf("write file: %w", err)
		}
		return jwtAuth, nil
	}
	defer f.Close()
	var jwtAuth bla.ServerJWTAuth
	if err := json.NewDecoder(f).Decode(&jwtAuth); err != nil {
		return bla.ServerJWTAuth{}, fmt.Errorf("decode: %w", err)
	}
	if err := jwtAuth.LoadKeyPairs(); err != nil {
		return bla.ServerJWTAuth{}, fmt.Errorf("load key pairs: %w", err)
	}
	return jwtAuth, nil
}
