package analytics

import (
	"math"
	"v1/pkg/execute"
)

func mean(data []float64) float64 {
	var sum float64
	for _, value := range data {
		sum += value
	}
	return sum / float64(len(data))
}

func stdDev(data []float64) float64 {
	mean := mean(data)
	var sqDiffSum float64
	for _, value := range data {
		diff := value - mean
		sqDiffSum += diff * diff
	}
	variance := sqDiffSum / float64(len(data)-1)
	return math.Sqrt(variance)
}

func calculateReturns(s *execute.SignalEvents) []float64 {
	var returns []float64
	var buyPrice float64

	for _, signal := range s.Signals {
		if signal.Side == "BUY" {
			buyPrice = signal.Price
		} else if signal.Side == "SELL" && buyPrice != 0 {
			returns = append(returns, (signal.Price-buyPrice)/buyPrice)
			buyPrice = 0 // Reset buy price after a sell
		}
	}

	return returns
}

func SharpeRatio(s *execute.SignalEvents, riskFreeRate float64) float64 {

	if s == nil {
		return 0.0
	}
	// Calculate the returns
	returns := calculateReturns(s)

	// Calculate the excess returns
	excessReturns := make([]float64, len(returns))
	for i, ret := range returns {
		excessReturns[i] = ret - riskFreeRate
	}

	// Calculate the mean and standard deviation of the excess returns
	meanExcessReturn := mean(excessReturns)
	stdDevExcessReturn := stdDev(excessReturns)

	// Calculate the Sharpe Ratio
	sharpeRatio := meanExcessReturn / stdDevExcessReturn

	return sharpeRatio
}
