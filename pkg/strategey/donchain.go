package strategey

import (
	"fmt"
	"log"
	"v1/pkg/analytics"
	"v1/pkg/execute"
	"v1/pkg/indicator/indicators"
)

var AccountBalance = NewAccount(1000.00)

func GetStrageyName() string {
	return "DBO"
}

func (df *DataFrameCandle) DonchainStrategy(period int, account *Account) *execute.SignalEvents {
	var StrategyName = "DBO"

	// fmt.Println("アカウントバランス", account.Balance)
	lenCandles := len(df.Candles)
	if lenCandles <= period {
		return nil
	}

	signalEvents := execute.NewSignalEvents()

	donchain := indicators.Donchain(df.Highs(), df.Low(), period)

	close := df.Closes()

	buySize := 0.0
	isHolding := false

	for i := 1; i < lenCandles; i++ {

		if i < period {
			continue
		}
		if close[i] > donchain.High[i-1] && !isHolding {
			buySize = account.TradeSize(0.2) / df.Candles[i].Close
			if account.Buy(df.Candles[i].Close, buySize) {
				signalEvents.Buy(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, buySize, true)
				isHolding = true
			}
		}
		if close[i] < donchain.Low[i-1] && isHolding {
			if account.Sell(df.Candles[i].Close) {
				signalEvents.Sell(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, buySize, true)
				isHolding = false

			}
		}

	}
	return signalEvents

}

func (df *DataFrameCandle) OptimizeDonchainProfit() (performance float64, bestPeriod int) {
	if df == nil {
		return 0.0, 0
	}
	bestPeriod = 40

	for period := 10; period < 333; period++ {

		signalEvents := df.DonchainStrategy(period, AccountBalance)
		if signalEvents == nil {
			continue
		}
		profit := Profit(signalEvents)
		if performance < profit {
			performance = profit
			bestPeriod = period

		}

	}

	fmt.Println("最高利益", performance, "最適なピリオド", bestPeriod)

	return performance, bestPeriod
}

func (df *DataFrameCandle) OptimizeDonchainWinRate() (performance float64, bestPeriod int) {
	bestPeriod = 40

	for period := 10; period < 333; period++ {

		signalEvents := df.DonchainStrategy(period, AccountBalance)
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

func RunBacktestDonchain() {

	strategyName := GetStrageyName()
	assetName := "OPUSDT"
	duration := "1h"

	df, _ := GetCandleData(assetName, duration)

	tableName := strategyName + "_" + assetName + "_" + duration

	_, err := execute.CreateDBTable(tableName)
	if err != nil {
		log.Fatal(err)
	}

	// df, _ := strategey.GetCandleData(assetName, duration)

	// profit, period := df.OptimizeProfitDonchain()

	// if profit > 0 {

	// 	df.Signal = df.DonchainStrategy(period)

	// }

	// winrate, period := df.OptimizeWinRateDonchain()

	// if winrate > 0 {

	// 	df.Signal = df.DonchainStrategy(period)

	// }

	performance, bestPeriod := df.OptimizeDonchainProfit()

	if performance > 0 {

		df.Signal = df.DonchainStrategy(bestPeriod, AccountBalance)

	}

	Result(df.Signal)

}
