package server_test

import (
	"bufio"
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"net/http/httptest"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/susji/ruuvi/data/rawv2"
	"github.com/susji/ruuvimetrics/internal/config"
	"github.com/susji/ruuvimetrics/internal/server"
	"github.com/susji/ruuvimetrics/internal/state"
	"github.com/susji/ruuvimetrics/internal/testhelpers"
)

type metric struct {
	Name string
	MAC  [6]byte
	TS   time.Time
	Val  string
}

func metricequal(m1, m2 metric) bool {
	return m1.Name == m2.Name &&
		m1.MAC == m2.MAC &&
		m1.TS.Round(time.Second).Equal(m2.TS.Round(time.Millisecond)) &&
		m1.Val == m2.Val
}

// parsemetric does a very simplified parsing of our text-based metric dumps.
func parsemetric(l string) metric {
	// The lines look like this:
	//
	//	ruuvi_temperature{mac="d9:7a:c8:a5:4f:4d"} 5.4849997 171268946557
	//
	// We always have timestamps, we know the metric names and the form of the
	// filter expression, so this is easy to parse.
	name := l[:strings.Index(l, "{")]
	macstr := l[strings.Index(l, `"`)+1 : strings.Index(l, "}")-1]
	val := l[strings.Index(l, "}")+2 : strings.LastIndex(l, " ")]
	milliepoch, _ := strconv.Atoi(l[strings.LastIndex(l, " ")+1:])
	ts := time.UnixMilli(int64(milliepoch))
	var m1, m2, m3, m4, m5, m6 byte
	fmt.Sscanf(macstr, "%02x:%02x:%02x:%02x:%02x:%02x", &m1, &m2, &m3, &m4, &m5, &m6)
	mac := [6]byte{m1, m2, m3, m4, m5, m6}
	return metric{
		Name: name,
		MAC:  mac,
		Val:  val,
		TS:   ts.Round(time.Second),
	}
}

