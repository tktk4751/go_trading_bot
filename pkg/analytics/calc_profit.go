package analytics

import "v1/pkg/execute"

func Profit(s *execute.SignalEvents) float64 {
	var profit float64 = 0.0
	var buyPrice, sellPrice float64
	var buySize, sellSize float64

	for _, signal := range s.Signals {
		if signal.Side == "BUY" {
			buyPrice = signal.Price
			buySize = signal.Size
		} else if signal.Side == "SELL" {
			sellPrice = signal.Price
			sellSize = signal.Size
			profit += (sellPrice - buyPrice) * min(buySize, sellSize)
		}
	}

	return profit
}

func NetProfit(s *execute.SignalEvents) float64 {
	totalProfit := TotalProfit(s)
	totalLoss := TotalLoss(s)

	return totalProfit - totalLoss
}

func FinalBalance(s *execute.SignalEvents) (float64, float64) {

	finalBlanceValue := AccountBalance + NetProfit(s)
	finalBlanceRatio := finalBlanceValue / AccountBalance

	return finalBlanceValue, finalBlanceRatio
}
