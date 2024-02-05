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
	"v1/pkg/management/risk"
	"v1/pkg/trader"
)

func (df *DataFrameCandle) SuperTrendChoppyStrategy(atrPeriod int, factor float64, choppy int, account *trader.Account) *execute.SignalEvents {

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

	// stUp := superTrend.UpperBand
	// stLow := superTrend.UpperBand
	st := superTrend.SuperTrend

	// rsiValue := talib.Rsi(df.Closes(), 14)

	buySize := 0.0
	buyPrice := 0.0
	slRatio := 0.9

	index := risk.ChoppySlice(c, h, l)
	choppyEma := risk.ChoppyEma(index, choppy)

	isBuyHolding := false

	for i := 1; i < len(choppyEma); i++ {

		if i < atrPeriod {
			// fmt.Printf("Skipping iteration %d due to insufficient data.\n", i)
			continue
		}
		if c[i-1] < st[i-1] && c[i] >= st[i] && choppyEma[i] > 50 && !isBuyHolding {

			accountBalance := account.GetBalance()
			buySize = account.TradeSize(riskSize) / c[i]
			buyPrice = c[i]
			if account.Buy(c[i], buySize) {

				signalEvents.Buy(StrategyName, df.AssetName, df.Duration, t[i], c[i], buySize, accountBalance, false)
				isBuyHolding = true

			}
		}
		if c[i-1] > st[i-1] && c[i] <= st[i] || (c[i] <= buyPrice*slRatio) && isBuyHolding {
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

func (df *DataFrameCandle) OptimizeSuperTrend() (performance float64, bestAtrPeriod int, bestFactor float64, bestChoppy int) {
	runtime.GOMAXPROCS(10)
	bestAtrPeriod = 21
	bestFactor = 3.0
	bestChoppy = 13

	limit := 1000
	slots := make(chan struct{}, limit)

	// a := trader.NewAccount(1000)
	// marketDefault, _ := BuyAndHoldingStrategy(a)

	var mu sync.Mutex
	var wg sync.WaitGroup

	for atrPeriod := 9; atrPeriod < 40; atrPeriod += 1 {
		for factor := 2.0; factor < 8.0; factor += 0.2 {
			for choppy := 5; choppy < 18; choppy += 1 {

				wg.Add(1)
				slots <- struct{}{}

				go func(atrPeriod int, factor float64, choppy int) {
					defer wg.Done()
					account := trader.NewAccount(1000) // Move this line inside the goroutine
					signalEvents := df.SuperTrendChoppyStrategy(atrPeriod, factor, choppy, account)

					if signalEvents == nil {
						return
					}

					if analytics.TotalTrades(signalEvents) < 5 {
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

					p := analytics.SQN(signalEvents)
					mu.Lock()
					if performance == 0 || performance < p {
						performance = p
						bestAtrPeriod = atrPeriod
						bestFactor = factor
						bestChoppy = choppy

					}
					<-slots
					mu.Unlock()

				}(atrPeriod, factor, choppy)

			}
		}
	}

	wg.Wait()

	fmt.Println("最高のSQN", performance, "最適なATR", bestAtrPeriod, "最適なファクター", bestFactor, "最適なチョッピー", bestChoppy)

	return performance, bestAtrPeriod, bestFactor, bestChoppy
}

func RunBacktestSuperTrend() {

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

	performance, bestAtrPeriod, bestFactor, bestChoppy := df.OptimizeSuperTrend()

	if performance > 0 {

		df.Signal = df.SuperTrendChoppyStrategy(bestAtrPeriod, bestFactor, bestChoppy, account)
		Result(df.Signal)

	}

}
