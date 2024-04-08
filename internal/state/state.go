// package state maintains and exposes all sensor values we have aggregated. We
// prefer simplicity over efficiency.
package state

import (
	"sync"
	"time"

	"github.com/susji/ruuvi/data/rawv2"
	"golang.org/x/exp/maps"
)

type Pair[T comparable] struct {
	Timestamp time.Time
	Value     T
}

func (p Pair[T]) Equal(p2 Pair[T]) bool {
	return p.Value == p2.Value && p.Timestamp.Equal(p2.Timestamp)
}

type Floats map[rawv2.MAC]Pair[float32]
type Uint32s map[rawv2.MAC]Pair[uint32]
type Int16s map[rawv2.MAC]Pair[int16]
type Uint16s map[rawv2.MAC]Pair[uint16]
type Octets map[rawv2.MAC]Pair[uint8]

type State struct {
	temp  Floats
	volt  Floats
	humid Floats
	pres  Uint32s
	accx  Int16s
	accy  Int16s
	accz  Int16s
	txpwr Int16s
	mov   Octets
	seq   Uint16s
	m     *sync.RWMutex
}

func (s *State) Update(d *rawv2.RuuviRawV2) {
	s.m.Lock()
	defer s.m.Unlock()
	if d.Temperature.Valid {
		s.temp[d.MAC] = Pair[float32]{Timestamp: d.Timestamp, Value: d.Temperature.Value}
	}
	if d.BatteryVoltage.Valid {
		s.volt[d.MAC] = Pair[float32]{Timestamp: d.Timestamp, Value: d.BatteryVoltage.Value}
	}
	if d.Humidity.Valid {
		s.humid[d.MAC] = Pair[float32]{Timestamp: d.Timestamp, Value: d.Humidity.Value}
	}
	if d.Pressure.Valid {
		s.pres[d.MAC] = Pair[uint32]{Timestamp: d.Timestamp, Value: d.Pressure.Value}
	}
	if d.AccelerationX.Valid {
		s.accx[d.MAC] = Pair[int16]{Timestamp: d.Timestamp, Value: d.AccelerationX.Value}
	}
	if d.AccelerationY.Valid {
		s.accy[d.MAC] = Pair[int16]{Timestamp: d.Timestamp, Value: d.AccelerationY.Value}
	}
	if d.AccelerationZ.Valid {
		s.accz[d.MAC] = Pair[int16]{Timestamp: d.Timestamp, Value: d.AccelerationZ.Value}
	}
	if d.TransmitPower.Valid {
		s.txpwr[d.MAC] = Pair[int16]{Timestamp: d.Timestamp, Value: d.TransmitPower.Value}
	}
	if d.MovementCounter.Valid {
		s.mov[d.MAC] = Pair[uint8]{Timestamp: d.Timestamp, Value: d.MovementCounter.Value}
	}
	if d.SequenceNumber.Valid {
		s.seq[d.MAC] = Pair[uint16]{Timestamp: d.Timestamp, Value: d.SequenceNumber.Value}
	}
}

func New() *State {
	return &State{
		temp:  Floats{},
		volt:  Floats{},
		humid: Floats{},
		pres:  Uint32s{},
		accx:  Int16s{},
		accy:  Int16s{},
		accz:  Int16s{},
		txpwr: Int16s{},
		mov:   Octets{},
		seq:   Uint16s{},
		m:     &sync.RWMutex{},
	}
}

func (s *State) Temperatures() Floats {
	s.m.RLock()
	defer s.m.RUnlock()
	return maps.Clone(s.temp)
}

func (s *State) Voltages() Floats {
	s.m.RLock()
	defer s.m.RUnlock()
	return maps.Clone(s.volt)
}

func (s *State) Humidities() Floats {
	s.m.RLock()
	defer s.m.RUnlock()
	return maps.Clone(s.humid)
}

func (s *State) Pressures() Uint32s {
	s.m.RLock()
	defer s.m.RUnlock()
	return maps.Clone(s.pres)
}

func (s *State) AccelerationXs() Int16s {
	s.m.RLock()
	defer s.m.RUnlock()
	return maps.Clone(s.accx)
}

func (s *State) AccelerationYs() Int16s {
	s.m.RLock()
	defer s.m.RUnlock()
	return maps.Clone(s.accy)
}

func (s *State) AccelerationZs() Int16s {
	s.m.RLock()
	defer s.m.RUnlock()
	return maps.Clone(s.accz)
}

func (s *State) TransmitPowers() Int16s {
	s.m.RLock()
	defer s.m.RUnlock()
	return maps.Clone(s.txpwr)
}

func (s *State) MovementCounters() Octets {
	s.m.RLock()
	defer s.m.RUnlock()
	return maps.Clone(s.mov)
}

func (s *State) SequenceNumbers() Uint16s {
	s.m.RLock()
	defer s.m.RUnlock()
	return maps.Clone(s.seq)
}
