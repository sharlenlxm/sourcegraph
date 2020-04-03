package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/inconshreveable/log15"
	"github.com/sourcegraph/sourcegraph/cmd/precise-code-intel-bundle-manager/server"
	"github.com/sourcegraph/sourcegraph/internal/debugserver"
	"github.com/sourcegraph/sourcegraph/internal/env"
	"github.com/sourcegraph/sourcegraph/internal/trace/ot"
	"github.com/sourcegraph/sourcegraph/internal/tracer"
)

var (
	storageDir = env.Get("LSIF_STORAGE_ROOT", "/lsif-storage", "Root dir containing uploads and converted bundles.")
)

func main() {
	env.Lock()
	env.HandleHelpFlag()
	tracer.Init()

	if storageDir == "" {
		log.Fatal("precise-code-intel-bundle-manager: LSIF_STORAGE_ROOT is required")
	}
	for _, dir := range []string{"", "uploads", "dbs"} {
		if err := os.MkdirAll(filepath.Join(storageDir, dir), os.ModePerm); err != nil {
			log.Fatalf("failed to create LSIF_STORAGE_ROOT: %s", err)
		}
	}

	bundleManager := server.Server{
		StorageDir: storageDir,
	}
	// bundleManager.RegisterMetrics()

	handler := ot.Middleware(bundleManager.Handler())

	go debugserver.Start()

	// TODO - janitorial stuff

	port := "3187"
	host := ""
	if env.InsecureDev {
		host = "127.0.0.1"
	}
	addr := net.JoinHostPort(host, port)
	srv := &http.Server{Addr: addr, Handler: handler}
	log15.Info("precise-code-intel-bundle-manager: listening", "addr", srv.Addr)

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Listen for shutdown signals. When we receive one attempt to clean up,
	// but do an insta-shutdown if we receive more than one signal.
	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGHUP)
	<-c
	go func() {
		<-c
		os.Exit(0)
	}()

	// Stop accepting requests. In the future we should use graceful shutdown.
	srv.Close()
}
