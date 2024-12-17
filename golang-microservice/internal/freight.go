package internal

import "math"

type FreightService struct{}

func (fs *FreightService) CalculateFreight(distance int) float64 {
	return math.Floor(float64(distance)*0.15+0.3*(100)) / 100
}

func NewFreightService() *FreightService {
	return &FreightService{}
}
