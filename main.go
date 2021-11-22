package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/dmitriivoitovich/test-assignment-sliide/app"
	"github.com/dmitriivoitovich/test-assignment-sliide/app/config"
	"github.com/dmitriivoitovich/test-assignment-sliide/app/provider"
)

var (
	addr = flag.String("addr", "127.0.0.1:8080", "the TCP address for the server to listen on, in the form 'host:port'")

	// app gets initialised with configuration.
	// as an example we've added 3 providers and a default configuration
	handler = app.App{
		ContentClients: map[provider.Provider]provider.Client{
			provider.Provider1: &provider.SampleContentProvider{Source: provider.Provider1},
			provider.Provider2: &provider.SampleContentProvider{Source: provider.Provider2},
			provider.Provider3: &provider.SampleContentProvider{Source: provider.Provider3},
		},
		Config: config.DefaultContentMix,
	}
)

func main() {
	flag.Parse()
	log.Printf("initalising server on %s", *addr)

	srv := http.Server{
		Addr:    *addr,
		Handler: handler,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}
