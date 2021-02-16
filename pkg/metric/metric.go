package metric

import (
	"log"
	"strconv"

	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics interface {
	IncHits(status int, method, path string)
	ObserveResponseTime(status int, method, path string, observeTime float64)
}

type PrometheusMetrics struct {
	HitsTotal prometheus.Counter
	Hits      *prometheus.CounterVec
	Times     *prometheus.HistogramVec
}

func CreateMetrics(address string, name string) (Metrics, error) {
	var pm PrometheusMetrics
	pm.HitsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: name + "_hits_total",
	})
	if err := prometheus.Register(pm.HitsTotal); err != nil {
		return nil, err
	}

	pm.Hits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: name + "_hits",
		},
		[]string{"status", "method", "path"},
	)
	if err := prometheus.Register(pm.Hits); err != nil {
		return nil, err
	}

	pm.Times = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: name + "_times",
		},
		[]string{"status", "method", "path"},
	)
	if err := prometheus.Register(pm.Times); err != nil {
		return nil, err
	}

	if err := prometheus.Register(prometheus.NewBuildInfoCollector()); err != nil {
		return nil, err
	}

	grpcPrometheus.EnableHandlingTimeHistogram()
	go func() {
		router := echo.New()
		router.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
		if err := router.Start(address); err != nil {
			log.Fatal(err)
		}
	}()

	return &pm, nil
}

func (pm *PrometheusMetrics) IncHits(status int, method, path string) {
	pm.HitsTotal.Inc()
	pm.Hits.WithLabelValues(strconv.Itoa(status), method, path).Inc()
}

func (pm *PrometheusMetrics) ObserveResponseTime(status int, method, path string, observeTime float64) {
	pm.Times.WithLabelValues(strconv.Itoa(status), method, path).Observe(observeTime)
}
