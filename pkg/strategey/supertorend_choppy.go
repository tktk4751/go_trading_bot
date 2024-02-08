package strategey

import (
	"fmt"
	"runtime"
	"sync"
	"v1/pkg/analytics"
	"v1/pkg/execute"
	"v1/pkg/indicator/indicators"
	"v1/pkg/management/risk"
	"v1/pkg/trader"
)

func (df *DataFrameCandle) SuperTrendChoppyStrategy(atrPeriod int, factor float64, choppy int, duration int, account *trader.Account, simple bool) *execute.SignalEvents {

	var StrategyName = "SUPERTREND_CHOPPY"
	// var err error

	lenCandles := len(df.Candles)

	if lenCandles <= atrPeriod {
		return nil
	}

	signalEvents := execute.NewSignalEvents()

	h := df.Highs()
	l := df.Lows()
	c := df.Closes()
	// hlc3 := df.Hlc3()

	superTrend, _ := indicators.SuperTrend(atrPeriod, factor, h, l, c)

	// stUp := superTrend.UpperBand
	// stLow := superTrend.UpperBand
	st := superTrend.SuperTrend

	// rsi := talib.Rsi(hlc3, 14)

	buySize := 0.0
	buyPrice := 0.0
	slRatio := 0.9

	index := risk.ChoppySlice(duration, c, h, l)
	choppyEma := risk.ChoppyEma(index, choppy)

	isHolding := false

	for i := 1; i < len(choppyEma); i++ {

		if i < atrPeriod {
			// fmt.Printf("Skipping iteration %d due to insufficient data.\n", i)
			continue
		}
		if (c[i-1] < st[i-1] && c[i] >= st[i]) && choppyEma[i] > 50 && !isHolding {

			// fee := 1 - 0.01
			if simple {
				buySize = account.SimpleTradeSize(1)
				buyPrice = c[i]
				accountBalance := account.GetBalance()

				signalEvents.Buy(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, buySize, accountBalance, false)
				isHolding = true

			} else {
				buySize = account.TradeSize(riskSize) / df.Candles[i].Close
				buyPrice = c[i]
				accountBalance := account.GetBalance()
				if account.Buy(df.Candles[i].Close, buySize) {
					signalEvents.Buy(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, buySize, accountBalance, false)
					isHolding = true
				}
			}

		}
		if ((c[i-1] > st[i-1] && c[i] <= st[i]) || (c[i] <= buyPrice*slRatio)) && isHolding {
			if simple {
				accountBalance := 1000.0

				signalEvents.Sell(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, buySize, accountBalance, false)
				isHolding = false

			} else {
				accountBalance := account.GetBalance()
				if account.Sell(df.Candles[i].Close) {
					signalEvents.Sell(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, buySize, accountBalance, false)
					isHolding = false
					buySize = 0.0
					account.PositionSize = buySize

				}
			}
		}
	}

	// fmt.Println(signalEvents)
	return signalEvents
}

func (df *DataFrameCandle) OptimizeSuperTrend() (performance float64, bestAtrPeriod int, bestFactor float64, bestChoppy int, bestDuration int) {
	runtime.GOMAXPROCS(10)
	bestAtrPeriod = 21
	bestFactor = 3.0
	bestChoppy = 13
	bestDuration = 30

	limit := 3000
	slots := make(chan struct{}, limit)

	// a := trader.NewAccount(1000)
	// marketDefault, _ := BuyAndHoldingStrategy(a)

	var mu sync.Mutex
	var wg sync.WaitGroup

	for atrPeriod := 5; atrPeriod < 40; atrPeriod += 2 {
		for factor := 2.0; factor < 8.0; factor += 0.5 {
			for choppy := 5; choppy < 18; choppy += 2 {
				for duration := 10; duration < 200; duration += 10 {

					wg.Add(1)
					slots <- struct{}{}

					go func(atrPeriod int, factor float64, choppy int, duration int) {
						defer wg.Done()
						account := trader.NewAccount(1000) // Move this line inside the goroutine
						signalEvents := df.SuperTrendChoppyStrategy(atrPeriod, factor, choppy, duration, account, simple)

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
							bestAtrPeriod = atrPeriod
							bestFactor = factor
							bestChoppy = choppy
							bestDuration = duration

						}
						<-slots
						mu.Unlock()

					}(atrPeriod, factor, choppy, duration)

				}
			}
		}
	}

	wg.Wait()

	fmt.Println("最高のパフォーマンス", performance, "最適なATR", bestAtrPeriod, "最適なファクター", bestFactor, "最適なチョッピー", bestChoppy, "最適なチョッピー期間", bestDuration)

	return performance, bestAtrPeriod, bestFactor, bestChoppy, bestDuration
}

func RunSTOptimize() {

	df, account, _ := RadyBacktest()

	performance, bestAtrPeriod, bestFactor, bestChoppy, bestDuration := df.OptimizeSuperTrend()

	if performance > 0 {

		df.Signal = df.SuperTrendChoppyStrategy(bestAtrPeriod, bestFactor, bestChoppy, bestDuration, account, simple)
		Result(df.Signal)

	} else {
		fmt.Println("💸マイナスです")
		df.Signal = df.SuperTrendChoppyStrategy(bestAtrPeriod, bestFactor, bestChoppy, bestDuration, account, simple)
		Result(df.Signal)

	}

}

func SuperTrendBacktest() {

	df, account, _ := RadyBacktest()

	df.Signal = df.SuperTrendChoppyStrategy(13, 7.5, 11, 120, account, simple)
	Result(df.Signal)

}
