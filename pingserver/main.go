package main

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PingServer struct {
	responseTime prometheus.Histogram // ping_response_time_ms
}

func (ps *PingServer) DelayedPing(rw http.ResponseWriter, _ *http.Request) {
	start := time.Now()
	time.Sleep(time.Millisecond)
	_, _ = rw.Write([]byte("PONG"))
	ps.responseTime.Observe(float64(time.Now().Sub(start).Nanoseconds()) / 1_000_000)
}

func main() {
	ps := PingServer{}
	prom := prometheus.NewRegistry()
	ps.responseTime = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "ping_response_time_ms",
		Buckets: []float64{0, 0.5, 1, 1.5, 2, 2.5, 3, 3.5, 4, 4.5, 5, 5.5, 6, 6.5, 7, 10, 20, 40, 80, 160},
	})
	prom.MustRegister(ps.responseTime)
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", ps.DelayedPing)
	mux.Handle("/metrics", promhttp.HandlerFor(prom, promhttp.HandlerOpts{Registry: prom}))

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
