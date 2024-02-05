package strategey

import (
	"fmt"
	"log"
	"sync"
	"v1/pkg/analytics"
	"v1/pkg/config"
	"v1/pkg/execute"
	"v1/pkg/indicator/indicators"
	"v1/pkg/management/risk"
	"v1/pkg/trader"
)

// func getStrageyNameDonchain() string {
// 	return "DBO"
// }

func (df *DataFrameCandle) DonchainChoppyStrategy(period int, choppy int, account *trader.Account) *execute.SignalEvents {
	var StrategyName = "DBO_CHOPPY"

	lenCandles := len(df.Candles)
	if lenCandles <= period {
		return nil
	}

	signalEvents := execute.NewSignalEvents()

	donchain := indicators.Donchain(df.Highs(), df.Lows(), period)
	// atr := talib.Atr(df.Highs(), df.Low(), df.Closes(), 21)

	close := df.Closes()

	buySize := 0.0
	isHolding := false

	index := risk.ChoppySlice(df.Closes(), df.Highs(), df.Lows())
	choppyEma := risk.ChoppyEma(index, choppy)

	for i := 30; i < lenCandles; i++ {

		if i < period || i >= len(choppyEma) {
			continue
		}

		if close[i] > donchain.High[i-1] && choppyEma[i] > 50 && !isHolding {

			buySize = account.TradeSize(riskSize) / df.Candles[i].Close
			accountBalance := account.GetBalance()
			if account.Buy(df.Candles[i].Close, buySize) {
				signalEvents.Buy(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, buySize, accountBalance, false)
				isHolding = true
			}
		}
		if close[i] < donchain.Low[i-1] && isHolding {
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

func (df *DataFrameCandle) OptimizeDonchainChoppyGoroutin() (performance float64, bestPeriod int, bestChoppy int) {

	bestPeriod = 40
	bestChoppy = 13
	var mu sync.Mutex
	var wg sync.WaitGroup

	a := trader.NewAccount(1000)
	marketDefault, _ := BuyAndHoldingStrategy(a)

	limit := 1000
	slots := make(chan struct{}, limit)

	for period := 10; period < 250; period += 3 {
		for choppy := 5; choppy < 18; choppy += 1 {
			wg.Add(1)
			slots <- struct{}{}

			go func(period int, choppy int) {
				defer wg.Done()
				account := trader.NewAccount(1000)
				signalEvents := df.DonchainChoppyStrategy(period, choppy, account)

				if signalEvents == nil {
					return
				}

				if analytics.TotalTrades(signalEvents) < 3 {
					<-slots
					return
				}

				if analytics.NetProfit(signalEvents) < marketDefault {
					// <-slots
					return
				}

				// if analytics.WinRate(signalEvents) < 0.45 {
				// <-slots
				// 	return
				// }

				// if analytics.ProfitFactor(signalEvents) < 3 {
				// <-slots
				// 	return
				// }

				pf := analytics.SQN(signalEvents)
				mu.Lock()
				if performance < pf {
					performance = pf
					bestPeriod = period
					bestChoppy = choppy
				}
				<-slots
				mu.Unlock()
			}(period, choppy)
		}
	}

	wg.Wait()

	fmt.Println("最高利益", performance, "最適なピリオド", bestPeriod, "最適なチョッピー", bestChoppy)

	return performance, bestPeriod, bestChoppy
}

func RunBacktestDonchainChoppy() {

	var err error

	// account := trader.NewAccount(1000)
	btcfg, err := config.Yaml()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Println("--------------------------------------------")

	// strategyName := getStrageyNameDonchain()
	assetName := btcfg.AssetName
	duration := btcfg.Dration

	// limit := btcfg.Limit

	account := trader.NewAccount(1000)

	df, _ := GetCandleData(assetName, duration)

	p, bestPeriod, bestChoppy := df.OptimizeDonchainChoppyGoroutin()

	if p > 0 {

		df.Signal = df.DonchainChoppyStrategy(bestPeriod, bestChoppy, account)
		Result(df.Signal)

	}

}
