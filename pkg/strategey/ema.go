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

func (df *DataFrameCandle) EmaStrategy(period1, period2 int, account *trader.Account) *execute.SignalEvents {

	var StrategyName = "EMA"
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

	isBuyHolding := false
	for i := 1; i < lenCandles; i++ {
		if i < period1 || i < period2 {
			continue
		}
		if emaValue1[i-1] < emaValue2[i-1] && emaValue1[i] >= emaValue2[i] && !isBuyHolding {
			accountBalance := account.GetBalance()
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

func (df *DataFrameCandle) OptimizeEma() (performance float64, bestPeriod1 int, bestPeriod2 int) {
	runtime.GOMAXPROCS(10)
	bestPeriod1 = 5
	bestPeriod2 = 21

	limit := 1000
	slots := make(chan struct{}, limit)

	a := trader.NewAccount(1000)
	marketDefault, _ := BuyAndHoldingStrategy(a)

	var mu sync.Mutex
	var wg sync.WaitGroup

	for period1 := 3; period1 < 100; period1 += 3 {
		for period2 := 5; period2 < 240; period2 += 3 {

			wg.Add(1)
			slots <- struct{}{}

			go func(period1 int, period2 int) {
				defer wg.Done()
				account := trader.NewAccount(1000) // Move this line inside the goroutine
				signalEvents := df.EmaStrategy(period1, period2, account)

				if signalEvents == nil {
					return
				}

				if analytics.TotalTrades(signalEvents) < 3 {
					<-slots
					return
				}

				if analytics.NetProfit(signalEvents) < marketDefault {
					<-slots
					return
				}

				// if analytics.WinRate(signalEvents) < 0.50 {
				// <-slots

				// 	return
				// }

				// if analytics.PayOffRatio(signalEvents) < 1 {
				// <-slots

				// 	return
				// }

				p := analytics.ProfitFactor(signalEvents)
				mu.Lock()
				if performance == 0 || performance < p {
					performance = p
					bestPeriod1 = period1
					bestPeriod2 = period2

				}
				<-slots
				mu.Unlock()

			}(period1, period2)

		}
	}

	wg.Wait()

	fmt.Println("最高利益", performance, "最適な短期線", bestPeriod1, "最適な長期線", bestPeriod2)

	return performance, bestPeriod1, bestPeriod2
}

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

func RunBacktestEma() {

	var err error

	// account := trader.NewAccount(1000)
	btcfg, err := config.Yaml()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Println("--------------------------------------------")

	assetName := btcfg.AssetName
	duration := btcfg.Dration
	// limit := btcfg.Limit

	account := trader.NewAccount(1000)

	df, _ := GetCandleData(assetName, duration)

	performance, bestPeriod1, bestPeriod2 := df.OptimizeEma()

	if performance > 0 {

		df.Signal = df.EmaStrategy(bestPeriod1, bestPeriod2, account)
		Result(df.Signal)

	}

}
