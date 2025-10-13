package utils

import "math"

func FloatTo2(num float64) float64 {
	return math.Round(num*100) / 100
}
