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

func (df *DataFrameCandle) STDonchain(atrPeriod int, factor float64, dcPeriod int, duration int, account *trader.Account) *execute.SignalEvents {

	var StrategyName = "STDONCHAIN"
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

	index := risk.ChoppySlice(duration, c, h, l)
	choppyEma := risk.ChoppyEma(index, 11)

	donchain := indicators.Donchain(h, l, dcPeriod)
	ema := talib.Ema(df.Hlc3(), 89)

	isBuyHolding := false

	for i := 1; i < len(choppyEma); i++ {

		if i < atrPeriod {
			// fmt.Printf("Skipping iteration %d due to insufficient data.\n", i)
			continue
		}
		if c[i-1] < st[i-1] && c[i] >= st[i] && choppyEma[i] > 50 && c[i] > ema[i] && !isBuyHolding {

			accountBalance := account.GetBalance()
			buySize = account.TradeSize(riskSize) / c[i]
			buyPrice = c[i]
			if account.Buy(c[i], buySize) {

				signalEvents.Buy(StrategyName, df.AssetName, df.Duration, t[i], c[i], buySize, accountBalance, false)
				isBuyHolding = true

			}
		}

		if c[i] > donchain.High[i-1] && choppyEma[i] > 50 && c[i] > ema[i] && !isBuyHolding {

			accountBalance := account.GetBalance()
			buySize = account.TradeSize(riskSize) / c[i]
			buyPrice = c[i]
			if account.Buy(c[i], buySize) {

				signalEvents.Buy(StrategyName, df.AssetName, df.Duration, t[i], c[i], buySize, accountBalance, false)
				isBuyHolding = true

			}
		}

		if (c[i-1] > st[i-1] && c[i] <= st[i] || c[i] < donchain.Low[i-1] || (c[i] <= buyPrice*slRatio)) && isBuyHolding {
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

func (df *DataFrameCandle) OptimizeSTDonchain() (performance float64, bestAtrPeriod int, bestFactor float64, bestDc int, bestDuration int) {
	runtime.GOMAXPROCS(10)
	bestAtrPeriod = 21
	bestFactor = 3.0
	bestDc = 40
	bestDuration = 30

	limit := 1000
	slots := make(chan struct{}, limit)

	// var accountPool = sync.Pool{
	// 	New: func() interface{} {
	// 		return trader.NewAccount(1000)
	// 	},
	// }

	// a := trader.NewAccount(1000)
	// marketDefault, _ := BuyAndHoldingStrategy(a)

	var mu sync.Mutex
	var wg sync.WaitGroup

	for atrPeriod := 5; atrPeriod < 40; atrPeriod += 2 {
		for factor := 2.0; factor < 8.0; factor += 0.5 {
			for dc := 10; dc < 60; dc += 10 {
				for duration := 30; duration < 150; duration += 10 {

					wg.Add(1)
					slots <- struct{}{}

					go func(atrPeriod int, factor float64, dc int, duration int) {
						defer wg.Done()
						// account := accountPool.Get().(*trader.Account)
						// defer accountPool.Put(account)
						account := trader.NewAccount(1000)
						signalEvents := df.STDonchain(atrPeriod, factor, dc, duration, account)

						if analytics.TotalTrades(signalEvents) < 30 {
							<-slots
							return
						}

						// if analytics.NetProfit(signalEvents) < marketDefault {
						// 	<-slots
						// 	return
						// }

						// if analytics.SQN(signalEvents) < 3.2 {
						// 	<-slots
						// 	return
						// }

						// if analytics.PayOffRatio(signalEvents) < 1 {
						// <-slots

						// 	return
						// }

						// pf := analytics.SortinoRatio(signalEvents, 0.02)
						p := analytics.SQN(signalEvents)
						mu.Lock()
						if performance == 0 || performance < p {
							performance = p
							bestAtrPeriod = atrPeriod
							bestFactor = factor
							bestDc = dc
							bestDuration = duration

						}
						<-slots
						mu.Unlock()

					}(atrPeriod, factor, dc, duration)

				}
			}
		}
	}

	wg.Wait()

	fmt.Println("最高のパフォーマンス", performance, "最適なATR", bestAtrPeriod, "最適なファクター", bestFactor, "最適なドンチャン", bestDc, "最適なチョッピー期間", bestDuration)

	return performance, bestAtrPeriod, bestFactor, bestDc, bestDuration
}

func RunSTDonchainOptimize() {

	df, account, _ := RadyBacktest()

	performance, bestAtrPeriod, bestFactor, bestDc, bestDuration := df.OptimizeSTDonchain()

	if performance > 0 {

		df.Signal = df.STDonchain(bestAtrPeriod, bestFactor, bestDc, bestDuration, account)
		Result(df.Signal)

	}

}

func STDonchainBacktest() {

	df, account, _ := RadyBacktest()

	df.Signal = df.STDonchain(13, 3.0, 60, 30, account)
	Result(df.Signal)

}
