package server

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	"github.com/numaproj/numaflow-go/pkg/info"
)

func getSDKVersion() string {
	version := "unknown"
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return version
	}
	for _, d := range info.Deps {
		if strings.HasSuffix(d.Path, "/numaflow-go") {
			version = d.Version
			break
		}
	}
	return version
}

func information(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	info := info.ServerInfo{Language: info.Go, Version: getSDKVersion(), Protocol: info.UDS}
	b, err := json.Marshal(&info)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
	}
}

// Start starts the info server, and waits for the context to be done.
func Start(ctx context.Context, opts ...Option) error {
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	if err := os.Remove(options.socketAdress); !os.IsNotExist(err) && err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		information(w, r)
	})

	server := &http.Server{
		Handler: mux,
	}
	listener, err := net.Listen("unix", options.socketAdress)
	if err != nil {
		return err
	}
	defer func() { _ = listener.Close() }()
	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
	log.Println("Info server is ready")
	// wait for signal
	<-ctx.Done()
	log.Println("Info server is now shutting down")
	defer log.Println("Info server has exited")

	// let's not wait indefinitely
	stopCtx, cancel := context.WithTimeout(context.Background(), options.drainTimeout)
	defer cancel()
	if err := server.Shutdown(stopCtx); err != nil {
		return err
	}
	return nil
}