func TestBasicMetrics(t *testing.T) {
	// This was a pretty tedious test to write, but it's fairly flexible and
	// doesn't demand static ordering from the metrics endpoint. It would
	// have been easier to just generate a comparison string of the True
	// Comparison Point but that would be less flexible.
	//
	// First we prime the state with two packets.
	state := state.New()
	now := time.Now()
	past := now.Add(-10 * time.Second)
	macs1 := "123456"
	macs2 := "abcdef"
	mac1 := rawv2.MAC{Value: [6]byte([]byte(macs1))}
	mac2 := rawv2.MAC{Value: [6]byte([]byte(macs2))}
	p1 := testhelpers.Packet1(mac1, now)
	p2 := testhelpers.Packet2(mac2, past)
	state.Update(p1)
	state.Update(p2)
	// Then we fire the request.
	req, err := http.NewRequest("GET", "/metrics", nil)
	if err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()
	handler := server.GenerateMetricsHandler(state, server.MetricsOptions{
		ContentType: config.CT,
		Endpoint:    config.ENDPOINT,
		MetricFmt:   config.METRICFMT,
		Verbose:     config.VERBOSE,
	})
	handler.ServeHTTP(rec, req)
	res := rec.Result()
	t.Run("validate content type header", func(t *testing.T) {
		ct := res.Header["Content-Type"]
		if !slices.Equal(ct, []string{config.CT}) {
			t.Error(ct)
		}
	})
	metriclineno := map[string][]int{}
	t.Run("validate the contents of the metrics printout", func(t *testing.T) {
		// Now we verify the contents of the metrics dump. For now we will all
		// comment lines and make sure the state values are actually there and
		// that each metric type is consecutively printed.
		lines := []string{}
		sc := bufio.NewScanner(rec.Body)
		// We'll keep track of the values we find in these two maps.
		n := func(mac, metric string) string { return mac + "_" + metric }
		foundmacs := map[[6]byte]bool{}
		foundmetrics := map[string]metric{}
		lineno := 0
		for sc.Scan() {
			l := sc.Text()
			if l[0] == '#' {
				continue
			}
			lineno++
			lines = append(lines, l)
			m := parsemetric(l)
			foundmacs[m.MAC] = true
			foundmetrics[n(string(m.MAC[:]), m.Name)] = m
			metriclineno[m.Name] = append(metriclineno[m.Name], lineno)
		}
		// First validate we have the two MAC addresses.
		wantmacs := map[[6]byte]bool{
			mac1.Value: true,
			mac2.Value: true,
		}
		if !maps.Equal(foundmacs, wantmacs) {
			t.Error(wantmacs)
		}
		// Then the metric values themselves...
		gm := func(p *rawv2.RuuviRawV2, name string, value any) metric {
			return metric{
				Name: fmt.Sprintf(config.METRICFMT, name),
				MAC:  p.MAC.Value,
				TS:   p.Timestamp.Round(time.Second),
				Val:  fmt.Sprintf("%v", value),
			}
		}
		n2 := func(mac, metric string) string { return mac + "_" + fmt.Sprintf(config.METRICFMT, metric) }
		wantmetrics := map[string]metric{
			n2(macs1, "temperature"):      gm(p1, "temperature", p1.Temperature.Value),
			n2(macs2, "temperature"):      gm(p2, "temperature", p2.Temperature.Value),
			n2(macs1, "voltage"):          gm(p1, "voltage", p1.BatteryVoltage.Value),
			n2(macs2, "voltage"):          gm(p2, "voltage", p2.BatteryVoltage.Value),
			n2(macs1, "humidity"):         gm(p1, "humidity", p1.Humidity.Value),
			n2(macs2, "humidity"):         gm(p2, "humidity", p2.Humidity.Value),
			n2(macs1, "pressure"):         gm(p1, "pressure", p1.Pressure.Value),
			n2(macs2, "pressure"):         gm(p2, "pressure", p2.Pressure.Value),
			n2(macs1, "acceleration_x"):   gm(p1, "acceleration_x", p1.AccelerationX.Value),
			n2(macs2, "acceleration_x"):   gm(p2, "acceleration_x", p2.AccelerationX.Value),
			n2(macs1, "acceleration_y"):   gm(p1, "acceleration_y", p1.AccelerationY.Value),
			n2(macs2, "acceleration_y"):   gm(p2, "acceleration_y", p2.AccelerationY.Value),
			n2(macs1, "acceleration_z"):   gm(p1, "acceleration_z", p1.AccelerationZ.Value),
			n2(macs2, "acceleration_z"):   gm(p2, "acceleration_z", p2.AccelerationZ.Value),
			n2(macs1, "transmit_power"):   gm(p1, "transmit_power", p1.TransmitPower.Value),
			n2(macs2, "transmit_power"):   gm(p2, "transmit_power", p2.TransmitPower.Value),
			n2(macs1, "movement_counter"): gm(p1, "movement_counter", p1.MovementCounter.Value),
			n2(macs2, "movement_counter"): gm(p2, "movement_counter", p2.MovementCounter.Value),
			n2(macs1, "sequence_number"):  gm(p1, "sequence_number", p1.SequenceNumber.Value),
			n2(macs2, "sequence_number"):  gm(p2, "sequence_number", p2.SequenceNumber.Value),
		}
		if !maps.EqualFunc(wantmetrics, foundmetrics, metricequal) {
			t.Log(mustJson(wantmetrics))
			t.Log(mustJson(foundmetrics))
			t.Error()
		}
	})
	t.Run("make sure line numbers are consecutive for the same metric", func(t *testing.T) {
		for metric, lineno := range metriclineno {
			d := lineno[0] - lineno[1]
			if d != -1 && d != 1 {
				t.Error(metric, lineno)
			}
		}
	})
}

func mustJson(v any) string {
	b, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(b)
}
