package analytics

import (
	"math"
	"v1/pkg/execute"
)

func calculateStandardDeviation(s *execute.SignalEvents) float64 {
	if s == nil || s.Signals == nil || len(s.Signals) == 0 {
		return 0.0
	}

	var profits []float64
	var buyPrice float64

	for _, signal := range s.Signals {
		if signal.Side != "BUY" && signal.Side != "SELL" {
			continue
		}
		if signal.Side == "BUY" {
			buyPrice = signal.Price
		} else if signal.Side == "SELL" && buyPrice != 0 {
			profit := (signal.Price - buyPrice) * signal.Size
			profits = append(profits, profit)
			buyPrice = 0 // Reset buy price after a sell
		}
	}

	// Calculate mean of profits
	mean := 0.0
	for _, profit := range profits {
		mean += profit
	}
	mean /= float64(len(profits))

	// Calculate standard deviation of profits
	variance := 0.0
	for _, profit := range profits {
		difference := profit - mean
		squaredDifference := difference * difference
		variance += squaredDifference
	}
	variance /= float64(len(profits) - 1)

	standardDeviation := math.Sqrt(variance)

	return standardDeviation
}

func SQN(s *execute.SignalEvents) float64 {
	if s == nil {
		return 0.0
	}

	e := ExpectedValue(s)
	tt := TotalTrades(s)

	sdv := calculateStandardDeviation(s)

	sqn := math.Sqrt(float64(tt)) * e / sdv

	return sqn
}
