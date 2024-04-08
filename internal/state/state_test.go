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
		MovementCounter: rawv2.MovementCounter{
			Valid: true,
			Value: 100,
		},
	})
	s.Update(&rawv2.RuuviRawV2{
		Timestamp: past,
		MAC:       mac,
		MovementCounter: rawv2.MovementCounter{
			Valid: true,
			Value: 20,
		},
	})
	mov := s.MovementCounters()
	if got, _ := mov[mac]; got.Value != 100 {
		t.Error()
	}
}
