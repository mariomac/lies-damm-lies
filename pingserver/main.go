package main

import (
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	arg      = "msg"
	delayArg = "delay"
)

type Adder struct {
}

type Subtractor struct {
}

func (add *Adder) Mutate(output *int) {
	*output += 1
}

func (sub *Subtractor) Mutate(output *int) {
	*output -= 1
}

func doStruct(N int, value bool, output *int) {
	if value {
		mut := &Adder{}
		for i := 0; i < N; i++ {
			mut.Mutate(output)
		}
	} else {
		mut := &Subtractor{}
		for i := 0; i < N; i++ {
			mut.Mutate(output)
		}
	}
}

func pingHandler(rw http.ResponseWriter, req *http.Request) {
	slog.Debug("connection established", "remoteAddr", req.RemoteAddr)

	var value int
	for i := 0; i < rand.Intn(100)*10; i++ {
		doStruct(100099, true, &value)
		doStruct(49999, false, &value)
	}

	ret := "PONG!"
	if req.URL.Query().Has(arg) {
		ret = req.URL.Query().Get(arg)
	}

	if req.URL.Query().Has(delayArg) {
		delay, _ := time.ParseDuration(req.URL.Query().Get(delayArg))
		if delay > 0 {
			time.Sleep(delay)
		}
	}
	rw.WriteHeader(http.StatusOK)
	b, err := rw.Write([]byte(ret))
	if err != nil {
		slog.Error("writing response", err, "url", req.URL)
		return
	}
	slog.Debug(fmt.Sprintf("%T", rw))
	slog.Debug("written response", "url", req.URL, slog.Int("bytes", b))
}

type PingServer struct {
	responseTime prometheus.Histogram // ping_response_time_ms
}

func (ps *PingServer) DelayedPing(rw http.ResponseWriter, req *http.Request) {
	start := time.Now()
	pingHandler(rw, req)
	// spend around 5ms in CPU time
	//hash := fnv.New64()
	//fmt.Fprintf(hash, "%x", start)
	//i := byte(0)
	//runtime.Gosched()
	//for time.Now().Sub(start) < 5*time.Millisecond {
	//	hash.Write([]byte{i})
	//	i++
	//}
	//runtime.Gosched()
	//_, _ = rw.Write([]byte(fmt.Sprint(hash.Sum64())))
	ps.responseTime.Observe(time.Now().Sub(start).Seconds() * 1000)
}

func main() {
	ps := PingServer{}
	prom := prometheus.NewRegistry()
	ps.responseTime = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "ping_response_time_ms",
		Buckets: []float64{0, 0.5, 1, 1.5, 2, 2.5, 3, 3.5, 4, 4.5, 5, 5.5, 6, 6.5, 7, 10, 20, 40, 80, 160, 320, 640, 1280, 2560, 5120, 10240},
	})
	prom.MustRegister(ps.responseTime)
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", ps.DelayedPing)
	mux.Handle("/metrics", promhttp.HandlerFor(prom, promhttp.HandlerOpts{Registry: prom}))

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
