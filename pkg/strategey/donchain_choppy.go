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

//	func getStrageyNameDonchain() string {
//		return "DBO"
//	}

func (df *DataFrameCandle) DonchainChoppyStrategy(period int, choppy int, duration int, account *trader.Account, simple bool) *execute.SignalEvents {
	var StrategyName = "DBO_CHOPPY"

	lenCandles := len(df.Candles)
	if lenCandles <= period {
		return nil
	}

	signalEvents := execute.NewSignalEvents()

	donchain := indicators.Donchain(df.Highs(), df.Lows(), period)
	// atr := talib.Atr(df.Highs(), df.Low(), df.Closes(), 21)

	ema := talib.Ema(df.Hlc3(), 89)

	close := df.Closes()

	buySize := 0.0
	buyPrice := 0.0
	slRatio := 0.9
	isHolding := false

	index := risk.ChoppySlice(duration, df.Closes(), df.Highs(), df.Lows())
	choppyEma := risk.ChoppyEma(index, choppy)

	for i := 30; i < lenCandles; i++ {

		if i < period || i >= len(choppyEma) {
			continue
		}

		if close[i] > donchain.High[i-1] && choppyEma[i] > 50 && close[i] > ema[i] && !isHolding {
			// fee := 1 - 0.01
			if simple {
				buySize = account.SimpleTradeSize(1)
				buyPrice = close[i]
				accountBalance := account.GetBalance()

				signalEvents.Buy(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, buySize, accountBalance, false)
				isHolding = true

			} else {
				buySize = account.TradeSize(riskSize) / df.Candles[i].Close
				buyPrice = close[i]
				accountBalance := account.GetBalance()
				if account.Buy(df.Candles[i].Close, buySize) {
					signalEvents.Buy(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, buySize, accountBalance, false)
					isHolding = true
				}
			}

		}
		if (close[i] < donchain.Low[i-1] || (close[i] <= buyPrice*slRatio)) && isHolding {

			if simple {
				accountBalance := 1000.0

				signalEvents.Sell(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, buySize, accountBalance, false)
				isHolding = false

			} else {
				accountBalance := account.GetBalance()
				if account.Sell(df.Candles[i].Close) {
					signalEvents.Sell(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, buySize, accountBalance, false)
					isHolding = false
					buySize = 0.0
					account.PositionSize = buySize

				}
			}

		}

	}
	return signalEvents

}

func (df *DataFrameCandle) OptimizeDonchain() (performance float64, bestPeriod int, bestChoppy int, bestDuration int) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	bestPeriod = 40
	bestChoppy = 13
	var mu sync.Mutex
	var wg sync.WaitGroup

	// a := trader.NewAccount(1000)
	// marketDefault, _ := BuyAndHoldingStrategy(a)

	limit := 3000
	slots := make(chan struct{}, limit)

	for period := 20; period < 300; period += 5 {
		for duration := 30; duration < 80; duration += 10 {
			for choppy := 3; choppy < 11; choppy += 2 {
				wg.Add(1)
				slots <- struct{}{}

				go func(period int, choppy int, duration int) {
					defer wg.Done()
					account := trader.NewAccount(1000)
					signalEvents := df.DonchainChoppyStrategy(period, choppy, duration, account, simple)

					if signalEvents == nil {
						return
					}

					// if analytics.TotalTrades(signalEvents) < 10 {
					// 	<-slots
					// 	return
					// }

					// if analytics.NetProfit(signalEvents) < marketDefault {
					// 	// <-slots
					// 	return
					// }

					// if analytics.SQN(signalEvents) < 3.2 {
					// 	<-slots
					// 	return
					// }

					// if analytics.ProfitFactor(signalEvents) < 3 {
					// <-slots
					// 	return
					// }

					pf := analytics.SortinoRatio(signalEvents, 0.02)
					// pf := analytics.Prr(signalEvents)
					mu.Lock()
					if performance < pf {
						performance = pf
						bestPeriod = period
						bestChoppy = choppy
						bestDuration = duration
					}
					<-slots
					mu.Unlock()
				}(period, choppy, duration)
			}
		}
	}

	wg.Wait()

	fmt.Println("最高パフォーマンス", performance, "最適なピリオド", bestPeriod, "最適なチョッピー", bestChoppy, "最適なチョッピー期間", bestDuration)

	return performance, bestPeriod, bestChoppy, bestDuration
}

// func (df *DataFrameCandle) OptimizeDonchain() (performance float64, bestPeriod int, bestChoppy int, bestDuration int) {

// 	// オブジェクティブ関数を定義
// 	objective := func(trial goptuna.Trial) (float64, error) {
// 		// ハイパーパラメータの候補をサンプリング
// 		period, _ := trial.SuggestStepInt("atrPeriod", 20, 300, 10)
// 		choppy, _ := trial.SuggestStepInt("choppy", 5, 18, 1)
// 		duration, _ := trial.SuggestStepInt("duration", 10, 200, 10)

// 		account := trader.NewAccount(1000) // Move this line inside the objective function
// 		signalEvents := df.DonchainChoppyStrategy(period, choppy, duration, account, simple)

// 		if signalEvents == nil {
// 			return 0.0, nil
// 		}

// 		if analytics.TotalTrades(signalEvents) < 10 {
// 			return 0.0, nil
// 		}

// 		// p := analytics.SortinoRatio(signalEvents, 0.02)
// 		p := analytics.Prr(signalEvents)
// 		return p, nil // パフォーマンスを返す
// 	}

// 	// ベイズ最適化の設定
// 	study, err := goptuna.CreateStudy(
// 		"donchain-choppy-optimization",
// 		goptuna.StudyOptionSampler(tpe.NewSampler()),                 // 獲得関数としてTPEを使用
// 		goptuna.StudyOptionDirection(goptuna.StudyDirectionMaximize), // 最大化問題として定義
// 		goptuna.StudyOptionLogger(nil),
// 	)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// ベイズ最適化の実行
// 	err = study.Optimize(objective, 800)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// 最適化結果の取得
// 	v, _ := study.GetBestValue()
// 	params, _ := study.GetBestParams()
// 	performance = v
// 	bestPeriod = params["atrPeriod"].(int)
// 	bestChoppy = params["choppy"].(int)
// 	bestDuration = params["duration"].(int)

// 	fmt.Println("最高パフォーマンス", performance, "最適なピリオド", bestPeriod, "最適なチョッピー", bestChoppy, "最適なチョッピー期間", bestDuration)

// 	return performance, bestPeriod, bestChoppy, bestDuration
// }

func RunDonchainOptimize() {

	df, account, _ := RadyBacktest()

	p, bestPeriod, bestChoppy, bestDuration := df.OptimizeDonchain()

	if p > 0 {

		df.Signal = df.DonchainChoppyStrategy(bestPeriod, bestChoppy, bestDuration, account, simple)
		Result(df.Signal)

	} else {
		fmt.Println("💸マイナスです")
		df.Signal = df.DonchainChoppyStrategy(bestPeriod, bestChoppy, bestDuration, account, simple)
		Result(df.Signal)

	}

}

func DonchainBacktest() {

	df, account, _ := RadyBacktest()

	df.Signal = df.DonchainChoppyStrategy(15, 16, 30, account, simple)
	if df.Signal.Signals == nil {
		fmt.Println("トレード結果がありません")
	}
	Result(df.Signal)
}
