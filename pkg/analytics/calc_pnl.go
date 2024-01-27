package analytics

import (
	"math"
	"v1/pkg/execute"
)

func TotalProfit(s *execute.SignalEvents) float64 {
	var totalProfit float64 = 0.0
	var buyPrice, sellPrice float64
	var buySize, sellSize float64

	for _, signal := range s.Signals {
		if signal.Side == "BUY" {
			buyPrice = signal.Price
			buySize = signal.Size
		} else if signal.Side == "SELL" {
			sellPrice = signal.Price
			sellSize = signal.Size
			profit := (sellPrice - buyPrice) * min(buySize, sellSize) / buyPrice * AccountBalance
			if profit > 0 {
				totalProfit += profit
			}
		}
	}

	return totalProfit
}

func TotalLoss(s *execute.SignalEvents) float64 {
	var totalLoss float64 = 0.0
	var buyPrice, sellPrice float64
	var buySize, sellSize float64

	for _, signal := range s.Signals {
		if signal.Side == "BUY" {
			buyPrice = signal.Price
			buySize = signal.Size
		} else if signal.Side == "SELL" {
			sellPrice = signal.Price
			sellSize = signal.Size
			profit := (sellPrice - buyPrice) * min(buySize, sellSize) / buyPrice * AccountBalance
			if profit < 0 {
				totalLoss -= profit
			}
		}
	}

	return totalLoss
}

func ProfitFactor(s *execute.SignalEvents) float64 {
	totalProfit := TotalProfit(s)
	totalLoss := TotalLoss(s)

	if totalLoss == 0 {
		return math.Inf(1)
	}

	return totalProfit / totalLoss
}
