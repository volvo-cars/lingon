package main

import (
	"fmt"
	"ncp/oidcp"
	"net/http"
	"os"

	"log/slog"

	"github.com/zitadel/oidc/v2/example/server/exampleop"
	"github.com/zitadel/oidc/v2/example/server/storage"
)

func main() {
	port := "9998"
	issuer := fmt.Sprintf("http://localhost:%s/", port)

	storage := storage.NewStorage(oidcp.NewUserStore(oidcp.WithAdminUser()))
	router := exampleop.SetupServer(issuer, storage)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
	slog.Info("server listening", "url", issuer)
	slog.Info("press ctrl+c to stop")
	err := server.ListenAndServe()
	if err != nil {
		slog.Error("listen and serve", "error", err)
		os.Exit(1)
	}
}
