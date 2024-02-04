package strategey

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"v1/pkg/analytics"
	"v1/pkg/config"
	"v1/pkg/execute"
	"v1/pkg/indicator/indicators"
	"v1/pkg/trader"
)

func (df *DataFrameCandleCsv) SuperTrend(atrPeriod int, factor float64, account *trader.Account) *execute.SignalEvents {

	var StrategyName = "SUPERTREND_CHOPPY"
	// var err error

	lenCandles := len(df.Candles)

	if lenCandles <= atrPeriod {
		return nil
	}

	signalEvents := execute.NewSignalEvents()
	t := df.Time()
	h := df.Highs()
	l := df.Lows()
	c := df.Closes()

	superTrend, _ := indicators.SuperTrend(atrPeriod, factor, h, l, c)

	up := superTrend.UpperBand
	low := superTrend.UpperBand
	// st := superTrend.SuperTrend

	// rsiValue := talib.Rsi(df.Closes(), 14)

	buySize := 0.0
	buyPrice := 0.0
	slRatio := 0.9

	isBuyHolding := false

	for i := 1; i < lenCandles; i++ {

		if i < atrPeriod {
			// fmt.Printf("Skipping iteration %d due to insufficient data.\n", i)
			continue
		}
		if c[i-1] < up[i-1] && c[i] >= up[i] && !isBuyHolding {

			accountBalance := account.GetBalance()
			buySize = account.TradeSize(riskSize) / c[i]
			buyPrice = c[i]
			if account.Buy(c[i], buySize) {

				signalEvents.Buy(StrategyName, df.AssetName, df.Duration, t[i], c[i], buySize, accountBalance, false)
				isBuyHolding = true

			}
		}
		if c[i-1] > low[i-1] && c[i] <= low[i] || (c[i] <= buyPrice*slRatio) && isBuyHolding {
			accountBalance := account.GetBalance()
			if account.Sell(c[i]) {
				signalEvents.Sell(StrategyName, df.AssetName, df.Duration, t[i], c[i], buySize, accountBalance, false)
				isBuyHolding = false
				buySize = 0.0
				account.PositionSize = buySize

			}
		}
	}

	// fmt.Println(signalEvents)
	return signalEvents
}

func (df *DataFrameCandleCsv) OptimizeST() (performance float64, bestAtrPeriod int, bestFactor float64) {
	runtime.GOMAXPROCS(10)
	bestAtrPeriod = 12
	bestFactor = 21

	limit := 1000
	slots := make(chan struct{}, limit)

	// a := trader.NewAccount(1000)
	// marketDefault, _ := BuyAndHoldingStrategy(a)

	var mu sync.Mutex
	var wg sync.WaitGroup

	for atrPeriod := 9; atrPeriod < 40; atrPeriod += 1 {
		for factor := 2.0; factor < 8.0; factor += 0.2 {

			wg.Add(1)
			slots <- struct{}{}

			go func(atrPeriod int, factor float64) {
				defer wg.Done()
				account := trader.NewAccount(1000) // Move this line inside the goroutine
				signalEvents := df.SuperTrend(atrPeriod, factor, account)

				if signalEvents == nil {
					return
				}

				if analytics.TotalTrades(signalEvents) < 3 {
					<-slots
					return
				}

				// if analytics.NetProfit(signalEvents) < marketDefault {
				// 	<-slots
				// 	return
				// }

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
					bestAtrPeriod = atrPeriod
					bestFactor = factor
				}
				<-slots
				mu.Unlock()

			}(atrPeriod, factor)

		}
	}

	wg.Wait()

	fmt.Println("最高のプロフィットファクター", performance, "最適なATR", bestAtrPeriod, "最適なファクター", bestFactor)

	return performance, bestAtrPeriod, bestFactor
}

func RunBacktestST() {

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

	df, _ := GetCsvDataFrame(assetName, duration, "2023-01", "2023-13")

	performance, bestAtrPeriod, bestFactor := df.OptimizeST()

	if performance > 0 {

		df.Signal = df.SuperTrend(bestAtrPeriod, bestFactor, account)
		Result(df.Signal)

	}

}
