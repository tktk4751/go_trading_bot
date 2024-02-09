package strategey

import (
	"fmt"
	"runtime"
	"sync"
	"v1/pkg/analytics"
	"v1/pkg/execute"
	"v1/pkg/management/risk"
	"v1/pkg/trader"

	"github.com/c-bata/goptuna"

	"github.com/c-bata/goptuna/tpe"

	"github.com/markcheno/go-talib"
)

func (df *DataFrameCandle) EmaChoppyStrategy(period1, period2 int, choppy int, duration int, account *trader.Account, simple bool) *execute.SignalEvents {

	var StrategyName = "EMA_CHOPPY"
	lenCandles := len(df.Candles)
	if lenCandles <= period1 || lenCandles <= period2 {
		return nil
	}
	signalEvents := execute.NewSignalEvents()

	c := df.Closes()

	emaValue1 := talib.Ema(df.Hlc3(), period1)
	emaValue2 := talib.Ema(df.Hlc3(), period2)
	// rsiValue := talib.Rsi(df.Closes(), 14)

	buySize := 0.0
	buyPrice := 0.0
	slRatio := 0.9
	index := risk.ChoppySlice(duration, df.Closes(), df.Highs(), df.Lows())
	choppyEma := risk.ChoppyEma(index, choppy)

	isBuyHolding := false
	for i := 1; i < lenCandles; i++ {
		if i < period1 || i < period2 || i >= len(choppyEma) {
			continue
		}

		if emaValue1[i-1] < emaValue2[i-1] && emaValue1[i] >= emaValue2[i] && choppyEma[i] > 50 && !isBuyHolding {

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
		if emaValue1[i-1] > emaValue2[i-1] && emaValue1[i] <= emaValue2[i] || (df.Candles[i].Close <= buyPrice*slRatio) && isBuyHolding {
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

func (df *DataFrameCandle) OptimizeEmaChoppy() (performance float64, bestPeriod1 int, bestPeriod2 int, bestChoppy int, bestDuration int) {

	// オブジェクティブ関数を定義
	objective := func(trial goptuna.Trial) (float64, error) {
		// ハイパーパラメータの候補をサンプリング
		period1, _ := trial.SuggestInt("period1", 3, 69)
		period2, _ := trial.SuggestInt("period2", 13, 200)
		choppy, _ := trial.SuggestInt("choppy", 5, 18)
		duration, _ := trial.SuggestInt("duration", 10, 200)

		account := trader.NewAccount(1000) // Move this line inside the objective function
		signalEvents := df.EmaChoppyStrategy(period1, period2, choppy, duration, account, simple)

		if signalEvents == nil {
			return 0.0, nil
		}

		if analytics.TotalTrades(signalEvents) < 10 {
			return 0.0, nil
		}

		// p := analytics.SortinoRatio(signalEvents, 0.02)
		p := analytics.Prr(signalEvents)
		return p, nil // パフォーマンスを返す
	}

	// ベイズ最適化の設定
	study, err := goptuna.CreateStudy(
		"ema-choppy-optimization",
		goptuna.StudyOptionSampler(tpe.NewSampler()),                 // 獲得関数としてTPEを使用
		goptuna.StudyOptionDirection(goptuna.StudyDirectionMaximize), // 最大化問題として定義
		goptuna.StudyOptionLogger(nil),
	)
	if err != nil {
		panic(err)
	}

	// ベイズ最適化の実行
	err = study.Optimize(objective, 500)
	if err != nil {
		panic(err)
	}

	// 最適化結果の取得
	v, _ := study.GetBestValue()
	params, _ := study.GetBestParams()
	performance = v
	bestPeriod1 = params["period1"].(int)
	bestPeriod2 = params["period2"].(int)
	bestChoppy = params["choppy"].(int)
	bestDuration = params["duration"].(int)

	fmt.Println("最高パフォーマンス", performance, "最適な短期線", bestPeriod1, "最適な長期線", bestPeriod2, "最適なチョッピーEMA", bestChoppy, "最適なチョッピー期間", bestDuration)

	return performance, bestPeriod1, bestPeriod2, bestChoppy, bestDuration
}

func (df *DataFrameCandle) OptimizeEmaChoppy2() (performance float64, bestPeriod1 int, bestPeriod2 int, bestChoppy int, bestDuration int) {
	runtime.GOMAXPROCS(10)
	bestPeriod1 = 5
	bestPeriod2 = 21
	bestDuration = 30

	limit := 3000
	slots := make(chan struct{}, limit)

	// a := trader.NewAccount(1000)
	// marketDefault, _ := BuyAndHoldingStrategy(a)

	var mu sync.Mutex
	var wg sync.WaitGroup

	for period1 := 3; period1 < 34; period1 += 2 {
		for period2 := 13; period2 < 89; period2 += 4 {
			for choppy := 5; choppy < 18; choppy += 1 {
				for duration := 10; duration < 200; duration += 10 {

					wg.Add(1)
					slots <- struct{}{}

					go func(period1 int, period2 int, choppy int, duration int) {
						defer wg.Done()
						account := trader.NewAccount(1000) // Move this line inside the goroutine
						signalEvents := df.EmaChoppyStrategy(period1, period2, choppy, duration, account, simple)

						if signalEvents == nil {
							return
						}

						if analytics.TotalTrades(signalEvents) < 10 {
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

						// p := analytics.SortinoRatio(signalEvents, 0.02)
						p := analytics.Prr(signalEvents)
						mu.Lock()
						if performance == 0 || performance < p {
							performance = p
							bestPeriod1 = period1
							bestPeriod2 = period2
							bestChoppy = choppy
							bestDuration = duration

						}
						<-slots
						mu.Unlock()

					}(period1, period2, 13, duration)

				}
			}
		}
	}

	wg.Wait()

	fmt.Println("最高パフォーマンス", performance, "最適な短期線", bestPeriod1, "最適な長期線", bestPeriod2, "最適なチョッピーEMA", bestChoppy, "最適なチョッピー期間", bestDuration)

	return performance, bestPeriod1, bestPeriod2, bestChoppy, bestDuration
}

func RunEmaOptimize() {

	df, account, _ := RadyBacktest()

	performance, bestPeriod1, bestPeriod2, bestChoppy, bestDuration := df.OptimizeEmaChoppy()

	if performance > 0 {

		df.Signal = df.EmaChoppyStrategy(bestPeriod1, bestPeriod2, bestChoppy, bestDuration, account, simple)
		Result(df.Signal)

	} else {
		fmt.Println("💸マイナスです")
		df.Signal = df.EmaChoppyStrategy(bestPeriod1, bestPeriod2, bestChoppy, bestDuration, account, simple)
		Result(df.Signal)

	}

}

func RunEmaOptimize2() {

	df, account, _ := RadyBacktest()

	performance, bestPeriod1, bestPeriod2, bestChoppy, bestDuration := df.OptimizeEmaChoppy2()

	if performance > 0 {

		df.Signal = df.EmaChoppyStrategy(bestPeriod1, bestPeriod2, bestChoppy, bestDuration, account, simple)
		Result(df.Signal)

	} else {
		fmt.Println("💸マイナスです")
		df.Signal = df.EmaChoppyStrategy(bestPeriod1, bestPeriod2, bestChoppy, bestDuration, account, simple)
		Result(df.Signal)

	}

}

func EmaBacktest() {

	df, account, _ := RadyBacktest()

	df.Signal = df.EmaChoppyStrategy(13, 33, 13, 30, account, simple)
	Result(df.Signal)

}
