package strategey

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"v1/pkg/analytics"
	"v1/pkg/config"
	"v1/pkg/execute"
	"v1/pkg/trader"

	"github.com/markcheno/go-talib"
)

func (df *DataFrameCandle) MacdStrategy(macdFastPeriod, macdSlowPeriod, macdSignalPeriod int, account *trader.Account) *execute.SignalEvents {
	var StrategyName = "MACD"

	lenCandles := len(df.Candles)

	if lenCandles <= macdFastPeriod || lenCandles <= macdSlowPeriod || lenCandles <= macdSignalPeriod {
		return nil
	}

	signalEvents := execute.NewSignalEvents()
	outMACD, outMACDSignal, _ := talib.Macd(df.Closes(), macdFastPeriod, macdSlowPeriod, macdSignalPeriod)

	buySize := 0.0
	buyPrice := 0.0
	slRatio := 0.9
	isBuyHolding := false

	for i := 1; i < lenCandles; i++ {
		if outMACD[i] < 0 &&
			outMACDSignal[i] < 0 &&
			outMACD[i-1] < outMACDSignal[i-1] &&
			outMACD[i] >= outMACDSignal[i] &&
			!isBuyHolding {

			accountBalance := account.GetBalance()
			buySize = account.TradeSize(riskSize) / df.Candles[i].Close
			buyPrice = df.Candles[i].Close

			if account.Buy(df.Candles[i].Close, buySize) {

				signalEvents.Buy(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, buySize, accountBalance, false)
				isBuyHolding = true

			}
		}

		if outMACD[i] > 0 &&
			outMACDSignal[i] > 0 &&
			outMACD[i-1] > outMACDSignal[i-1] &&
			outMACD[i] <= outMACDSignal[i] ||
			(df.Candles[i].Close <= buyPrice*slRatio) && isBuyHolding {
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

func (df *DataFrameCandle) OptimizeMacd() (performance float64, bestMacdFastPeriod, bestMacdSlowPeriod, bestMacdSignalPeriod int) {
	runtime.GOMAXPROCS(10)

	bestMacdFastPeriod = 12
	bestMacdSlowPeriod = 26
	bestMacdSignalPeriod = 9

	limit := 1000
	slots := make(chan struct{}, limit)

	// a := trader.NewAccount(1000)
	// marketDefault, _ := BuyAndHoldingStrategy(a)

	var mu sync.Mutex
	var wg sync.WaitGroup

	for fastPeriod := 3; fastPeriod < 13; fastPeriod++ {
		for slowPeriod := 5; slowPeriod < 60; slowPeriod++ {
			for signalPeriod := 3; signalPeriod < 21; signalPeriod++ {
				wg.Add(1)
				slots <- struct{}{}
				go func(fastPeriod int, slowPeriod int, signalPeriod int) {
					defer wg.Done()
					account := trader.NewAccount(1000) // Move this line inside the goroutine
					signalEvents := df.MacdStrategy(fastPeriod, slowPeriod, signalPeriod, account)

					if signalEvents == nil {
						<-slots
						return
					}

					// if analytics.TotalTrades(signalEvents) < 20 {
					// 	return
					// }

					// if analytics.NetProfit(signalEvents) < marketDefault {
					// 	<-slots
					// 	return
					// }

					// if analytics.WinRate(signalEvents) < 0.50 {
					// 	return
					// }

					// if analytics.PayOffRatio(signalEvents) < 1 {
					// 	return
					// }

					p := analytics.ProfitFactor(signalEvents)
					mu.Lock()
					if performance == 0 || performance < p {
						performance = p
						bestMacdFastPeriod = fastPeriod
						bestMacdSlowPeriod = slowPeriod
						bestMacdSignalPeriod = signalPeriod
					}
					<-slots
					mu.Unlock()

				}(fastPeriod, slowPeriod, signalPeriod)
			}
		}
	}

	wg.Wait()

	fmt.Println("最高利益", performance, "最適なfastPeriod", bestMacdFastPeriod, "最適なslowPeriod", bestMacdSlowPeriod, "最適なsignalPeriod", bestMacdSignalPeriod)

	return performance, bestMacdFastPeriod, bestMacdSlowPeriod, bestMacdSignalPeriod
}

func RunBacktestMacd() {

	var err error

	// account := trader.NewAccount(1000)
	btcfg, err := config.Yaml()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Println(btcfg.AssetName)

	assetName := btcfg.AssetName
	duration := btcfg.Dration

	account := trader.NewAccount(1000)

	df, _ := GetCandleData(assetName, duration)

	performance, bestMacdFastPeriod, bestMacdSlowPeriod, bestMacdSignalPeriod := df.OptimizeMacd()

	if performance > 0 {

		df.Signal = df.MacdStrategy(bestMacdFastPeriod, bestMacdSlowPeriod, bestMacdSignalPeriod, account)
		Result(df.Signal)

	}

}
