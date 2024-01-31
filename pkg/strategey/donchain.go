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

func getStrageyNameDonchain() string {
	return "DBO"
}

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

			if account.Buy(df.Candles[i].Close, buySize) {
				signalEvents.Buy(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, buySize, false)
				isHolding = true
			}
		}
		if close[i] < donchain.Low[i-1] && isHolding {
			if account.Sell(df.Candles[i].Close) {
				signalEvents.Sell(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, buySize, false)
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

	for period := 10; period < 333; period++ {
		wg.Add(1)
		go func(period int) {
			defer wg.Done()
			account := trader.NewAccount(1000)
			signalEvents := df.DonchainStrategy(period, account)
			if signalEvents == nil {
				return
			}
			pf := analytics.PayOffRatio(signalEvents)
			mu.Lock()
			if performance < pf {
				performance = pf
				bestPeriod = period
			}
			mu.Unlock()
		}(period)
	}

	wg.Wait()

	fmt.Println("ペイオフレシオ", performance, "最適なピリオド", bestPeriod)

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

	strategyName := getStrageyNameDonchain()
	assetName := btcfg.AssetName
	duration := btcfg.Dration

	account := trader.NewAccount(1000)

	df, _ := GetCandleData(assetName, duration)

	tableName := strategyName + "_" + assetName + "_" + duration

	_, err = execute.CreateDBTable(tableName)
	if err != nil {
		log.Fatal(err)
	}

	// performanceProfit, bestProfit := df.OptimizeDonchainProfit()

	// if performanceProfit > 0 {

	// 	df.Signal = df.DonchainStrategy(bestProfit, account)
	// 	fmt.Println("🔺利益最適化")
	// 	Result(df.Signal)

	// }

	// performanceLoss, bestLoss := df.OptimizeDonchainLoss()

	// if performanceLoss > 0 {

	// 	df.Signal = df.DonchainStrategy(bestLoss, account)
	// 	fmt.Println("🔺損失最適化")
	// 	Result(df.Signal)

	// }

	// performanceWinRate, bestWinRate := df.OptimizeDonchainWinRate()

	// if performanceWinRate > 0 {

	// 	df.Signal = df.DonchainStrategy(bestWinRate, account)
	// 	fmt.Println("🔺勝率最適化")
	// 	Result(df.Signal)

	// }
	// performanceProfitPeriod, bestProfitPeriod := df.OptimizeDonchainProfitFactor()

	// if performanceProfitPeriod > 0 {

	// 	df.Signal = df.DonchainStrategy(bestProfitPeriod, account)
	// 	fmt.Println("🔺プロフィットファクター最適化")
	// 	Result(df.Signal)

	// }
	// goperformancePayOffRatio, gobestPayOffRatioPeriod := df.OptimizeDonchainGoroutin()

	// if goperformancePayOffRatio > 0 {

	// 	df.Signal = df.DonchainStrategy(gobestPayOffRatioPeriod, account)
	// 	fmt.Println("🔺goroutinペイオフレシオ最適化")
	// 	Result(df.Signal)

	// }

	performancePayOffRatio, bestPayOffRatioPeriod := df.OptimizeDonchainPayOffRatio()

	if performancePayOffRatio > 0 {

		df.Signal = df.DonchainStrategy(bestPayOffRatioPeriod, account)
		fmt.Println("🔺ペイオフレシオ最適化")
		Result(df.Signal)

	}

}
