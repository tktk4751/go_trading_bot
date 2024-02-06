package strategey

import (
	"fmt"
	"sync"
	"v1/pkg/analytics"
	"v1/pkg/execute"
	"v1/pkg/indicator/indicators"
	"v1/pkg/management/risk"
	"v1/pkg/trader"

	"github.com/markcheno/go-talib"
)

// func getStrageyNameDonchain() string {
// 	return "DBO"
// }

func (df *DataFrameCandle) DonchainChoppyStrategy(period int, choppy int, duration int, account *trader.Account) *execute.SignalEvents {
	var StrategyName = "DBO_CHOPPY"

	lenCandles := len(df.Candles)
	if lenCandles <= period {
		return nil
	}

	signalEvents := execute.NewSignalEvents()

	donchain := indicators.Donchain(df.Highs(), df.Lows(), period)
	// atr := talib.Atr(df.Highs(), df.Low(), df.Closes(), 21)

	ema := talib.Ema(df.Hlc3(), 89)

	close := df.Closes()

	buySize := 0.0
	buyPrice := 0.0
	slRatio := 0.9
	isHolding := false

	index := risk.ChoppySlice(duration, df.Closes(), df.Highs(), df.Lows())
	choppyEma := risk.ChoppyEma(index, choppy)

	for i := 30; i < lenCandles; i++ {

		if i < period || i >= len(choppyEma) {
			continue
		}

		if close[i] > donchain.High[i-1] && choppyEma[i] > 50 && close[i] > ema[i] && !isHolding {

			buySize = account.TradeSize(riskSize) / df.Candles[i].Close
			buyPrice = close[i]
			accountBalance := account.GetBalance()
			if account.Buy(df.Candles[i].Close, buySize) {
				signalEvents.Buy(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, buySize, accountBalance, false)
				isHolding = true
			}
		}
		if (close[i] < donchain.Low[i-1] || (close[i] <= buyPrice*slRatio)) && isHolding {
			accountBalance := account.GetBalance()
			if account.Sell(df.Candles[i].Close) {
				signalEvents.Sell(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, buySize, accountBalance, false)
				isHolding = false
				buySize = 0.0
				account.PositionSize = buySize

			}
		}

	}
	return signalEvents

}

func (df *DataFrameCandle) OptimizeDonchainChoppyGoroutin() (performance float64, bestPeriod int, bestChoppy int, bestDuration int) {

	bestPeriod = 40
	bestChoppy = 13
	var mu sync.Mutex
	var wg sync.WaitGroup

	// a := trader.NewAccount(1000)
	// marketDefault, _ := BuyAndHoldingStrategy(a)

	limit := 3000
	slots := make(chan struct{}, limit)

	for period := 5; period < 250; period += 10 {
		for duration := 10; duration < 200; duration += 10 {
			for choppy := 6; choppy < 18; choppy += 2 {
				wg.Add(1)
				slots <- struct{}{}

				go func(period int, choppy int, duration int) {
					defer wg.Done()
					account := trader.NewAccount(1000)
					signalEvents := df.DonchainChoppyStrategy(period, choppy, duration, account)

					if signalEvents == nil {
						return
					}

					if analytics.TotalTrades(signalEvents) < 30 {
						<-slots
						return
					}

					// if analytics.NetProfit(signalEvents) < marketDefault {
					// 	// <-slots
					// 	return
					// }

					// if analytics.SQN(signalEvents) < 3.2 {
					// 	<-slots
					// 	return
					// }

					// if analytics.ProfitFactor(signalEvents) < 3 {
					// <-slots
					// 	return
					// }

					// pf := analytics.SortinoRatio(signalEvents, 0.02)
					pf := analytics.SQN(signalEvents)
					mu.Lock()
					if performance < pf {
						performance = pf
						bestPeriod = period
						bestChoppy = choppy
						bestDuration = duration
					}
					<-slots
					mu.Unlock()
				}(period, choppy, duration)
			}
		}
	}

	wg.Wait()

	fmt.Println("最高パフォーマンス", performance, "最適なピリオド", bestPeriod, "最適なチョッピー", bestChoppy, "最適なチョッピー期間", bestDuration)

	return performance, bestPeriod, bestChoppy, bestDuration
}

func RunDonchainOptimize() {

	df, account, _ := RadyBacktest()

	p, bestPeriod, bestChoppy, bestDuration := df.OptimizeDonchainChoppyGoroutin()

	if p > 0 {

		df.Signal = df.DonchainChoppyStrategy(bestPeriod, bestChoppy, bestDuration, account)
		Result(df.Signal)

	}

}

func DonchainBacktest() {

	df, account, _ := RadyBacktest()

	df.Signal = df.DonchainChoppyStrategy(215, 12, 60, account)
	Result(df.Signal)
}
