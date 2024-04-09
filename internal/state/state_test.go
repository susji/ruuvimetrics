package state_test

import (
	"maps"
	"testing"
	"time"

	"github.com/susji/ruuvi/data/rawv2"
	"github.com/susji/ruuvimetrics/internal/state"
	"github.com/susji/ruuvimetrics/internal/testhelpers"
)

func TestEmptyStateGetting(t *testing.T) {
	s := state.New()
	if !maps.Equal(s.Temperatures(), state.Floats{}) {
		t.Error()
	}
	if !maps.Equal(s.Voltages(), state.Floats{}) {
		t.Error()
	}
	if !maps.Equal(s.Humidities(), state.Floats{}) {
		t.Error()
	}
	if !maps.Equal(s.Pressures(), state.Uint32s{}) {
		t.Error()
	}
	if !maps.Equal(s.AccelerationXs(), state.Int16s{}) {
		t.Error()
	}
	if !maps.Equal(s.AccelerationYs(), state.Int16s{}) {
		t.Error()
	}
	if !maps.Equal(s.AccelerationZs(), state.Int16s{}) {
		t.Error()
	}
	if !maps.Equal(s.TransmitPowers(), state.Int16s{}) {
		t.Error()
	}
	if !maps.Equal(s.MovementCounters(), state.Octets{}) {
		t.Error()
	}
	if !maps.Equal(s.SequenceNumbers(), state.Uint16s{}) {
		t.Error()
	}
}

func TestUpdateWithOlder(t *testing.T) {
	s := state.New()
	now := time.Now()
	past := now.Add(-1 * time.Hour)
	mac := rawv2.MAC{Value: [6]byte([]byte("123456"))}
	dnow := testhelpers.Packet1(mac, now)
	dpast := testhelpers.Packet2(mac, past)
	s.Update(dnow)
	s.Update(dpast)
	temp := s.Temperatures()
	if got, _ := temp[mac]; got.Value != dnow.Temperature.Value {
		t.Error(got.Value)
	}
	volt := s.Voltages()
	if got, _ := volt[mac]; got.Value != dnow.BatteryVoltage.Value {
		t.Error(got.Value)
	}
	hum := s.Humidities()
	if got, _ := hum[mac]; got.Value != dnow.Humidity.Value {
		t.Error(got.Value)
	}
	pres := s.Pressures()
	if got, _ := pres[mac]; got.Value != dnow.Pressure.Value {
		t.Error(got.Value)
	}
	accx := s.AccelerationXs()
	if got, _ := accx[mac]; got.Value != dnow.AccelerationX.Value {
		t.Error(got.Value)
	}
	accy := s.AccelerationYs()
	if got, _ := accy[mac]; got.Value != dnow.AccelerationY.Value {
		t.Error(got.Value)
	}
	accz := s.AccelerationZs()
	if got, _ := accz[mac]; got.Value != dnow.AccelerationZ.Value {
		t.Error(got.Value)
	}
	txpwr := s.TransmitPowers()
	if got, _ := txpwr[mac]; got.Value != dnow.TransmitPower.Value {
		t.Error(got.Value)
	}
	mov := s.MovementCounters()
	if got, _ := mov[mac]; got.Value != dnow.MovementCounter.Value {
		t.Error(got.Value)
	}
	seq := s.SequenceNumbers()
	if got, _ := seq[mac]; got.Value != dnow.SequenceNumber.Value {
		t.Error(got.Value)
	}
}

func TestUpdateInvalidAndMissing(t *testing.T) {
	s := state.New()
	now := time.Now()
	newer := now.Add(1 * time.Hour)
	mac := rawv2.MAC{Value: [6]byte([]byte("123456"))}
	dnow := &rawv2.RuuviRawV2{
		Timestamp: now,
		MAC:       mac,
		Temperature: rawv2.Temperature{
			Valid: true,
			Value: 321,
		},
		BatteryVoltage: rawv2.BatteryVoltage{
			Valid: true,
			Value: 2.0,
		},
	}
	dnewer := &rawv2.RuuviRawV2{
		Timestamp: newer,
		MAC:       mac,
		Temperature: rawv2.Temperature{
			Valid: true,
			Value: 642,
		},
		BatteryVoltage: rawv2.BatteryVoltage{
			Valid: false,
			Value: 1.8,
		},
		Humidity: rawv2.Humidity{
			Valid: true,
			Value: 15,
		},
	}
	s.Update(dnow)
	s.Update(dnewer)
	temp := s.Temperatures()
	if got, _ := temp[mac]; got.Value != dnewer.Temperature.Value {
		t.Error(got.Value)
	}
	volt := s.Voltages()
	if got, _ := volt[mac]; got.Value != dnow.BatteryVoltage.Value {
		t.Error(got.Value)
	}
	hum := s.Humidities()
	if got, _ := hum[mac]; got.Value != dnewer.Humidity.Value {
		t.Error(got.Value)
	}
}
