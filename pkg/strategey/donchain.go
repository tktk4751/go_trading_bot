package strategey

import (
	"fmt"
	"log"
	"sync"
	"v1/pkg/analytics"
	"v1/pkg/config"
	"v1/pkg/execute"
	"v1/pkg/indicator/indicators"
	"v1/pkg/trader"
)

// func getStrageyNameDonchain() string {
// 	return "DBO"
// }

func (df *DataFrameCandle) DonchainStrategy(period int, account *trader.Account) *execute.SignalEvents {
	var StrategyName = "DBO"

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

	for i := 1; i < lenCandles; i++ {

		if i < period {
			continue
		}

		if close[i] > donchain.High[i-1] && !isHolding {

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
func (df *DataFrameCandle) OptimizeDonchainGoroutin() (performance float64, bestPeriod int) {

	bestPeriod = 40
	var mu sync.Mutex
	var wg sync.WaitGroup

	// a := trader.NewAccount(1000)
	// marketDefault, _ := BuyAndHoldingStrategy(a)

	for period := 10; period < 333; period++ {
		wg.Add(1)
		go func(period int) {
			defer wg.Done()
			account := trader.NewAccount(1000)
			signalEvents := df.DonchainStrategy(period, account)

			if signalEvents == nil {
				return
			}

			if analytics.TotalTrades(signalEvents) < 3 {
				return
			}

			// if analytics.NetProfit(signalEvents) < marketDefault {
			// 	return
			// }

			// if analytics.WinRate(signalEvents) < 0.45 {
			// 	return
			// }

			// if analytics.ProfitFactor(signalEvents) < 3 {
			// 	return
			// }

			pf := analytics.SQN(signalEvents)
			mu.Lock()
			if performance < pf {
				performance = pf
				bestPeriod = period
			}
			mu.Unlock()
		}(period)
	}

	wg.Wait()

	fmt.Println("最高利益", performance, "最適なピリオド", bestPeriod)

	return performance, bestPeriod
}

func RunBacktestDonchain() {

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

	performancePayOffRatio, bestPayOffRatioPeriod := df.OptimizeDonchainGoroutin()

	if performancePayOffRatio > 0 {

		df.Signal = df.DonchainStrategy(bestPayOffRatioPeriod, account)
		Result(df.Signal)

	}

}
