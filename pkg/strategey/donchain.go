package strategey

import (
	"fmt"
	"log"
	"math"
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

	donchain := indicators.Donchain(df.Highs(), df.Low(), period)
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
			if account.Sell(df.Candles[i].Close, 0.0) {
				signalEvents.Sell(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, buySize, accountBalance, false)
				isHolding = false
				buySize = 0.0
				account.PositionSize = buySize

			}
		}

	}
	return signalEvents

}

func (df *DataFrameCandle) OptimizeDonchainProfit() (performance float64, bestPeriod int) {
	if df == nil {
		return 0.0, 0
	}

	account := trader.NewAccount(1000)

	bestPeriod = 40

	for period := 5; period < 350; period++ {

		signalEvents := df.DonchainStrategy(period, account)
		if signalEvents == nil {
			continue
		}
		profit := analytics.NetProfit(signalEvents)
		if performance < profit {
			performance = profit
			bestPeriod = period

		}

	}

	fmt.Println("最高利益", performance, "最適なピリオド", bestPeriod)

	return performance, bestPeriod
}

func (df *DataFrameCandle) OptimizeDonchainLoss() (performance float64, bestPeriod int) {
	if df == nil {
		return 0.0, 0
	}

	account := trader.NewAccount(1000)

	bestPeriod = 40
	performance = math.MaxFloat64

	for period := 5; period < 350; period++ {

		signalEvents := df.DonchainStrategy(period, account)
		if signalEvents == nil {
			continue
		}
		loss := analytics.Loss(signalEvents)
		if performance > loss {
			performance = loss
			bestPeriod = period

		}

	}

	fmt.Println("最高利益", performance, "最適なピリオド", bestPeriod)

	return performance, bestPeriod
}

func (df *DataFrameCandle) OptimizeDonchainWinRate() (performance float64, bestPeriod int) {
	bestPeriod = 40

	account := trader.NewAccount(1000)

	for period := 10; period < 333; period++ {

		signalEvents := df.DonchainStrategy(period, account)
		if signalEvents == nil {
			continue
		}
		winrate := analytics.WinRate(signalEvents)
		if performance < winrate {
			performance = winrate
			bestPeriod = period

		}

	}

	fmt.Println("最高勝率", performance*100, "% ", "最適なピリオド", bestPeriod)
	return performance, bestPeriod
}

func (df *DataFrameCandle) OptimizeDonchainProfitFactor() (performance float64, bestPeriod int) {

	account := trader.NewAccount(1000)
	bestPeriod = 40

	for period := 10; period < 333; period++ {

		signalEvents := df.DonchainStrategy(period, account)
		if signalEvents == nil {
			continue
		}
		pf := analytics.ProfitFactor(signalEvents)
		if performance < pf {
			performance = pf
			bestPeriod = period

		}

	}

	fmt.Println("プロフィットファクター", performance, "最適なピリオド", bestPeriod)

	return performance, bestPeriod
}

func (df *DataFrameCandle) OptimizeDonchainPayOffRatio() (performance float64, bestPeriod int) {

	account := trader.NewAccount(1000)
	bestPeriod = 40

	for period := 10; period < 333; period++ {

		signalEvents := df.DonchainStrategy(period, account)
		if signalEvents == nil {
			continue
		}
		pf := analytics.PayOffRatio(signalEvents)
		if performance < pf {
			performance = pf
			bestPeriod = period

		}

	}

	fmt.Println("ペイオフレシオ", performance, "最適なピリオド", bestPeriod)

	return performance, bestPeriod
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

			// if analytics.TotalTrades(signalEvents) < 20 {
			// 	return
			// }

			// if analytics.NetProfit(signalEvents) < marketDefault {
			// 	return
			// }

			// if analytics.WinRate(signalEvents) < 0.45 {
			// 	return
			// }

			// if analytics.ProfitFactor(signalEvents) < 3 {
			// 	return
			// }

			pf := analytics.ProfitFactor(signalEvents)
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

	fmt.Println(btcfg.AssetName)

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
