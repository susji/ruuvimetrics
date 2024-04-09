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
	"github.com/susji/ruuvimetrics/internal/server"
	"github.com/susji/ruuvimetrics/internal/state"
)

var (
	STATE     = state.New()
	VERBOSE   = false
	CT        = "text/plain; version=0.0.4"
	ENDPOINT  = "/metrics"
	METRICFMT = "ruuvi_%s"
)

func main() {
	var l string
	flag.StringVar(&l, "listen", "localhost:9900", "Listen address")
	flag.StringVar(&CT, "content-type", CT, "Content-Type header in responses")
	flag.StringVar(&ENDPOINT, "endpoint", ENDPOINT, "HTTP endpoint for metrics")
	flag.StringVar(&METRICFMT, "metric-format", METRICFMT, "Format string for metric name generation")
	flag.BoolVar(&VERBOSE, "verbose", false, "Verbose output")
	flag.Parse()
	log.Println("starting to listen at", l, "and verbosity is", VERBOSE)
	http.HandleFunc(ENDPOINT, server.GenerateMetricsHandler(STATE, server.MetricsOptions{
		ContentType: CT,
		Endpoint:    ENDPOINT,
		MetricFmt:   METRICFMT,
		Verbose:     VERBOSE,
	}))
	wg := &sync.WaitGroup{}
	wg.Add(2)
	s := &http.Server{Addr: l}
	go func() {
		defer wg.Done()
		err := s.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			log.Println("http server shutting down")
		} else if err != nil {
			log.Println("http server errored:", err)
			os.Exit(1)
		}
	}()
	go func() {
		defer wg.Done()
		defer s.Close()
		s := bufio.NewScanner(os.Stdin)
		log.Println("reading standard input")
		for s.Scan() {
			d := rawv2.RuuviRawV2{}
			err := json.Unmarshal(s.Bytes(), &d)
			if err != nil {
				log.Println("cannot unmarshal Ruuvi data:", err)
				continue
			}
			if VERBOSE {
				log.Printf("update[%s]\n", d.MAC)
			}
			STATE.Update(&d)
		}
		if err := s.Err(); err != nil {
			log.Println("reading input failed:", err)
		}
	}()
	wg.Wait()
}
