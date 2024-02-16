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
func SharpeRatio(s *execute.SignalEvents, riskFreeRate float64) float64 {

	if s == nil {
		return 0.0
	}
	// Calculate the returns
	returns := PLSlice(s)

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

func downsideDeviation(data []float64, target float64) float64 {
	var sqDiffSum float64
	for _, value := range data {
		diff := math.Min(value-target, 0) // Only consider negative deviations
		sqDiffSum += diff * diff
	}
	variance := sqDiffSum / float64(len(data))
	return math.Sqrt(variance)
}

// Modified function to calculate the Sortino Ratio
func SortinoRatio(s *execute.SignalEvents, riskFreeRate float64) float64 {

	if s == nil {
		return 0.0
	}
	// Calculate the returns
	returns := PLSlice(s)

	// Calculate the excess returns
	excessReturns := make([]float64, len(returns))
	for i, ret := range returns {
		excessReturns[i] = ret - riskFreeRate
	}

	// Calculate the mean and downside deviation of the excess returns
	meanExcessReturn := mean(excessReturns)
	downsideDeviationExcessReturn := downsideDeviation(excessReturns, 0)

	// Calculate the Sortino Ratio
	sortinoRatio := meanExcessReturn / downsideDeviationExcessReturn

	return sortinoRatio
}
