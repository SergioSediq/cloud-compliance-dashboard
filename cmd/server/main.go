package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SergioSediq/security-portfolio-projects/cloud-compliance-dashboard/internal/api"
	"github.com/SergioSediq/security-portfolio-projects/cloud-compliance-dashboard/internal/config"
	"github.com/SergioSediq/security-portfolio-projects/cloud-compliance-dashboard/internal/server"
	"github.com/SergioSediq/security-portfolio-projects/cloud-compliance-dashboard/internal/store"
	"github.com/SergioSediq/security-portfolio-projects/cloud-compliance-dashboard/pkg/version"
)

//go:embed static/*
var staticFiles embed.FS

func main() {
	cfg := config.Load()
	st, err := store.Load(cfg.ChecksPath)
	if err != nil {
		log.Fatalf("load checks: %v", err)
	}

	h := &api.Handler{Store: st}
	sub, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}

	handler := server.New(h, sub)
	addr := cfg.ListenAddr()
	srv := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("listening on %s (checks=%s) version=%s commit=%s", addr, cfg.ChecksPath, version.Version, version.Commit)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}
