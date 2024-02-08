package strategey

import (
	"fmt"
	"math"
	"runtime"
	"sync"
	"v1/pkg/analytics"
	"v1/pkg/execute"
	"v1/pkg/indicator/indicators"

	"v1/pkg/trader"

	"github.com/markcheno/go-talib"
)

func (df *DataFrameCandle) RsiDonchainStrategy(rsiPeriod int, donchainPeriod int, buyThread float64, account *trader.Account) *execute.SignalEvents {

	var StrategyName = "RSI_DONCHAIN"
	lenCandles := len(df.Candles)
	if lenCandles <= donchainPeriod {
		return nil
	}

	signalEvents := execute.NewSignalEvents()
	close := df.Closes()

	values := talib.Rsi(close, rsiPeriod)

	donchain := indicators.Donchain(df.Highs(), df.Lows(), donchainPeriod)

	buySize := 0.0
	buyPrice := 0.0
	slRatio := 0.1

	isBuyHolding := false
	for i := 1; i < lenCandles; i++ {
		if values[i-1] == 0 || values[i-1] == 100 {
			continue
		}

		if values[i-1] < buyThread && values[i] >= buyThread && !isBuyHolding {
			accountBalance := account.GetBalance()
			buySize = account.TradeSize(riskSize) / df.Candles[i].Close
			buyPrice = df.Candles[i].Close
			if account.Buy(df.Candles[i].Close, buySize) {

				signalEvents.Buy(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, buySize, accountBalance, false)
				isBuyHolding = true

			}
		}

		if close[i] > donchain.High[i-1] || (df.Candles[i].Close < buyPrice-buyPrice*slRatio) || values[i-1] > 90 && isBuyHolding {
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

func (df *DataFrameCandle) OptimizeRsiDonchainGoroutin() (performance float64, bestRsiPeriod int, bestDonchainPeriod int, bestBuyThread float64) {
	runtime.GOMAXPROCS(12)

	bestRsiPeriod = 13
	bestDonchainPeriod = 20
	bestBuyThread = 20.0

	limit := 1000
	slots := make(chan struct{}, limit)

	// var pool = sync.Pool{
	// 	New: func() interface{} {
	// 		return trader.NewAccount(1000)
	// 	},
	// }

	// a := trader.NewAccount(1000)

	// marketDefault, _ := BuyAndHoldingStrategy(a)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for rsiPeriod := 2; rsiPeriod < 14; rsiPeriod++ {
		for buyThread := 30.0; buyThread > 8.0; buyThread -= 1 {
			for donchainPeriod := 20; donchainPeriod < 100; donchainPeriod += 5 {
				wg.Add(1)
				slots <- struct{}{}
				go func(rsiPeriod int, buyThread float64, donchainPeriod int) {
					defer wg.Done()
					account := trader.NewAccount(1000) // Move this line inside the goroutine
					// account := pool.Get().(*trader.Account)
					// defer pool.Put(account)
					signalEvents := df.RsiDonchainStrategy(rsiPeriod, donchainPeriod, buyThread, account)

					if signalEvents == nil {
						return
					}

					if analytics.TotalTrades(signalEvents) < 10 {
						return
					}

					// if analytics.NetProfit(signalEvents) < marketDefault {
					// 	return
					// }

					// if analytics.WinRate(signalEvents) < 0.50 {
					// 	return
					// }

					// if analytics.GainPainRatio(signalEvents) < 1 {
					// 	return
					// }

					p := analytics.LongNetProfit(signalEvents)
					mu.Lock()
					if performance == 0 || performance < p {
						performance = p
						bestRsiPeriod = rsiPeriod
						bestBuyThread = buyThread
						bestDonchainPeriod = donchainPeriod
					}
					<-slots
					mu.Unlock()
				}(rsiPeriod, buyThread, donchainPeriod)
			}
		}
	}

	wg.Wait()

	fmt.Println("最高利益", performance, "最適なRSI", bestRsiPeriod, "最適な買いライン", bestBuyThread, "最適なドンチャン", bestDonchainPeriod)

	return performance, bestRsiPeriod, bestDonchainPeriod, bestBuyThread
}

func (df *DataFrameCandle) OptimizeRsiDonchainDrawDownGoroutin() (performance float64, bestPeriod int, bestBuyThread, bestSellThread float64) {
	runtime.GOMAXPROCS(10)

	bestPeriod = 13
	bestBuyThread, bestSellThread = 20.0, 80.0

	performance = math.MaxFloat64

	a := trader.NewAccount(1000)

	marketDefault, _ := BuyAndHoldingStrategy(a)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for period := 2; period < 28; period++ {
		for buyThread := 30.0; buyThread > 10; buyThread -= 1 {
			for sellThread := 70.0; sellThread < 96; sellThread += 1 {
				wg.Add(1)
				go func(period int, buyThread, sellThread float64) {
					defer wg.Done()
					account := trader.NewAccount(1000) // Move this line inside the goroutine
					signalEvents := df.RsiStrategy(period, buyThread, sellThread, account)

					if signalEvents == nil {
						return
					}

					if analytics.TotalTrades(signalEvents) < 5 {
						return
					}

					if analytics.LongNetProfit(signalEvents) < marketDefault {
						return
					}

					dd := analytics.MaxDrawdownUSD(signalEvents)
					mu.Lock()
					if performance > dd {
						performance = dd
						bestPeriod = period
						bestBuyThread = buyThread
						bestSellThread = sellThread
					}
					mu.Unlock()
				}(period, buyThread, sellThread)
			}
		}
	}

	wg.Wait()

	fmt.Println("ドローダウン", performance, "最適なピリオド", bestPeriod, "最適な買いライン", bestBuyThread, "最適な売りライン", bestSellThread)

	return performance, bestPeriod, bestBuyThread, bestSellThread
}

func RunBacktestRsiDonchain() {

	df, account, _ := RadyBacktest()

	performance, bestRsiPeriod, bestDonchainPeriod, bestBuyThread := df.OptimizeRsiDonchainGoroutin()

	if performance > 0 {

		df.Signal = df.RsiDonchainStrategy(bestRsiPeriod, bestDonchainPeriod, bestBuyThread, account)
		Result(df.Signal)

	}

}
