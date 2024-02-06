package strategey

import (
	"fmt"
	"runtime"
	"sync"
	"v1/pkg/analytics"
	"v1/pkg/execute"
	"v1/pkg/management/risk"

	"v1/pkg/trader"

	"github.com/markcheno/go-talib"
)

func (df *DataFrameCandle) RsiStrategy2(period int, buyThread float64, tpRatio float64, slRatio float64, account *trader.Account) *execute.SignalEvents {

	var StrategyName = "RSI_SL&TP_RATIO"
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

	atr := talib.Atr(h, l, c, 13)
	ema := talib.Ema(hl3, 200)

	index := risk.ChoppySlice(70, df.Closes(), df.Highs(), df.Lows())
	choppyEma := risk.ChoppyEma(index, 11)

	for i := 1; i < lenCandles; i++ {
		if values[i-1] == 0 || values[i-1] == 100 || i >= len(choppyEma) {
			continue
		}

		sl := atr[i] * slRatio
		tp := atr[i] * tpRatio

		if values[i-1] < buyThread && values[i] >= buyThread && choppyEma[i] > 40 && c[i] < ema[i] && !isBuyHolding {
			accountBalance := account.GetBalance()
			buySize = account.TradeSize(riskSize) / df.Candles[i].Close
			buyPrice = df.Candles[i].Close
			if account.Buy(df.Candles[i].Close, buySize) {

				signalEvents.Buy(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, buySize, accountBalance, false)
				isBuyHolding = true

			}
		}

		if (df.Candles[i].Close > buyPrice+tp || df.Candles[i].Close < buyPrice-sl) && isBuyHolding {
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
func (df *DataFrameCandle) OptimizeRsi2() (performance float64, bestPeriod int, bestBuyThread, bestTpRatio, bestSlRatio float64) {
	runtime.GOMAXPROCS(10)

	bestPeriod = 13
	bestBuyThread = 20.0
	bestTpRatio = 3
	bestSlRatio = 1

	limit := 3000
	slots := make(chan struct{}, limit)

	// a := trader.NewAccount(1000)
	// marketDefault, _ := BuyAndHoldingStrategy(a)

	var mu sync.Mutex
	var wg sync.WaitGroup

	for period := 8; period < 18; period++ {
		for buyThread := 25.0; buyThread > 8.0; buyThread -= 2 {
			for slRatio := 1.0; slRatio < 5.0; slRatio += 0.5 {
				for tpRatio := 2.0; tpRatio < 15.0; tpRatio += 1.0 {
					wg.Add(1)
					slots <- struct{}{}
					go func(period int, buyThread, tpRatio, slRatio float64) {
						defer wg.Done()
						account := trader.NewAccount(1000) // Move this line inside the goroutine
						signalEvents := df.RsiStrategy2(period, buyThread, tpRatio, slRatio, account)

						if signalEvents == nil {
							<-slots
							return
						}

						if analytics.TotalTrades(signalEvents) < 35 {
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

						p := analytics.SortinoRatio(signalEvents, 0.02)
						mu.Lock()
						if performance == 0 || performance < p {
							performance = p
							bestPeriod = period
							bestBuyThread = buyThread
							bestTpRatio = tpRatio
							bestSlRatio = slRatio
						}

						mu.Unlock()
						<-slots
					}(period, buyThread, tpRatio, slRatio)
				}
			}
		}
	}

	wg.Wait()

	fmt.Println("最高パフォーマンス", performance, "最適なピリオド", bestPeriod, "最適な買いライン", bestBuyThread, "最適なTPレシオ", bestTpRatio, "最適なSLレシオ", bestSlRatio)

	return performance, bestPeriod, bestBuyThread, bestTpRatio, bestSlRatio
}

func RunRsi2Optimize() {

	df, account, _ := RadyBacktest()

	performance, bestPeriod, bestBuyThread, bestTpRatio, bestSlRatio := df.OptimizeRsi2()

	if performance > 0 {

		df.Signal = df.RsiStrategy2(bestPeriod, bestBuyThread, bestTpRatio, bestSlRatio, account)
		Result(df.Signal)

	} else {
		fmt.Println("マイナス利益です")
	}

}
