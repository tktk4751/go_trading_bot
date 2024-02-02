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

func getStrageyNameBb() string {
	return "BB"
}

func (df *DataFrameCandle) BbStrategy(n int, k float64, account *trader.Account) *execute.SignalEvents {

	var StrategyName = "BB"
	lenCandles := len(df.Candles)

	if lenCandles <= n {
		return nil
	}

	signalEvents := execute.NewSignalEvents()
	bbUp, _, bbDown := talib.BBands(df.Closes(), n, k, k, 0)

	buySize := 0.0
	isBuyHolding := false
	for i := 1; i < lenCandles; i++ {
		if i < n {
			continue
		}
		if bbDown[i-1] > df.Candles[i-1].Close && bbDown[i] <= df.Candles[i].Close && !isBuyHolding {
			buySize = account.TradeSize(0.9) / df.Candles[i].Close
			accountBalance := account.GetBalance()
			signalEvents.Buy(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, buySize, accountBalance, false)

			isBuyHolding = true
		}
		if bbUp[i-1] < df.Candles[i-1].Close && bbUp[i] >= df.Candles[i].Close && isBuyHolding {
			accountBalance := account.GetBalance()
			signalEvents.Sell(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, buySize, accountBalance, false)
			isBuyHolding = false
		}
	}
	return signalEvents
}

func (df *DataFrameCandle) OptimizeBbGoroutin() (performance float64, bestN int, bestK float64) {
	runtime.GOMAXPROCS(10)

	bestN = 20
	bestK = 2.0

	a := trader.NewAccount(1000)

	marketDefault, _ := BuyAndHoldingStrategy(a)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for n := 10; n < 300; n++ {
		for k := 2.0; k < 3.5; k += 0.1 {

			wg.Add(1)
			go func(n int, k float64) {
				defer wg.Done()
				account := trader.NewAccount(1000) // Move this line inside the goroutine
				signalEvents := df.BbStrategy(n, k, account)

				if signalEvents == nil {
					return
				}

				if analytics.TotalTrades(signalEvents) < 20 {
					return
				}

				if analytics.NetProfit(signalEvents) < marketDefault {
					return
				}

				// if analytics.WinRate(signalEvents) < 0.50 {
				// 	return
				// }

				// if analytics.PayOffRatio(signalEvents) < 1 {
				// 	return
				// }

				p := analytics.NetProfit(signalEvents)
				mu.Lock()
				if performance == 0 || performance < p {
					performance = p
					bestN = n
					bestK = k

				}
				mu.Unlock()
			}(n, k)

		}
	}

	wg.Wait()

	fmt.Println("最高利益", performance, "最適なN", bestN, "最適なK", bestK)

	return performance, bestN, bestK
}
func RunBacktestBb() {

	var err error

	// account := trader.NewAccount(1000)
	btcfg, err := config.Yaml()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Println(btcfg.AssetName)

	assetName := btcfg.AssetName
	duration := btcfg.Dration
	// limit := btcfg.Limit

	account := trader.NewAccount(1000)

	df, _ := GetCandleData(assetName, duration)

	performance, bestN, bestK := df.OptimizeBbGoroutin()

	if performance > 0 {

		df.Signal = df.BbStrategy(bestN, bestK, account)
		Result(df.Signal)

	}

}
