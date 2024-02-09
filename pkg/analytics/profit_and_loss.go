package analytics

import (
	"time"

	dbquery "v1/pkg/data/query"
	"v1/pkg/execute"
	"v1/pkg/trader"
)

// 課題 エグジットフラグメントを実装して､空売りにも対応するProfit関数を作ろう

// func Profit(s *execute.SignalEvents) float64 {
// 	if s == nil {
// 		return 0.0
// 	}
// 	total := 0.0
// 	beforeSell := 0.0
// 	isHolding := false
// 	for i, signalEvent := range s.Signals {
// 		if i == 0 && signalEvent.Side == "SELL" {
// 			continue
// 		}
// 		if signalEvent.Side == "BUY" {
// 			total -= signalEvent.Price * signalEvent.Size
// 			isHolding = true
// 		}
// 		if signalEvent.Side == "SELL" {
// 			total += signalEvent.Price * signalEvent.Size
// 			isHolding = false
// 			beforeSell = total
// 		}
// 	}
// 	if isHolding {
// 		return beforeSell
// 	}
// 	return total
// }

// これは純利益を出力する関数
func Profi2(s *execute.SignalEvents) float64 {
	if s == nil {
		return 0.0
	}
	profit := 0.0
	beforeSell := 0.0
	isHolding := false

	if s.Signals == nil || len(s.Signals) == 0 {
		return 0.0
	}
	for i, signalEvent := range s.Signals {
		if i == 0 && signalEvent.Side == "SELL" {
			continue
		}
		if signalEvent.Side == "BUY" {
			profit -= signalEvent.Price * signalEvent.Size
			isHolding = true
		}
		if signalEvent.Side == "SELL" {
			profit += signalEvent.Price * signalEvent.Size
			isHolding = false
			beforeSell = profit
		}
	}
	if isHolding {
		return beforeSell
	}
	return profit
}

func Profit(s *execute.SignalEvents) float64 {
	if s == nil {
		return 0.0
	}
	var profit float64 = 0.0
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
			if signal.Price > buyPrice {
				profit += (signal.Price - buyPrice) * signal.Size
			}
			buyPrice = 0 // Reset buy price after a sell
		}
	}

	return profit
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

	// if totalLoss == 0 {
	// 	return math.Inf(1)
	// }

	return totalProfit / totalLoss
}

func FinalBalance(s *execute.SignalEvents) (float64, float64) {
	if s == nil {
		return 0.0, 0.0
	}

	accountBalance := 1000.00

	if accountBalance == 0 {
		return 0, 0
	}

	finalBlanceValue := accountBalance + NetProfit(s)
	finalBlanceRatio := finalBlanceValue / accountBalance

	return finalBlanceValue, finalBlanceRatio
}

func MaxLossTrade(s *execute.SignalEvents) (float64, time.Time) {
	if s == nil || s.Signals == nil || len(s.Signals) == 0 {
		return 0.0, time.Time{}
	}
	var maxLossTrade float64 = 0.0
	var lossTime time.Time
	var buyPrice float64
	for _, signal := range s.Signals {
		if signal.Side != "BUY" && signal.Side != "SELL" {
			continue
		}
		if signal.Side == "BUY" {
			buyPrice = signal.Price
		} else if signal.Side == "SELL" && buyPrice != 0 {
			if signal.Price < buyPrice {
				loss := (buyPrice - signal.Price) * signal.Size
				if loss > maxLossTrade {
					maxLossTrade = loss
					lossTime = signal.Time
				}
			}
			buyPrice = 0 // Reset buy price after a sell
		}
	}
	return maxLossTrade, lossTime
}

func MaxProfitTrade(s *execute.SignalEvents) (float64, time.Time) {
	if s == nil || s.Signals == nil || len(s.Signals) == 0 {
		return 0.0, time.Time{}
	}
	var maxProfitTrade float64 = 0.0
	var profitTime time.Time
	var total float64 = 0.0
	var buyPrice float64
	for _, signal := range s.Signals {
		if signal.Side != "BUY" && signal.Side != "SELL" {
			continue
		}
		if signal.Side == "BUY" {
			total -= signal.Price * signal.Size
			buyPrice = signal.Price
		} else if signal.Side == "SELL" && buyPrice != 0 {
			total += signal.Price * signal.Size
			profit := (signal.Price - buyPrice) * signal.Size
			if profit > maxProfitTrade {
				maxProfitTrade = profit
				profitTime = signal.Time
			}
		}
	}
	return maxProfitTrade, profitTime
}

func HoldingReturn(s *execute.SignalEvents) (float64, float64) {

	acount := trader.NewAccount(initialBalance)

	assetName := s.Signals[0].AssetName
	duration := s.Signals[0].Duration
	close, _ := dbquery.GetCloseData(assetName, duration)
	size := acount.Balance

	buyPrice := close[len(close)-1]

	holdingReturn := (close[0] / buyPrice) * size
	holdingReturnRatio := holdingReturn / size

	return holdingReturn, holdingReturnRatio
}

func GainPainRatio(s *execute.SignalEvents) float64 {
	if s == nil {
		return 0.0
	}

	gain := AveregeProfit(s)
	pain := MaxDrawdownUSD(s)

	gainPainRatio := gain / pain

	return gainPainRatio

}

func ReturnProfitLoss(s *execute.SignalEvents) []float64 {
	if s == nil {
		return nil
	}
	var pl []float64 // profit or loss slice
	var buyPrice float64

	if s.Signals == nil || len(s.Signals) == 0 {
		return nil
	}
	for _, signal := range s.Signals {

		if signal.Side != "BUY" && signal.Side != "SELL" {
			return nil
		}
		if signal.Side == "BUY" {
			buyPrice = signal.Price
		} else if signal.Side == "SELL" && buyPrice != 0 {
			pl = append(pl, (signal.Price-buyPrice)*signal.Size) // append the profit or loss of the trade to the slice
			buyPrice = 0                                         // Reset buy price after a sell
		}
	}

	return pl
}
