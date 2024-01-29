package strategey

// import (
// 	"fmt"
// 	"math"
// 	"v1/pkg/analytics"
// 	"v1/pkg/execute"

// 	"github.com/markcheno/go-talib"
// )

// func (df *DataFrameCandle) RsiStrategy(period int, buyThread float64, sellThread float64, account *Account) *execute.SignalEvents {

// 	var StrategyName = "RSI"
// 	lenCandles := len(df.Candles)
// 	if lenCandles <= period {
// 		return nil
// 	}

// 	signalEvents := execute.NewSignalEvents()

// 	values := talib.Rsi(df.Closes(), period)

// 	buySize := 0.0
// 	isBuyHolding := false
// 	for i := 1; i < lenCandles; i++ {
// 		if values[i-1] == 0 || values[i-1] == 100 {
// 			continue
// 		}
// 		if values[i-1] < buyThread && values[i] >= buyThread && !isBuyHolding {
// 			buySize = account.TradeSize(0.9) / df.Candles[i].Close
// 			signalEvents.Buy(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, buySize, true)
// 			isBuyHolding = true
// 		}

// 		if values[i-1] > float64(sellThread) && values[i] <= float64(sellThread) && isBuyHolding {
// 			signalEvents.Sell(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, buySize, true)
// 			isBuyHolding = false
// 		}
// 	}

// 	return signalEvents
// }

// func (df *DataFrameCandle) OptimizeRsi() (performance float64, bestPeriod int, bestBuyThread, bestSellThread float64) {
// 	bestPeriod = 14
// 	bestBuyThread, bestSellThread = 15.0, 80.0

// 	for period := 5; period < 60; period++ {

// 		signalEvents := df.RsiStrategy(period, bestBuyThread, bestSellThread, accountBlance)
// 		if signalEvents == nil {
// 			continue
// 		}

// 		profit := Profit(signalEvents)
// 		if performance < profit {
// 			performance = profit
// 			bestPeriod = period
// 			bestBuyThread = bestBuyThread
// 			bestSellThread = bestSellThread
// 		}

// 	}

// 	fmt.Println("最高利益", performance, "最適なピリオド", bestPeriod, "最適な上限ライン", bestBuyThread, "最適な下限ライン", bestSellThread)

// 	return performance, bestPeriod, bestBuyThread, bestSellThread
// }

// func (df *DataFrameCandle) OptimizeRsi2() (performance float64, bestPeriod int, bestBuyThread, bestSellThread float64) {
// 	bestPeriod = 14
// 	bestBuyThread, bestSellThread = 20.0, 80.0

// 	for period := 5; period < 30; period++ {
// 		for buyThread := 25.0; buyThread > 10; buyThread -= 1 {

// 			for sellThread := 75.0; sellThread < 95; sellThread += 1 {
// 				signalEvents := df.RsiStrategy(period, buyThread, sellThread, accountBlance)
// 				if signalEvents == nil {
// 					continue
// 				}

// 				profit := Profit(signalEvents)
// 				if performance < profit {
// 					performance = profit
// 					bestPeriod = period
// 					bestBuyThread = buyThread
// 					bestSellThread = sellThread
// 				}
// 			}
// 		}

// 	}

// 	fmt.Println("最高利益", performance, "最適なピリオド", bestPeriod, "最適な買いライン", bestBuyThread, "最適な売りライン", bestSellThread)

// 	return performance, bestPeriod, bestBuyThread, bestSellThread
// }

// func (df *DataFrameCandle) OptimizeRsiWinRate() (performance float64, bestPeriod int, bestBuyThread, bestSellThread float64) {
// 	bestPeriod = 14
// 	bestBuyThread, bestSellThread = 20.0, 80.0

// 	for period := 5; period < 30; period++ {
// 		for buyThread := 25.0; buyThread > 10; buyThread -= 1 {

// 			for sellThread := 75.0; sellThread < 95; sellThread += 1 {
// 				signalEvents := df.RsiStrategy(period, buyThread, sellThread, accountBlance)
// 				if signalEvents == nil {
// 					continue
// 				}

// 				if analytics.TotalTrades(signalEvents) < 20 {
// 					continue
// 				}

// 				winrate := analytics.WinRate(signalEvents)
// 				if performance < winrate {
// 					performance = winrate
// 					bestPeriod = period
// 					bestBuyThread = buyThread
// 					bestSellThread = sellThread
// 				}
// 			}
// 		}

// 	}

// 	fmt.Println("最高勝率", performance*100, "%", "最適なピリオド", bestPeriod, "最適な買いライン", bestBuyThread, "最適な売りライン", bestSellThread)

// 	return performance, bestPeriod, bestBuyThread, bestSellThread
// }

// func (df *DataFrameCandle) OptimizeRsiLoss() (performance float64, bestPeriod int, bestBuyThread, bestSellThread float64) {
// 	bestPeriod = 14
// 	bestBuyThread, bestSellThread = 20.0, 80.0
// 	performance = math.MaxFloat64

