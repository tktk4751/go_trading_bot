package analytics

import "v1/pkg/execute"

func MaxDrawdown(s *execute.SignalEvents) float64 {
	var maxPeakPrice float64 = 0.0
	var maxDrawdown float64 = 0.0

	for _, signal := range s.Signals {
		if signal.Side == "BUY" {
			if signal.Price > maxPeakPrice {
				maxPeakPrice = signal.Price
			}
		} else if signal.Side == "SELL" {
			drawdown := (maxPeakPrice - signal.Price) / maxPeakPrice
			if drawdown > maxDrawdown {
				maxDrawdown = drawdown
			}
			if signal.Price > maxPeakPrice {
				maxPeakPrice = signal.Price
			}
		}
	}

	return maxDrawdown
}
