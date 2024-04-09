package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/susji/ruuvi/data/rawv2"
	"github.com/susji/ruuvimetrics/internal/help"
	"github.com/susji/ruuvimetrics/internal/state"
)

var (
	STATE     = state.New()
	VERBOSE   = false
	CT        = "text/plain; version=0.0.4"
	ENDPOINT  = "/metrics"
	METRICFMT = "ruuvi_%s"
)

func dump[T comparable](w http.ResponseWriter, samples map[rawv2.MAC]state.Pair[T], kind, help string) {
	metric := fmt.Sprintf(METRICFMT, kind)
	fmt.Fprintf(w, "# HELP %s %s\n", metric, help)
	fmt.Fprintf(w, "# TYPE %s gauge\n", metric)
	for mac, v := range samples {
		fmt.Fprintf(w,
			"%s{mac=\"%s\"} %v %d\n",
			metric, mac.String(), v.Value, v.Timestamp.UnixMilli())
	}
}

func metrics(w http.ResponseWriter, r *http.Request) {
	if VERBOSE {
		log.Println(ENDPOINT)
	}
	// If performance ever becomes an issue, we could cache the metrics for
	// a configurable amount of time.
	w.Header().Add("content-type", CT)
	dump(w, STATE.Temperatures(), "temperature", help.TEMP)
	dump(w, STATE.Voltages(), "voltage", help.VOLT)
	dump(w, STATE.Humidities(), "humidity", help.HUM)
	dump(w, STATE.Pressures(), "pressure", help.PRES)
	dump(w, STATE.AccelerationXs(), "acceleration_x", help.ACCEL)
	dump(w, STATE.AccelerationYs(), "acceleration_y", help.ACCEL)
	dump(w, STATE.AccelerationZs(), "acceleration_z", help.ACCEL)
	dump(w, STATE.TransmitPowers(), "transmit_power", help.TX)
	dump(w, STATE.MovementCounters(), "movement_counter", help.MOV)
	dump(w, STATE.SequenceNumbers(), "sequence_number", help.SEQ)
}

func main() {
	var l string
	flag.StringVar(&l, "listen", "localhost:9900", "Listen address")
	flag.StringVar(&CT, "content-type", CT, "Content-Type header in responses")
	flag.StringVar(&ENDPOINT, "endpoint", ENDPOINT, "HTTP endpoint for metrics")
	flag.StringVar(&METRICFMT, "metric-format", METRICFMT, "Format string for metric name generation")
	flag.BoolVar(&VERBOSE, "verbose", false, "Verbose output")
	flag.Parse()
	log.Println("starting to listen at", l, "and verbosity is", VERBOSE)
	http.HandleFunc(ENDPOINT, metrics)
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
