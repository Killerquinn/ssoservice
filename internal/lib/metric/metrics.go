package metric

import (
	"fmt"
	"log"
	"strconv"

	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
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

func CreateMetrics(addr, name string) (Metrics, error) {
	const op = "metric.CreateMetrics"

	var metr PrometheusMetrics

	metr.HitsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: name + "_total_hits",
		},
	)

	if err := prometheus.Register(metr.HitsTotal); err != nil {

		return nil, fmt.Errorf("%s:%w", op, err)
	}

	metr.Hits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: name + "_hits",
		},

		[]string{"status", "method", "path"},
	)
	if err := prometheus.Register(metr.Hits); err != nil {

		return nil, fmt.Errorf("%s:%w", op, err)
	}

	metr.Times = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: name + "_times",
		},

		[]string{"status", "method", "path"},
	)

	if err := prometheus.Register(metr.Times); err != nil {

		return nil, fmt.Errorf("%s:%w", op, err)
	}

	if err := prometheus.Register(collectors.NewBuildInfoCollector()); err != nil {

		return nil, fmt.Errorf("%s:%w", op, err)
	}

	go func() {
		router := echo.New()

		router.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

		if err := router.Start(addr); err != nil {

			log.Fatal(err)
		}
	}()

	return &metr, nil

}

func (metr *PrometheusMetrics) IncHits(status int, method, path string) {

	metr.HitsTotal.Inc()

	metr.Hits.WithLabelValues(strconv.Itoa(status), method, path).Inc()
}

func (metr *PrometheusMetrics) ObserveResponseTime(status int, method, path string, observeTime float64) {

	metr.Times.WithLabelValues(strconv.Itoa(status), method, path).Observe(observeTime)
}
