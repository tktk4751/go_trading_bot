package strategey

import (
	"fmt"
	"log"
	"math"
	"runtime"
	"sync"
	"v1/pkg/analytics"
	"v1/pkg/config"
	"v1/pkg/execute"

	"v1/pkg/trader"

	"github.com/markcheno/go-talib"
)

func getStrageyNameRSI() string {
	return "RSI"
}

func (df *DataFrameCandle) RsiStrategy(period int, buyThread float64, sellThread float64, account *trader.Account) *execute.SignalEvents {

	var StrategyName = "RSI"
	lenCandles := len(df.Candles)
	if lenCandles <= period {
		return nil
	}

	signalEvents := execute.NewSignalEvents()
	close := df.Closes()

	values := talib.Rsi(close, period)

	buySize := 0.0
	buyPrice := 0.0
	slRatio := 0.9
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

		if values[i-1] > sellThread && values[i] <= sellThread || (df.Candles[i].Close <= buyPrice*slRatio) && isBuyHolding {
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
func (df *DataFrameCandle) OptimizeRsi() (performance float64, bestPeriod int, bestBuyThread, bestSellThread float64) {
	runtime.GOMAXPROCS(10)

	bestPeriod = 13
	bestBuyThread, bestSellThread = 20.0, 80.0

	limit := 1000
	slots := make(chan struct{}, limit)

	a := trader.NewAccount(1000)
	marketDefault, _ := BuyAndHoldingStrategy(a)

	var mu sync.Mutex
	var wg sync.WaitGroup

	for period := 2; period < 30; period++ {
		for buyThread := 25.0; buyThread > 9; buyThread -= 1 {
			for sellThread := 75.0; sellThread < 96; sellThread += 1 {
				wg.Add(1)
				slots <- struct{}{}
				go func(period int, buyThread, sellThread float64) {
					defer wg.Done()
					account := trader.NewAccount(1000) // Move this line inside the goroutine
					signalEvents := df.RsiStrategy(period, buyThread, sellThread, account)

					if signalEvents == nil {
						<-slots
						return
					}

					if analytics.TotalTrades(signalEvents) < 20 {
						<-slots
						return
					}

					if analytics.NetProfit(signalEvents) < marketDefault {
						return
					}

					// if analytics.WinRate(signalEvents) < 0.50 {
					// 	<-slots
					// 	return
					// }

					// if analytics.PayOffRatio(signalEvents) < 1 {
					// 	return
					// }

					p := analytics.SQN(signalEvents)
					mu.Lock()
					if performance == 0 || performance < p {
						performance = p
						bestPeriod = period
						bestBuyThread = buyThread
						bestSellThread = sellThread
					}

					mu.Unlock()
					<-slots
				}(period, buyThread, sellThread)
			}
		}
	}

	wg.Wait()

	fmt.Println("最高利益", performance, "最適なピリオド", bestPeriod, "最適な買いライン", bestBuyThread, "最適な売りライン", bestSellThread)

	return performance, bestPeriod, bestBuyThread, bestSellThread
}

func (df *DataFrameCandle) OptimizeRsiDrawDownGoroutin() (performance float64, bestPeriod int, bestBuyThread, bestSellThread float64) {
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

					if analytics.NetProfit(signalEvents) < marketDefault {
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

func RunBacktestRsi() {

	var err error

	// account := trader.NewAccount(1000)
	btcfg, err := config.Yaml()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Println(btcfg.AssetName)

	strategyName := getStrageyNameRSI()
	assetName := btcfg.AssetName
	duration := btcfg.Dration
	// limit := btcfg.Limit

	account := trader.NewAccount(1000)

	df, _ := GetCandleData(assetName, duration)

	tableName := strategyName + "_" + assetName + "_" + duration

	_, err = execute.CreateDBTable(tableName)
	if err != nil {
		log.Fatal(err)
	}

	performance, bestPeriod, bestBuyThread, bestSellThread := df.OptimizeRsi()

	if performance > 0 {

		df.Signal = df.RsiStrategy(bestPeriod, bestBuyThread, bestSellThread, account)
		Result(df.Signal)

	}

}
