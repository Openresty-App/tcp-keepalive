package utils

import (
	"github.com/grd/statistics"
)

type Metric struct {
	f []float64
	l uint64
}

func NewMetric() *Metric {
	return &Metric{
		f: []float64{},
		l: 0,
	}
}

func (m *Metric) Put(d float64) {
	m.f = append(m.f, d)
	m.l += 1
}

func (m *Metric) Len() uint64 {
	return m.l
}

func (m *Metric) Max() float64 {
	data := statistics.Float64(m.f)
	if data == nil || len(data) == 0 {
		return 0.0
	}
	d, _ := statistics.Max(&data)
	return d
}

func (m *Metric) Min() float64 {
	data := statistics.Float64(m.f)
	if data == nil || len(data) == 0 {
		return 0.0
	}
	d, _ := statistics.Min(&data)
	return d
}

func (m *Metric) Mean() float64 {
	data := statistics.Float64(m.f)
	if data == nil || len(data) == 0 {
		return 0.0
	}
	return statistics.Mean(&data)
}
