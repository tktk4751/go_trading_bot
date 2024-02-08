package strategey

import (
	"fmt"
	"runtime"
	"sync"
	"v1/pkg/analytics"
	"v1/pkg/execute"
	"v1/pkg/management/risk"
	"v1/pkg/trader"

	"github.com/markcheno/go-talib"
)

func (df *DataFrameCandle) EmaChoppyStrategy(period1, period2 int, choppy int, duration int, account *trader.Account) *execute.SignalEvents {

	var StrategyName = "EMA_CHOPPY"
	lenCandles := len(df.Candles)
	if lenCandles <= period1 || lenCandles <= period2 {
		return nil
	}
	signalEvents := execute.NewSignalEvents()

	emaValue1 := talib.Ema(df.Hlc3(), period1)
	emaValue2 := talib.Ema(df.Hlc3(), period2)
	// rsiValue := talib.Rsi(df.Closes(), 14)

	buySize := 0.0
	buyPrice := 0.0
	slRatio := 0.9
	index := risk.ChoppySlice(duration, df.Closes(), df.Highs(), df.Lows())
	choppyEma := risk.ChoppyEma(index, choppy)

	isBuyHolding := false
	for i := 1; i < lenCandles; i++ {
		if i < period1 || i < period2 || i >= len(choppyEma) {
			continue
		}

		if emaValue1[i-1] < emaValue2[i-1] && emaValue1[i] >= emaValue2[i] && choppyEma[i] > 50 && !isBuyHolding {

			accountBalance := account.GetBalance()
			// fee := 1 - 0.01

			buySize = account.TradeSize(riskSize) / df.Candles[i].Close
			buyPrice = df.Candles[i].Close
			if account.Buy(df.Candles[i].Close, buySize) {

				signalEvents.Buy(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, buySize, accountBalance, false)
				isBuyHolding = true

			}
		}
		if emaValue1[i-1] > emaValue2[i-1] && emaValue1[i] <= emaValue2[i] || (df.Candles[i].Close <= buyPrice*slRatio) && isBuyHolding {
			accountBalance := account.GetBalance()
			if account.Sell(df.Candles[i].Close) {
				signalEvents.Sell(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, buySize, accountBalance, false)
				isBuyHolding = false
				buySize = 0.0
				account.PositionSize = buySize

			}
		}
	}
	return signalEvents
}

func (df *DataFrameCandle) OptimizeEmaChoppy() (performance float64, bestPeriod1 int, bestPeriod2 int, bestChoppy int, bestDuration int) {
	runtime.GOMAXPROCS(10)
	bestPeriod1 = 5
	bestPeriod2 = 21
	bestDuration = 30

	limit := 3000
	slots := make(chan struct{}, limit)

	// a := trader.NewAccount(1000)
	// marketDefault, _ := BuyAndHoldingStrategy(a)

	var mu sync.Mutex
	var wg sync.WaitGroup

	for period1 := 3; period1 < 34; period1 += 2 {
		for period2 := 13; period2 < 89; period2 += 4 {
			for choppy := 5; choppy < 18; choppy += 1 {
				for duration := 10; duration < 200; duration += 10 {

					wg.Add(1)
					slots <- struct{}{}

					go func(period1 int, period2 int, choppy int, duration int) {
						defer wg.Done()
						account := trader.NewAccount(1000) // Move this line inside the goroutine
						signalEvents := df.EmaChoppyStrategy(period1, period2, choppy, duration, account)

						if signalEvents == nil {
							return
						}

						if analytics.TotalTrades(signalEvents) < 10 {
							<-slots
							return
						}

						// if analytics.NetProfit(signalEvents) < marketDefault {
						// 	<-slots
						// 	return
						// }

						// if analytics.SQN(signalEvents) < 3.2 {
						// 	<-slots
						// 	return
						// }

						// if analytics.PayOffRatio(signalEvents) < 1 {
						// <-slots

						// 	return
						// }

						p := analytics.SortinoRatio(signalEvents, 0.02)
						// p := analytics.SQN(signalEvents)
						mu.Lock()
						if performance == 0 || performance < p {
							performance = p
							bestPeriod1 = period1
							bestPeriod2 = period2
							bestChoppy = choppy
							bestDuration = duration

						}
						<-slots
						mu.Unlock()

					}(period1, period2, 13, duration)

				}
			}
		}
	}

	wg.Wait()

	fmt.Println("最高パフォーマンス", performance, "最適な短期線", bestPeriod1, "最適な長期線", bestPeriod2, "最適なチョッピーEMA", bestChoppy, "最適なチョッピー期間", bestDuration)

	return performance, bestPeriod1, bestPeriod2, bestChoppy, bestDuration
}

func RunEmaOptimize() {

	df, account, _ := RadyBacktest()

	performance, bestPeriod1, bestPeriod2, bestChoppy, bestDuration := df.OptimizeEmaChoppy()

	if performance > 0 {

		df.Signal = df.EmaChoppyStrategy(bestPeriod1, bestPeriod2, bestChoppy, bestDuration, account)
		Result(df.Signal)

	}

}

func EmaBacktest() {

	df, account, _ := RadyBacktest()

	df.Signal = df.EmaChoppyStrategy(11, 13, 13, 20, account)
	Result(df.Signal)

}
