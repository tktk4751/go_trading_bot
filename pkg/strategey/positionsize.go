package strategey

import (
	"v1/pkg/execute"
)

var AccountBalance float64 = 1000.000

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

func NetProfit(s *execute.SignalEvents) float64 {
	totalProfit := TotalProfit(s)
	totalLoss := TotalLoss(s)

	return totalProfit - totalLoss
}

// func RiskSizeCalculator(s *execute.SignalEvents) float64 {

// 	w := s.WinRate()
// 	r := s.ProfitFactor()
// 	d := s.MaxDrawdown()

// 	// f := (w*(r+1)-1)/r - (d*(d*2+1)-1)/r - 0.002
// 	f := (((w*(r+w+w)-(1+d))/(r-w*d) - 0.002) * w) / 1.618

// 	if f < 0 || r <= 1.05 || d > 0.45 {
// 		fmt.Print("トレード禁止")
// 		return 0
// 	}
// 	return f
// }
