package main

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
)

const (
	defaultListenAddress = "127.0.0.1"
	defaultMetricsPort   = 9696
)

type MetricsServer struct {
	metrics *telemetry.Metrics
	server  *http.Server
}

func (s *MetricsServer) metricsHandler(w http.ResponseWriter, _ *http.Request) {
	gr, err := s.metrics.Gather("prometheus")
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to gather metrics: %s", err))
		return
	}

	w.Header().Set("Content-Type", gr.ContentType)
	_, _ = w.Write(gr.Metrics)
}

func (s *MetricsServer) StartMetricsClient(config Config) {
	m, err := telemetry.New(telemetry.Config{
		ServiceName:             "loadtest-client",
		Enabled:                 true,
		EnableHostnameLabel:     true,
		EnableServiceLabel:      true,
		PrometheusRetentionTime: 600,
	})
	if err != nil {
		panic(err)
	}
	s.metrics = m
	http.HandleFunc("/healthz", s.healthzHandler)
	http.HandleFunc("/metrics", s.metricsHandler)

	metricsPort := config.MetricsPort
	if config.MetricsPort == 0 {
		metricsPort = defaultMetricsPort
	}

	listenAddr := fmt.Sprintf("%s:%d", defaultListenAddress, metricsPort)
	log.Printf("Listening for metrics scrapes on http://%s/metrics", listenAddr)

	s.server = &http.Server{
		Addr:              listenAddr,
		ReadHeaderTimeout: 3 * time.Second,
	}
	err = s.server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func (s *MetricsServer) healthzHandler(w http.ResponseWriter, _ *http.Request) {
	_, err := io.WriteString(w, "ok\n")
	if err != nil {
		panic(err)
	}
}

func WriteErrorResponse(w http.ResponseWriter, status int, err string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(legacy.Cdc.MustMarshalJSON(NewErrorResponse(0, err)))
}

// ErrorResponse defines the attributes of a JSON error response.
type ErrorResponse struct {
	Code  int    `json:"code,omitempty"`
	Error string `json:"error"`
}

// NewErrorResponse creates a new ErrorResponse instance.
func NewErrorResponse(code int, err string) ErrorResponse {
	return ErrorResponse{Code: code, Error: err}
}
