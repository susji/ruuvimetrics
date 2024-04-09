package testhelpers

import (
	"time"

	"github.com/susji/ruuvi/data/rawv2"
)

func Packet1(mac rawv2.MAC, ts time.Time) *rawv2.RuuviRawV2 {
	return &rawv2.RuuviRawV2{
		Timestamp: ts,
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
	}
}

func Packet2(mac rawv2.MAC, ts time.Time) *rawv2.RuuviRawV2 {
	return &rawv2.RuuviRawV2{
		Timestamp: ts,
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
	}
}
