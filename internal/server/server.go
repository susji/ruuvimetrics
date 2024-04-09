package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/susji/ruuvi/data/rawv2"
	"github.com/susji/ruuvimetrics/internal/config"
	"github.com/susji/ruuvimetrics/internal/state"
)

func dump[T comparable](w http.ResponseWriter, samples map[rawv2.MAC]state.Pair[T], metricfmt, kind, help string) {
	metric := fmt.Sprintf(metricfmt, kind)
	fmt.Fprintf(w, "# HELP %s %s\n", metric, help)
	fmt.Fprintf(w, "# TYPE %s gauge\n", metric)
	for mac, v := range samples {
		fmt.Fprintf(w,
			"%s{mac=\"%s\"} %v %d\n",
			metric, mac.String(), v.Value, v.Timestamp.UnixMilli())
	}
}

type MetricsOptions struct {
	ContentType string
	Endpoint    string
	MetricFmt   string
	Verbose     bool
}

func GenerateMetricsHandler(state *state.State, mo MetricsOptions) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if mo.Verbose {
			log.Println(mo.Endpoint)
		}
		// If performance ever becomes an issue, we could cache the metrics for
		// a configurable amount of time.
		w.Header().Add("content-type", mo.ContentType)
		dump(w, state.Temperatures(), mo.MetricFmt, "temperature", config.TEMP)
		dump(w, state.Voltages(), mo.MetricFmt, "voltage", config.VOLT)
		dump(w, state.Humidities(), mo.MetricFmt, "humidity", config.HUM)
		dump(w, state.Pressures(), mo.MetricFmt, "pressure", config.PRES)
		dump(w, state.AccelerationXs(), mo.MetricFmt, "acceleration_x", config.ACCEL)
		dump(w, state.AccelerationYs(), mo.MetricFmt, "acceleration_y", config.ACCEL)
		dump(w, state.AccelerationZs(), mo.MetricFmt, "acceleration_z", config.ACCEL)
		dump(w, state.TransmitPowers(), mo.MetricFmt, "transmit_power", config.TX)
		dump(w, state.MovementCounters(), mo.MetricFmt, "movement_counter", config.MOV)
		dump(w, state.SequenceNumbers(), mo.MetricFmt, "sequence_number", config.SEQ)
	}
}
