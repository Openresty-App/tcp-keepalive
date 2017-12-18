package main

import "github.com/grd/statistics"

type Metric struct {
	f []float64
}

func NewMetric() *Metric {
	return &Metric{
		f: make([]float64, 0),
	}
}

func (m *Metric) Put(d float64) {
	m.f = append(m.f, d)
}

func (m *Metric) Max() float64 {
	data := statistics.Float64(m.f)
	d, _ := statistics.Max(&data)
	return d
}

func (m *Metric) Min() float64 {
	data := statistics.Float64(m.f)
	d, _ := statistics.Min(&data)
	return d
}

func (m *Metric) Mean() float64 {
	data := statistics.Float64(m.f)
	return statistics.Mean(&data)
}
