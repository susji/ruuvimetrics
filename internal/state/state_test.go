package state_test

import (
	"maps"
	"testing"

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
