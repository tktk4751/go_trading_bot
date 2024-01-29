package strategey

// import (
// 	"fmt"
// 	"v1/pkg/analytics"
// 	"v1/pkg/execute"

// 	"github.com/markcheno/go-talib"
// )

// func (df *DataFrameCandle) EmaStrategy(period1, period2 int, account *Account) *execute.SignalEvents {

// 	var StrategyName = "EMA"
// 	lenCandles := len(df.Candles)
// 	if lenCandles <= period1 || lenCandles <= period2 {
// 		return nil
// 	}
// 	signalEvents := execute.NewSignalEvents()
// 	emaValue1 := talib.Ema(df.Closes(), period1)
// 	emaValue2 := talib.Ema(df.Closes(), period2)
// 	rsiValue := talib.Rsi(df.Closes(), 14)

// 	buySize := 0.0
// 	isBuyHolding := false
// 	for i := 1; i < lenCandles; i++ {
// 		if i < period1 || i < period2 {
// 			continue
// 		}
// 		if emaValue1[i-1] < emaValue2[i-1] && emaValue1[i] >= emaValue2[i] && !isBuyHolding {
// 			buySize = account.TradeSize(0.9) / df.Candles[i].Close
// 			signalEvents.Buy(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, buySize, true)
// 			isBuyHolding = true
// 		}
// 		if emaValue1[i-1] > emaValue2[i-1] && emaValue1[i] <= emaValue2[i] && isBuyHolding || rsiValue[i] < 30.0 {
// 			signalEvents.Sell(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, buySize, true)
// 			isBuyHolding = false
// 		}
// 	}
// 	return signalEvents
// }

// func (df *DataFrameCandle) OptimizeEma() (performance float64, bestPeriod1 int, bestPeriod2 int) {
// 	bestPeriod1 = 5
// 	bestPeriod2 = 21

// 	for period1 := 20; period1 < 120; period1++ {
// 		for period2 := 50; period2 < 250; period2++ {
// 			signalEvents := df.EmaStrategy(period1, period2, accountBlance)
// 			if signalEvents == nil {
// 				continue
// 			}
// 			profit := analytics.Profit(signalEvents)
// 			if performance < profit {
// 				performance = profit
// 				bestPeriod1 = period1
// 				bestPeriod2 = period2
// 			}
// 		}
// 	}

// 	fmt.Println("最高利益", performance, "最適なピリオド1", bestPeriod1, "最適なピリオド2", bestPeriod2)

// 	return performance, bestPeriod1, bestPeriod2
// }

// func (df *DataFrameCandle) OptimizeEmaWinRate() (performance float64, bestPeriod1 int, bestPeriod2 int) {
// 	bestPeriod1 = 5
// 	bestPeriod2 = 21

// 	for period1 := 10; period1 < 100; period1++ {
// 		for period2 := 20; period2 < 250; period2++ {
// 			signalEvents := df.EmaStrategy(period1, period2, accountBlance)
// 			if signalEvents == nil {
// 				continue
// 			}
// 			winrate := analytics.WinRate(signalEvents)
// 			if performance < winrate {
// 				performance = winrate
// 				bestPeriod1 = period1
// 				bestPeriod2 = period2
// 			}
// 		}
// 	}

// 	fmt.Println("最高勝率", performance, "最適なピリオド1", bestPeriod1, "最適なピリオド2", bestPeriod2)

// 	return performance, bestPeriod1, bestPeriod2
// }
