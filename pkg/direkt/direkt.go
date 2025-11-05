package direkt

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
)

const (
	url          = "https://iss.intinor.se/"
	unitEndpoint = "api/v1/units/"
)

var errUnitOffline = errors.New("unit offline")

func New(username, password string) *Direkt {
	return &Direkt{
		username: username,
		password: password,
		client: http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

type Direkt struct {
	username string
	password string
	client   http.Client
}

func (d *Direkt) Handle(w http.ResponseWriter, r *http.Request, l zerolog.Logger) {
	id, err := validateRequest(r)
	if err != nil {
		w.Write([]byte(err.Error()))
		// w.WriteHeader(http.StatusBadRequest)
		l.Err(err).Msg("Error validating request parameters")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	registry, err := d.gatherMetrics(ctx, l.With().Str("serial", id).Logger(), []metricGatherer{system, interfaces, decoders, outputs, encoders}, id)
	if err != nil {
		w.Write([]byte(err.Error()))
		// w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}

func (d *Direkt) doRequest(l zerolog.Logger, req *http.Request) ([]byte, error) {
	if d.username != "" && d.password != "" {
		l.Debug().Msg("Authentication set")
		req.SetBasicAuth(d.username, d.password)
	}
	l.Trace().Str("url", req.URL.String()).Msg("Sending request")
	res, err := d.client.Do(req)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, nil
	}
	defer res.Body.Close()

	if res.StatusCode == 503 {
		return nil, errUnitOffline
	}

	if res.StatusCode != 200 {
		l.Info().Int("status_code", res.StatusCode).Str("request", req.URL.String()).Msg("Non-OK status code returned")
		return nil, errors.New("non-okay request returned")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	l.Trace().Str("url", req.URL.String()).Msg("Finished request, returning body")
	return body, err
}

type metricGatherer func(ctx context.Context, l zerolog.Logger, registry prometheus.Registerer, doReq func(l zerolog.Logger, request *http.Request) ([]byte, error), id string) error

func (d *Direkt) gatherMetrics(ctx context.Context, l zerolog.Logger, metrics []metricGatherer, id string) (*prometheus.Registry, error) {
	l.Info().Msg("Requesting metrics for Direkt unit")
	start := time.Now()
	successGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "request_success",
		Help: "Displays whether or not the request was a success",
	})
	durationGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "request_duration_seconds",
		Help: "Returns how long the request took to complete in seconds",
	})

	baseRegistry := prometheus.NewRegistry()
	registry := prometheus.WrapRegistererWith(prometheus.Labels{"serial": id}, baseRegistry)
	registry.MustRegister(successGauge)
	registry.MustRegister(durationGauge)
	var retErr error
	for _, metric := range metrics {
		err := metric(ctx, l, registry, d.doRequest, id)
		if err != nil {
			l.Err(err).Msg("Error retrieving metrics")
			successGauge.Set(0)
			if errors.Is(err, errUnitOffline) {
				break
			}
			retErr = err
		}
	}
	if retErr == nil {
		successGauge.Set(1)
	}
	duration := time.Since(start).Seconds()
	durationGauge.Set(duration)
	l.Info().Float64("duration", duration).Msg("Finished gathering metrics")
	return baseRegistry, retErr
}

func validateRequest(r *http.Request) (string, error) {
	params := r.URL.Query()

	val := params.Get("serial")
	if val == "" {
		return "", errors.New("no serial provided")
	}

	if !strings.HasPrefix(val, "D0") {
		return "", errors.New("invalid serial provided")
	}
	return val, nil
}
