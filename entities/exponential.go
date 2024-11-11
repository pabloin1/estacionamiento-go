package entities

import "gonum.org/v1/gonum/stat/distuv"

type ExponentialDist struct {
}

func NewExponentialDist() *ExponentialDist {
	return &ExponentialDist{}
}

func (ed *ExponentialDist) Generate(rate float64) float64 {
	exp := distuv.Exponential{Rate: rate, Src: nil}
	return exp.Rand()
}
