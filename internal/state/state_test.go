package state_test

import (
	"maps"
	"testing"
	"time"

	"github.com/susji/ruuvi/data/rawv2"
	"github.com/susji/ruuvimetrics/internal/state"
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
	s.Update(&rawv2.RuuviRawV2{
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
		Humidity: rawv2.Humidity{
			Valid: true,
			Value: 25,
		},
		Pressure: rawv2.Pressure{
			Valid: true,
			Value: 75000,
		},
		AccelerationX: rawv2.Acceleration{
			Valid: true,
			Value: 10,
		},
		AccelerationY: rawv2.Acceleration{
			Valid: true,
			Value: 20,
		},
		AccelerationZ: rawv2.Acceleration{
			Valid: true,
			Value: 30,
		},
		TransmitPower: rawv2.TransmitPower{
			Valid: true,
			Value: 10,
		},
		MovementCounter: rawv2.MovementCounter{
			Valid: true,
			Value: 100,
		},
		SequenceNumber: rawv2.SequenceNumber{
			Valid: true,
			Value: 10000,
		},
	})
	s.Update(&rawv2.RuuviRawV2{
		Timestamp: past,
		MAC:       mac,
		Temperature: rawv2.Temperature{
			Valid: true,
			Value: 642,
		},
		BatteryVoltage: rawv2.BatteryVoltage{
			Valid: true,
			Value: 2.2,
		},
		Humidity: rawv2.Humidity{
			Valid: true,
			Value: 15,
		},
		Pressure: rawv2.Pressure{
			Valid: true,
			Value: 85000,
		},
		AccelerationX: rawv2.Acceleration{
			Valid: true,
			Value: 15,
		},
		AccelerationY: rawv2.Acceleration{
			Valid: true,
			Value: 25,
		},
		AccelerationZ: rawv2.Acceleration{
			Valid: true,
			Value: 35,
		},
		TransmitPower: rawv2.TransmitPower{
			Valid: true,
			Value: 8,
		},
		MovementCounter: rawv2.MovementCounter{
			Valid: true,
			Value: 20,
		},
		SequenceNumber: rawv2.SequenceNumber{
			Valid: true,
			Value: 20000,
		},
	})
	temp := s.Temperatures()
	if got, _ := temp[mac]; got.Value != 321 {
		t.Error(got.Value)
	}
	volt := s.Voltages()
	if got, _ := volt[mac]; got.Value != 2.0 {
		t.Error(got.Value)
	}
	hum := s.Humidities()
	if got, _ := hum[mac]; got.Value != 25 {
		t.Error(got.Value)
	}
	pres := s.Pressures()
	if got, _ := pres[mac]; got.Value != 75000 {
		t.Error(got.Value)
	}
	accx := s.AccelerationXs()
	if got, _ := accx[mac]; got.Value != 10 {
		t.Error(got.Value)
	}
	accy := s.AccelerationYs()
	if got, _ := accy[mac]; got.Value != 20 {
		t.Error(got.Value)
	}
	accz := s.AccelerationZs()
	if got, _ := accz[mac]; got.Value != 30 {
		t.Error(got.Value)
	}
	txpwr := s.TransmitPowers()
	if got, _ := txpwr[mac]; got.Value != 10 {
		t.Error(got.Value)
	}
	mov := s.MovementCounters()
	if got, _ := mov[mac]; got.Value != 100 {
		t.Error(got.Value)
	}
	seq := s.SequenceNumbers()
	if got, _ := seq[mac]; got.Value != 10000 {
		t.Error(got.Value)
	}
}

func TestUpdateInvalidAndMissing(t *testing.T) {
	s := state.New()
	now := time.Now()
	newer := now.Add(1 * time.Hour)
	mac := rawv2.MAC{Value: [6]byte([]byte("123456"))}
	s.Update(&rawv2.RuuviRawV2{
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
	})
	s.Update(&rawv2.RuuviRawV2{
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
	})
	temp := s.Temperatures()
	if got, _ := temp[mac]; got.Value != 642 {
		t.Error(got.Value)
	}
	volt := s.Voltages()
	if got, _ := volt[mac]; got.Value != 2.0 {
		t.Error(got.Value)
	}
	hum := s.Humidities()
	if got, _ := hum[mac]; got.Value != 15 {
		t.Error(got.Value)
	}
}
