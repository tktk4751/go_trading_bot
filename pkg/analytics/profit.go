package analytics

import (
	"math"
	"v1/pkg/execute"
)

// 課題 エグジットフラグメントを実装して､空売りにも対応するProfit関数を作ろう

func Profit(s *execute.SignalEvents) float64 {
	if s == nil {
		return 0.0
	}
	total := 0.0
	beforeSell := 0.0
	isHolding := false
	for i, signalEvent := range s.Signals {
		if i == 0 && signalEvent.Side == "SELL" {
			continue
		}
		if signalEvent.Side == "BUY" {
			total -= signalEvent.Price * signalEvent.Size
			isHolding = true
		}
		if signalEvent.Side == "SELL" {
			total += signalEvent.Price * signalEvent.Size
			isHolding = false
			beforeSell = total
		}
	}
	if isHolding {
		return beforeSell
	}
	return total
}

func Loss(s *execute.SignalEvents) float64 {

	if s == nil {
		return 0.0
	}
	var loss float64 = 0.0
	var buyPrice float64

	if s.Signals == nil || len(s.Signals) == 0 {
		return 0.0
	}
	for _, signal := range s.Signals {

		if signal.Side != "BUY" && signal.Side != "SELL" {
			return 0.0
		}
		if signal.Side == "BUY" {
			buyPrice = signal.Price
		} else if signal.Side == "SELL" && buyPrice != 0 {
			if signal.Price < buyPrice {
				loss += (buyPrice - signal.Price) * signal.Size
			}
			buyPrice = 0 // Reset buy price after a sell
		}
	}

	return loss
}

// // TotalProfit returns the total profit of a series of signal events
// func TotalProfit(s *execute.SignalEvents) float64 {
// 	var totalProfit float64 = 0.0
// 	for _, signal := range s.Signals {
// 		if signal.Side == "SELL" {
// 			totalProfit += Profit(s)
// 		}
// 	}
// 	return totalProfit
// }

// // TotalLoss returns the total loss of a series of signal events
// func TotalLoss(s *execute.SignalEvents) float64 {
// 	var totalLoss float64 = 0.0
// 	for _, signal := range s.Signals {
// 		if signal.Side == "BUY" {
// 			totalLoss -= Profit(s)
// 		}
// 	}
// 	return totalLoss
// }

func NetProfit(s *execute.SignalEvents) float64 {
	if s == nil {
		return 0.0
	}
	totalProfit := Profit(s)
	totalLoss := Loss(s)

	return totalProfit - totalLoss
}

func ProfitFactor(s *execute.SignalEvents) float64 {
	if s == nil {
		return 0.0
	}
	totalProfit := Profit(s)
	totalLoss := Loss(s)

	if totalLoss == 0 {
		return math.Inf(1)
	}

	return totalProfit / totalLoss
}

func FinalBalance(s *execute.SignalEvents) (float64, float64) {
	if s == nil {
		return 0.0, 0.0
	}

	if AccountBalance == 0 {
		return 0, 0
	}

	finalBlanceValue := AccountBalance + NetProfit(s)
	finalBlanceRatio := finalBlanceValue / AccountBalance

	return finalBlanceValue, finalBlanceRatio
}
