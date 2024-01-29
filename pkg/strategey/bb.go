package strategey

// import (
// 	"fmt"
// 	"math"
// 	"v1/pkg/analytics"
// 	"v1/pkg/execute"

// 	"github.com/markcheno/go-talib"
// )

// func (df *DataFrameCandle) BBStrategy(n int, k float64, account *Account) *execute.SignalEvents {

// 	var StrategyName = "BB"
// 	lenCandles := len(df.Candles)

// 	if lenCandles <= n {
// 		return nil
// 	}

// 	signalEvents := execute.NewSignalEvents()
// 	bbUp, _, bbDown := talib.BBands(df.Closes(), n, k, k, 0)

// 	buySize := 0.0
// 	isBuyHolding := false
// 	for i := 1; i < lenCandles; i++ {
// 		if i < n {
// 			continue
// 		}
// 		if bbDown[i-1] > df.Candles[i-1].Close && bbDown[i] <= df.Candles[i].Close && !isBuyHolding {
// 			buySize = account.TradeSize(0.9) / df.Candles[i].Close
// 			signalEvents.Buy(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, buySize, true)
// 			isBuyHolding = true
// 		}
// 		if bbUp[i-1] < df.Candles[i-1].Close && bbUp[i] >= df.Candles[i].Close && isBuyHolding {
// 			signalEvents.Sell(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, buySize, true)
// 			isBuyHolding = false
// 		}
// 	}
// 	return signalEvents
// }

// func (df *DataFrameCandle) OptimizeBbProfit() (performance float64, bestN int, bestK float64) {
// 	bestN = 20
// 	bestK = 2.0

// 	for n := 13; n < 200; n++ {
// 		for k := 2.0; k < 3.5; k += 0.1 {
// 			signalEvents := df.BBStrategy(n, k, accountBlance)
// 			if signalEvents == nil {
// 				continue
// 			}
// 			profit := analytics.Profit(signalEvents)
// 			if performance < profit {
// 				performance = profit
// 				bestN = n
// 				bestK = k
// 			}
// 		}
// 	}

// 	fmt.Println("最高利益", performance, "最適なピリオド", bestN, "最適な標準偏差", bestK)

// 	return performance, bestN, bestK
// }

// func (df *DataFrameCandle) OptimizeBbLoss() (performance float64, bestN int, bestK float64) {
// 	bestN = 20
// 	bestK = 2.0
// 	performance = math.MaxFloat64

// 	for n := 5; n < 120; n++ {
// 		for k := 1.8; k < 3.8; k += 0.1 {
// 			signalEvents := df.BBStrategy(n, k, accountBlance)
// 			if signalEvents == nil {
// 				continue
// 			}
// 			loss := analytics.Loss(signalEvents)
// 			if performance < loss {
// 				performance = loss
// 				bestN = n
// 				bestK = k
// 			}
// 		}
// 	}

// 	fmt.Println("損失", performance, "最適なピリオド", bestN, "最適な標準偏差", bestK)

// 	return performance, bestN, bestK
// }

// func (df *DataFrameCandle) OptimizeBbWinRate() (performance float64, bestN int, bestK float64) {
// 	bestN = 20
// 	bestK = 2.0

// 	for n := 13; n < 200; n++ {
// 		for k := 2.0; k < 3.5; k += 0.1 {
// 			signalEvents := df.BBStrategy(n, k, accountBlance)
// 			if signalEvents == nil {
// 				continue
// 			}

// 			if analytics.TotalTrades(signalEvents) < 5 {
// 				continue
// 			}
// 			winrate := analytics.WinRate(signalEvents)
// 			if performance < winrate {
// 				performance = winrate
// 				bestN = n
// 				bestK = k
// 			}
// 		}
// 	}

// 	fmt.Println("最高勝率", performance, "最適なピリオド", bestN, "最適な標準偏差", bestK)

// 	return performance, bestN, bestK
// }
