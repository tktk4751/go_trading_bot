package strategey

import (
	"fmt"
	"runtime"
	"sync"
	"v1/pkg/analytics"
	"v1/pkg/execute"
	"v1/pkg/indicator/indicators"
	"v1/pkg/management/risk"

	"v1/pkg/trader"

	"github.com/markcheno/go-talib"
)

func (df *DataFrameCandle) BetterRsiStrategy(period int, buyThread float64, dcPeriod int, account *trader.Account, simlpe bool) *execute.SignalEvents {

	var StrategyName = "RSI_BETTER"
	lenCandles := len(df.Candles)
	if lenCandles <= period {
		return nil
	}

	signalEvents := execute.NewSignalEvents()

	hl3 := df.Hlc3()

	h := df.Highs()
	l := df.Lows()
	c := df.Closes()

	values := talib.Rsi(hl3, period)

	buySize := 0.0
	buyPrice := 0.0
	isBuyHolding := false
	slRatio := 4.0

	atr := talib.Atr(h, l, c, 13)
	// ema := talib.Ema(hl3, 200)

	donchain := indicators.Donchain(h, l, dcPeriod)

	index := risk.ChoppySlice(30, df.Closes(), df.Highs(), df.Lows())
	choppyEma := risk.ChoppyEma(index, 10)

	for i := 1; i < lenCandles; i++ {
		if values[i-1] == 0 || values[i-1] == 100 || i >= len(choppyEma) {
			continue
		}

		sl := atr[i] * slRatio
		// tp := atr[i] * tpRatio

		if (values[i-1] < buyThread && values[i] <= buyThread && c[i-1] < donchain.Low[i-2] && c[i] <= donchain.Low[i-1]) && choppyEma[i] > 60 && !isBuyHolding {
			// fee := 1 - 0.01
			if simple {
				buySize = account.SimpleTradeSize(1)
				buyPrice = c[i]
				accountBalance := account.GetBalance()

				signalEvents.Buy(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, buySize, accountBalance, false)
				isBuyHolding = true

			} else {
				buySize = account.TradeSize(riskSize) / df.Candles[i].Close
				buyPrice = c[i]
				accountBalance := account.GetBalance()
				if account.Buy(df.Candles[i].Close, buySize) {
					signalEvents.Buy(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, buySize, accountBalance, false)
					isBuyHolding = true
				}
			}

		}

		if ((c[i-1] > donchain.Mid[i-1] && c[i] <= donchain.Mid[i]) || (values[i] < buyThread) || choppyEma[i] < 10 || df.Candles[i].Close < buyPrice-sl) && isBuyHolding {
			if simple {
				accountBalance := 1000.0

				signalEvents.Sell(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, buySize, accountBalance, false)
				isBuyHolding = false

			} else {
				accountBalance := account.GetBalance()
				if account.Sell(df.Candles[i].Close) {
					signalEvents.Sell(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, buySize, accountBalance, false)
					isBuyHolding = false
					buySize = 0.0
					account.PositionSize = buySize

				}
			}
		}

	}

	return signalEvents
}
func (df *DataFrameCandle) OptimizeBetterRsi() (performance float64, bestPeriod int, bestBuyThread float64, bestDcPeriod int) {
	runtime.GOMAXPROCS(10)

	bestPeriod = 13
	bestBuyThread = 20.0
	bestDcPeriod = 40

	limit := 3000
	slots := make(chan struct{}, limit)

	// a := trader.NewAccount(1000)
	// marketDefault, _ := BuyAndHoldingStrategy(a)

	var mu sync.Mutex
	var wg sync.WaitGroup

	for period := 5; period < 16; period++ {
		for buyThread := 20.0; buyThread > 10.0; buyThread -= 1 {
			for dc := 200; dc < 450; dc += 10 {

				wg.Add(1)
				slots <- struct{}{}
				go func(period int, buyThread float64, dc int) {
					defer wg.Done()
					account := trader.NewAccount(1000) // Move this line inside the goroutine
					signalEvents := df.BetterRsiStrategy(period, buyThread, dc, account, simple)

					if signalEvents == nil {
						<-slots
						return
					}

					if analytics.TotalTrades(signalEvents) < 10 {
						<-slots
						return
					}

					// if analytics.NetProfit(signalEvents) < marketDefault {
					// 	return
					// }

					// if analytics.PayOffRatio(signalEvents) < 1.50 {
					// 	<-slots
					// 	return
					// }

					// if analytics.SQN(signalEvents) < 3.2 {
					// 	<-slots
					// 	return
					// }

					// p := analytics.ExpectedValue(signalEvents)
					p := analytics.SortinoRatio(signalEvents, 0.02)
					mu.Lock()
					if performance == 0 || performance < p {
						performance = p
						bestPeriod = period
						bestBuyThread = buyThread
						// bestTpRatio = tpRatio
						bestDcPeriod = dc
					}

					mu.Unlock()
					<-slots
				}(period, buyThread, dc)

			}
		}
	}

	wg.Wait()

	fmt.Println("最高パフォーマンス", performance, "最適なピリオド", bestPeriod, "最適な買いライン", bestBuyThread, "最適なドンチャン ", bestDcPeriod)

	return performance, bestPeriod, bestBuyThread, bestDcPeriod
}

func RunBetterRsiOptimize() {

	df, account, _ := RadyBacktest()

	performance, bestPeriod, bestBuyThread, bestDcPeriod := df.OptimizeBetterRsi()

	if performance > 0 {

		df.Signal = df.BetterRsiStrategy(bestPeriod, bestBuyThread, bestDcPeriod, account, simple)
		Result(df.Signal)

	} else {
		fmt.Println("💸マイナスです")
		df.Signal = df.BetterRsiStrategy(bestPeriod, bestBuyThread, bestDcPeriod, account, simple)
		Result(df.Signal)
	}

}

func RSIBryyrtBacktest() {

	df, account, _ := RadyBacktest()

	df.Signal = df.BetterRsiStrategy(14, 20, 200, account, simple)
	Result(df.Signal)

}
