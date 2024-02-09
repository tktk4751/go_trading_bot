package analytics

import (
	"math"
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

func LongProfit(s *execute.SignalEvents) float64 {
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

func LongLoss(s *execute.SignalEvents) float64 {

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

func ShortProfit(s *execute.SignalEvents) float64 {
	if s == nil {
		return 0.0
	}
	var profit float64 = 0.0
	var sellPrice float64

	if s.Signals == nil || len(s.Signals) == 0 {
		return 0.0
	}
	for _, signal := range s.Signals {

		if signal.Side != "BUY" && signal.Side != "SELL" {
			return 0.0
		}
		if signal.Side == "SELL" {
			sellPrice = signal.Price
		} else if signal.Side == "BUY" && sellPrice != 0 {
			if signal.Price < sellPrice {
				profit += (sellPrice - signal.Price) * signal.Size
			}
			sellPrice = 0 // Reset sell price after a buy
		}
	}

	return profit
}

func ShortLoss(s *execute.SignalEvents) float64 {

	if s == nil {
		return 0.0
	}
	var loss float64 = 0.0
	var sellPrice float64

	if s.Signals == nil || len(s.Signals) == 0 {
		return 0.0
	}
	for _, signal := range s.Signals {

		if signal.Side != "BUY" && signal.Side != "SELL" {
			return 0.0
		}
		if signal.Side == "SELL" {
			sellPrice = signal.Price
		} else if signal.Side == "BUY" && sellPrice != 0 {
			if signal.Price > sellPrice {
				loss += (signal.Price - sellPrice) * signal.Size
			}
			sellPrice = 0 // Reset sell price after a buy
		}
	}

	return loss
}

func TotalProfit(s *execute.SignalEvents) float64 {
	if s == nil {
		return 0.0
	}

	long := LongProfit(s)
	short := ShortProfit(s)

	totalProfit := long + short

	return totalProfit
}

func TotalLoss(s *execute.SignalEvents) float64 {
	if s == nil {
		return 0.0
	}

	long := LongLoss(s)
	short := ShortLoss(s)

	totalLoss := long + short

	return totalLoss
}
func LongNetProfit(s *execute.SignalEvents) float64 {
	if s == nil {
		return 0.0
	}
	totalProfit := LongProfit(s)
	totalLoss := LongLoss(s)

	return totalProfit - totalLoss
}

func ShortNetProfit(s *execute.SignalEvents) float64 {
	if s == nil {
		return 0.0
	}
	totalProfit := ShortProfit(s)
	totalLoss := ShortLoss(s)

	return totalProfit - totalLoss
}

func TotalNetProfit(s *execute.SignalEvents) float64 {
	if s == nil {
		return 0.0
	}
	longProfit := LongNetProfit(s)
	shortProfit := ShortNetProfit(s)

	return longProfit + shortProfit
}

func ProfitFactor(s *execute.SignalEvents) float64 {
	if s == nil {
		return 0.0
	}
	totalProfit := TotalProfit(s)
	totalLoss := TotalLoss(s)

	if totalLoss == 0 {
		return 0.0
	}

	return totalProfit / totalLoss
}

func Prr(s *execute.SignalEvents) float64 {
	if s == nil {
		return 0.0
	}

	pf := ProfitFactor(s)
	n := float64(TotalTrades(s))

	prr := pf / ((n + 1.96*math.Sqrt(n)) / (n - 1.96*math.Sqrt(n)))

	return prr
}

func FinalBalance(s *execute.SignalEvents) (float64, float64) {
	if s == nil {
		return 0.0, 0.0
	}

	accountBalance := 1000.00

	if accountBalance == 0 {
		return 0, 0
	}

	finalBlanceValue := accountBalance + TotalNetProfit(s)
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

func PLSlice(s *execute.SignalEvents) []float64 {
	if s == nil {
		return nil
	}
	var pl []float64 // profit or loss slice
	var buyPrice float64
	var sellPrice float64

	if s.Signals == nil || len(s.Signals) == 0 {
		return nil
	}
	for _, signal := range s.Signals {

		if signal.Side != "BUY" && signal.Side != "SELL" {
			return nil
		}
		if signal.Side == "BUY" {
			buyPrice = signal.Price
			// if there is a previous sell price, calculate the profit or loss
			if sellPrice != 0 {
				pl = append(pl, (sellPrice-buyPrice)*signal.Size)
				// reset the sell price
				sellPrice = 0
			}
		}
		if signal.Side == "SELL" {
			sellPrice = signal.Price
			// if there is a previous buy price, calculate the profit or loss
			if buyPrice != 0 {
				pl = append(pl, (sellPrice-buyPrice)*signal.Size)
				// reset the buy price
				buyPrice = 0
			}
		}
	}

	return pl
}

func TotalProfitSlice(s *execute.SignalEvents) []float64 {
	if s == nil {
		return nil
	}
	var profit []float64 // total profit slice
	var buyPrice float64
	var sellPrice float64
	var longProfit float64
	var shortProfit float64

	if s.Signals == nil || len(s.Signals) == 0 {
		return nil
	}
	for _, signal := range s.Signals {

		if signal.Side != "BUY" && signal.Side != "SELL" {
			return nil
		}
		if signal.Side == "BUY" {
			buyPrice = signal.Price
			// if there is a previous sell price, calculate the short profit
			if sellPrice != 0 {
				shortProfit = (sellPrice - buyPrice) * signal.Size
				profit = append(profit, shortProfit)
				// reset the sell price and short profit
				sellPrice = 0
				shortProfit = 0
			}
		}
		if signal.Side == "SELL" {
			sellPrice = signal.Price
			// if there is a previous buy price, calculate the long profit
			if buyPrice != 0 {
				longProfit = (sellPrice - buyPrice) * signal.Size
				profit = append(profit, longProfit)
				// reset the buy price and long profit
				buyPrice = 0
				longProfit = 0
			}
		}
	}

	return profit
}
