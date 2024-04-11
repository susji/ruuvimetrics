package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/susji/ruuvi/data/rawv2"
	"github.com/susji/ruuvimetrics/internal/config"
	"github.com/susji/ruuvimetrics/internal/server"
	"github.com/susji/ruuvimetrics/internal/state"
)

func main() {
	st := state.New()
	l := ""
	ct := config.CT
	ep := config.ENDPOINT
	mf := config.METRICFMT
	lr := config.LOGREQUESTS
	li := config.LOGINPUT
	logger := log.Default()
	flag.StringVar(&l, "listen", "localhost:9900", "Listen address")
	flag.StringVar(&ct, "content-type", ct, "Content-Type header in responses")
	flag.StringVar(&ep, "endpoint", ep, "HTTP endpoint for metrics")
	flag.StringVar(&mf, "metric-format", mf, "Format string for metric name generation")
	flag.BoolVar(&lr, "log-requests", lr, "Log HTTP requests")
	flag.BoolVar(&li, "log-input", li, "Log successfully parsed input packets")
	flag.Parse()
	logger.Println("starting to listen at", l)
	h := server.GenerateMetricsHandler(st, server.MetricsOptions{
		ContentType: ct,
		Endpoint:    ep,
		MetricFmt:   mf,
	})
	if lr {
		h = server.GenerateRequestLogger(logger, h)
	}
	http.HandleFunc(config.ENDPOINT, h)
	wg := &sync.WaitGroup{}
	wg.Add(2)
	s := &http.Server{Addr: l}
	go func() {
		defer wg.Done()
		rc := 0
		err := s.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			logger.Println("HTTP server closed - reader probably shut down")
		} else if err != nil {
			logger.Println("HTTP server errored:", err)
			rc = 1
		}
		os.Exit(rc)
	}()
	go func() {
		defer wg.Done()
		defer s.Close()
		s := bufio.NewScanner(os.Stdin)
		logger.Println("reading standard input")
		for s.Scan() {
			d := rawv2.RuuviRawV2{}
			err := json.Unmarshal(s.Bytes(), &d)
			if err != nil {
				logger.Println("cannot unmarshal Ruuvi data:", err)
				continue
			}
			if li {
				logger.Printf("update[%s]\n", d.MAC)
			}
			st.Update(&d)
		}
		if err := s.Err(); err != nil {
			logger.Println("reading input failed:", err)
		}
	}()
	wg.Wait()
}
