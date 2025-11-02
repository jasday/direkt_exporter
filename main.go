package main

import (
	"context"
	"errors"
	"flag"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	versioncollector "github.com/prometheus/client_golang/prometheus/collectors/version"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/vividbroadcast/prometheus-direkt-exporter/pkg/direkt"
)

func init() {
	prometheus.MustRegister(versioncollector.NewCollector("direkt_exporter"))
}

func main() {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)

	var dev bool
	flag.BoolVar(&dev, "development", false, "Whether to enable development mode")
	flag.BoolVar(&dev, "dev", false, "Whether to enable development mode")
	flag.BoolVar(&dev, "d", false, "Whether to enable development mode")
	flag.Parse()

	baseLogger := zerolog.New(os.Stderr)
	if dev {
		baseLogger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	}

	logger := baseLogger.With().Timestamp().Logger()
	logger.Info().Msg("Starting exporter")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	username := os.Getenv("DIREKT_USERNAME")
	password := os.Getenv("DIREKT_PASSWORD")
	if username == "" || password == "" {
		logger.Info().Str("username", username).Msg("Username or password not set, authentication will not be used")
	}
	d := direkt.New(username, password)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/-/healthy", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Healthy"))
	})
	mux.HandleFunc("/probe", func(w http.ResponseWriter, r *http.Request) {
		d.Handle(w, r, logger.With().Str("endpoint", "probe").Logger())
	})

	httpServer := &http.Server{
		Addr:        ":9110",
		Handler:     mux,
		BaseContext: func(_ net.Listener) context.Context { return ctx },
	}

	go func() {
		<-exit
		cancel()
		httpServer.Shutdown(context.Background())
	}()

	err := httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Err(err).Msg("Handler exited with error")
	}
	logger.Info().Msg("Exiting exporter")
}
