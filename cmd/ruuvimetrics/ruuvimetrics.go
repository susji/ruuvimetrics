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
	v := config.VERBOSE
	flag.StringVar(&l, "listen", "localhost:9900", "Listen address")
	flag.StringVar(&ct, "content-type", ct, "Content-Type header in responses")
	flag.StringVar(&ep, "endpoint", ep, "HTTP endpoint for metrics")
	flag.StringVar(&mf, "metric-format", mf, "Format string for metric name generation")
	flag.BoolVar(&v, "verbose", v, "Verbose output")
	flag.Parse()
	log.Println("starting to listen at", l, "and verbosity is", v)
	http.HandleFunc(config.ENDPOINT, server.GenerateMetricsHandler(st, server.MetricsOptions{
		ContentType: ct,
		Endpoint:    ep,
		MetricFmt:   mf,
		Verbose:     v,
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
			if v {
				log.Printf("update[%s]\n", d.MAC)
			}
			st.Update(&d)
		}
		if err := s.Err(); err != nil {
			log.Println("reading input failed:", err)
		}
	}()
	wg.Wait()
}
