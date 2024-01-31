package strategey

import (
	"fmt"
	"log"
	"math"
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

func (df *DataFrameCandle) OptimizeBbProfit() (performance float64, bestN int, bestK float64) {

	account := trader.NewAccount(1000)
	bestN = 20
	bestK = 2.0

	for n := 13; n < 200; n++ {
		for k := 2.0; k < 3.5; k += 0.1 {
			signalEvents := df.BbStrategy(n, k, account)
			if signalEvents == nil {
				continue
			}
			if analytics.TotalTrades(signalEvents) < 50 {
				continue
			}
			profit := analytics.NetProfit(signalEvents)
			if performance < profit {
				performance = profit
				bestN = n
				bestK = k
			}
		}
	}

	fmt.Println("最高利益", performance, "最適なピリオド", bestN, "最適な標準偏差", bestK)

	return performance, bestN, bestK
}

func (df *DataFrameCandle) OptimizeBbLoss() (performance float64, bestN int, bestK float64) {

	account := trader.NewAccount(1000)

	bestN = 20
	bestK = 2.0
	performance = math.MaxFloat64

	for n := 5; n < 120; n++ {
		for k := 1.8; k < 3.8; k += 0.1 {
			signalEvents := df.BbStrategy(n, k, account)
			if signalEvents == nil {
				continue
			}
			loss := analytics.Loss(signalEvents)
			if performance < loss {
				performance = loss
				bestN = n
				bestK = k
			}
		}
	}

	fmt.Println("損失", performance, "最適なピリオド", bestN, "最適な標準偏差", bestK)

	return performance, bestN, bestK
}

func (df *DataFrameCandle) OptimizeBbWinRate() (performance float64, bestN int, bestK float64) {
	account := trader.NewAccount(1000)

	bestN = 20
	bestK = 2.0

	for n := 13; n < 200; n++ {
		for k := 1.8; k < 3.5; k += 0.1 {
			signalEvents := df.BbStrategy(n, k, account)
			if signalEvents == nil {
				continue
			}

			if analytics.TotalTrades(signalEvents) < 20 {
				continue
			}
			winrate := analytics.WinRate(signalEvents)
			if performance < winrate {
				performance = winrate
				bestN = n
				bestK = k
			}
		}
	}

	fmt.Println("最高勝率", performance, "最適なピリオド", bestN, "最適な標準偏差", bestK)

	return performance, bestN, bestK
}

func (df *DataFrameCandle) OptimizeBbProfitFactor() (performance float64, bestN int, bestK float64) {
	account := trader.NewAccount(1000)

	bestN = 20
	bestK = 2.0

	for n := 13; n < 200; n++ {
		for k := 1.8; k < 3.5; k += 0.1 {
			signalEvents := df.BbStrategy(n, k, account)
			if signalEvents == nil {
				continue
			}

			if analytics.TotalTrades(signalEvents) < 20 {
				continue
			}
			winrate := analytics.ProfitFactor(signalEvents)
			if performance < winrate {
				performance = winrate
				bestN = n
				bestK = k
			}
		}
	}

	fmt.Println("プロフィットファクター", performance, "最適なピリオド", bestN, "最適な標準偏差", bestK)

	return performance, bestN, bestK
}

func (df *DataFrameCandle) OptimizeBbPayOffRatio() (performance float64, bestN int, bestK float64) {
	account := trader.NewAccount(1000)

	bestN = 20
	bestK = 2.0

	for n := 13; n < 200; n++ {
		for k := 1.8; k < 3.5; k += 0.1 {
			signalEvents := df.BbStrategy(n, k, account)
			if signalEvents == nil {
				continue
			}

			if analytics.TotalTrades(signalEvents) < 20 {
				continue
			}
			winrate := analytics.PayOffRatio(signalEvents)
			if performance < winrate {
				performance = winrate
				bestN = n
				bestK = k
			}
		}
	}

	fmt.Println("ペイオフレシオ", performance, "最適なピリオド", bestN, "最適な標準偏差", bestK)

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

	strategyName := getStrageyNameBb()
	assetName := btcfg.AssetName
	duration := btcfg.Dration

	account := trader.NewAccount(1000)

	df, _ := GetCandleData(assetName, duration)

	tableName := strategyName + "_" + assetName + "_" + duration

	_, err = execute.CreateDBTable(tableName)
	if err != nil {
		log.Fatal(err)
	}

	performance, bestN, bestK := df.OptimizeBbProfitFactor()

	if performance > 0 {

		df.Signal = df.BbStrategy(bestN, bestK, account)
		Result(df.Signal)

	}

}
