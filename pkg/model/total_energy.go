package model

import "go.uber.org/atomic"

type TotalEnergy struct {
	Idx    int64
	UserId string
	//TotalEnergy       int
	TotalEnergyAtomic *atomic.Int32
}
