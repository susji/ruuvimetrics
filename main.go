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
	"github.com/susji/ruuvimetrics/internal/state"
)

var (
	STATE   = state.New()
	VERBOSE = false
)

func metrics(w http.ResponseWriter, r *http.Request) {
	if VERBOSE {
		log.Println("/metrics")
	}
	w.Header().Add("content-type", "text/plain; version=0.0.4")
	fmt.Fprintln(w, "# HELP ruuvi_temperature", HELP_TEMP)
	fmt.Fprintln(w, "# TYPE ruuvi_temperature gauge")
	for mac, v := range STATE.Temperatures() {
		fmt.Fprintf(w,
			"ruuvi_temperature{mac=\"%s\"} %f %d\n",
			mac.String(), v.Value, v.Timestamp.UnixMilli())
	}
	fmt.Fprintln(w, "# HELP ruuvi_voltage", HELP_VOLT)
	fmt.Fprintln(w, "# TYPE ruuvi_voltage gauge")
	for mac, v := range STATE.Voltages() {
		fmt.Fprintf(w,
			"ruuvi_voltage{mac=\"%s\"} %f %d\n",
			mac.String(), v.Value, v.Timestamp.UnixMilli())
	}
	fmt.Fprintln(w, "# HELP ruuvi_humidity", HELP_HUM)
	fmt.Fprintln(w, "# TYPE ruuvi_humidity gauge")
	for mac, v := range STATE.Humidities() {
		fmt.Fprintf(w,
			"ruuvi_humidity{mac=\"%s\"} %f %d\n",
			mac.String(), v.Value, v.Timestamp.UnixMilli())
	}
	fmt.Fprintln(w, "# HELP ruuvi_pressure", HELP_PRES)
	fmt.Fprintln(w, "# TYPE ruuvi_pressure gauge")
	for mac, v := range STATE.Pressures() {
		fmt.Fprintf(w,
			"ruuvi_pressure{mac=\"%s\"} %d %d\n",
			mac.String(), v.Value, v.Timestamp.UnixMilli())
	}
	fmt.Fprintln(w, "# HELP ruuvi_acceleration_x", HELP_ACCEL)
	fmt.Fprintln(w, "# TYPE ruuvi_acceleration_x gauge")
	for mac, v := range STATE.AccelerationXs() {
		fmt.Fprintf(w,
			"ruuvi_acceleration_x{mac=\"%s\"} %d %d\n",
			mac.String(), v.Value, v.Timestamp.UnixMilli())
	}
	fmt.Fprintln(w, "# HELP ruuvi_acceleration_y", HELP_ACCEL)
	fmt.Fprintln(w, "# TYPE ruuvi_acceleration_y gauge")
	for mac, v := range STATE.AccelerationYs() {
		fmt.Fprintf(w,
			"ruuvi_acceleration_y{mac=\"%s\"} %d %d\n",
			mac.String(), v.Value, v.Timestamp.UnixMilli())
	}
	fmt.Fprintln(w, "# HELP ruuvi_acceleration_z", HELP_ACCEL)
	fmt.Fprintln(w, "# TYPE ruuvi_acceleration_z gauge")
	for mac, v := range STATE.AccelerationZs() {
		fmt.Fprintf(w,
			"ruuvi_acceleration_z{mac=\"%s\"} %d %d\n",
			mac.String(), v.Value, v.Timestamp.UnixMilli())
	}
	fmt.Fprintln(w, "# HELP ruuvi_acceleration_z", HELP_ACCEL)
	fmt.Fprintln(w, "# TYPE ruuvi_acceleration_z gauge")
	for mac, v := range STATE.AccelerationZs() {
		fmt.Fprintf(w,
			"ruuvi_acceleration_z{mac=\"%s\"} %d %d\n",
			mac.String(), v.Value, v.Timestamp.UnixMilli())
	}
	fmt.Fprintln(w, "# HELP ruuvi_transmit_power", HELP_TX)
	fmt.Fprintln(w, "# TYPE ruuvi_transmit_power gauge")
	for mac, v := range STATE.TransmitPowers() {
		fmt.Fprintf(w,
			"ruuvi_transmit_power{mac=\"%s\"} %d %d\n",
			mac.String(), v.Value, v.Timestamp.UnixMilli())
	}
	fmt.Fprintln(w, "# HELP ruuvi_movement_counter", HELP_MOV)
	fmt.Fprintln(w, "# TYPE ruuvi_movement_counter gauge")
	for mac, v := range STATE.MovementCounters() {
		fmt.Fprintf(w,
			"ruuvi_movement_counter{mac=\"%s\"} %d %d\n",
			mac.String(), v.Value, v.Timestamp.UnixMilli())
	}
	fmt.Fprintln(w, "# HELP ruuvi_sequence_number", HELP_SEQ)
	fmt.Fprintln(w, "# TYPE ruuvi_sequence_number gauge")
	for mac, v := range STATE.SequenceNumbers() {
		fmt.Fprintf(w,
			"ruuvi_sequence_number{mac=\"%s\"} %d %d\n",
			mac.String(), v.Value, v.Timestamp.UnixMilli())
	}
}

func main() {
	var l string
	flag.StringVar(&l, "listen", "localhost:9900", "Listen address")
	flag.BoolVar(&VERBOSE, "verbose", false, "Verbose output")
	flag.Parse()
	log.Println("starting to listen at", l, "and verbosity is", VERBOSE)
	http.HandleFunc("/metrics", metrics)
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