// 	for period := 5; period < 30; period++ {
// 		for buyThread := 25.0; buyThread > 10; buyThread -= 1 {

// 			for sellThread := 75.0; sellThread < 96; sellThread += 1 {
// 				signalEvents := df.RsiStrategy(period, buyThread, sellThread, accountBlance)
// 				if signalEvents == nil {
// 					continue
// 				}

// 				if analytics.TotalTrades(signalEvents) < 20 {
// 					continue
// 				}

// 				loss := analytics.Loss(signalEvents)
// 				if performance > loss {
// 					performance = loss
// 					bestPeriod = period
// 					bestBuyThread = buyThread
// 					bestSellThread = sellThread
// 				}
// 			}
// 		}

// 	}

// 	fmt.Println("損失", performance, "最適なピリオド", bestPeriod, "最適な買いライン", bestBuyThread, "最適な売りライン", bestSellThread)

// 	return performance, bestPeriod, bestBuyThread, bestSellThread
// }

// func (df *DataFrameCandle) OptimizeRsiProfitFactor() (performance float64, bestPeriod int, bestBuyThread, bestSellThread float64) {
// 	bestPeriod = 14
// 	bestBuyThread, bestSellThread = 20.0, 80.0

// 	for period := 4; period < 30; period++ {
// 		for buyThread := 30.0; buyThread > 10; buyThread -= 1 {

// 			for sellThread := 70.0; sellThread < 96; sellThread += 1 {
// 				signalEvents := df.RsiStrategy(period, buyThread, sellThread, accountBlance)
// 				if signalEvents == nil {
// 					continue
// 				}

// 				if analytics.TotalTrades(signalEvents) < 40 {
// 					continue
// 				}

// 				profitFactor := analytics.ProfitFactor(signalEvents)
// 				if performance < profitFactor {
// 					performance = profitFactor
// 					bestPeriod = period
// 					bestBuyThread = buyThread
// 					bestSellThread = sellThread
// 				}
// 			}
// 		}

// 	}

// 	fmt.Println("プロフィットファクター", performance, "最適なピリオド", bestPeriod, "最適な買いライン", bestBuyThread, "最適な売りライン", bestSellThread)

// 	return performance, bestPeriod, bestBuyThread, bestSellThread
// }

// func (df *DataFrameCandle) OptimizeRsiPayOffRatio() (performance float64, bestPeriod int, bestBuyThread, bestSellThread float64) {
// 	bestPeriod = 14
// 	bestBuyThread, bestSellThread = 20.0, 80.0

// 	for period := 4; period < 25; period++ {
// 		for buyThread := 30.0; buyThread > 10; buyThread -= 1 {

// 			for sellThread := 75.0; sellThread < 96; sellThread += 1 {
// 				signalEvents := df.RsiStrategy(period, buyThread, sellThread, accountBlance)
// 				if signalEvents == nil {
// 					continue
// 				}

// 				if analytics.TotalTrades(signalEvents) < 40 {
// 					continue
// 				}

// 				payOffRatio := analytics.PayOffRatio(signalEvents)
// 				if performance < payOffRatio {
// 					performance = payOffRatio
// 					bestPeriod = period
// 					bestBuyThread = buyThread
// 					bestSellThread = sellThread
// 				}
// 			}
// 		}

// 	}

// 	fmt.Println("プロフィットファクター", performance, "最適なピリオド", bestPeriod, "最適な買いライン", bestBuyThread, "最適な売りライン", bestSellThread)

// 	return performance, bestPeriod, bestBuyThread, bestSellThread
// }

// func (df *DataFrameCandle) OptimizeRsiSharpRatio() (performance float64, bestPeriod int, bestBuyThread, bestSellThread float64) {
// 	bestPeriod = 14
// 	bestBuyThread, bestSellThread = 20.0, 80.0

// 	for period := 5; period < 30; period++ {
// 		for buyThread := 25.0; buyThread > 10; buyThread -= 1 {

// 			for sellThread := 75.0; sellThread < 95; sellThread += 1 {
// 				signalEvents := df.RsiStrategy(period, buyThread, sellThread, accountBlance)
// 				if signalEvents == nil {
// 					continue
// 				}

// 				if analytics.TotalTrades(signalEvents) < 40 {
// 					continue
// 				}

// 				sharpeRatio := analytics.SharpeRatio(signalEvents, 0.06)
// 				if performance < sharpeRatio {
// 					performance = sharpeRatio
// 					bestPeriod = period
// 					bestBuyThread = buyThread
// 					bestSellThread = sellThread
// 				}
// 			}
// 		}

// 	}

// 	fmt.Println("シャープレシオ", performance, "最適なピリオド", bestPeriod, "最適な買いライン", bestBuyThread, "最適な売りライン", bestSellThread)

// 	return performance, bestPeriod, bestBuyThread, bestSellThread
// }
