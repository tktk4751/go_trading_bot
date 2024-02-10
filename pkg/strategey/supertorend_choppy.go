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

	"github.com/c-bata/goptuna"
	"github.com/c-bata/goptuna/tpe"
)

func (df *DataFrameCandle) SuperTrendChoppyStrategy(atrPeriod int, factor float64, account *trader.Account, simple bool) *execute.SignalEvents {

	var StrategyName = "SUPERTREND_CHOPPY"
	// var err error

	lenCandles := len(df.Candles)

	if lenCandles <= atrPeriod {
		return nil
	}

	signalEvents := execute.NewSignalEvents()

	h := df.Highs()
	l := df.Lows()
	c := df.Closes()
	// hlc3 := df.Hlc3()

	superTrend, _ := indicators.SuperTrend(atrPeriod, factor, h, l, c)

	// stUp := superTrend.UpperBand
	// stLow := superTrend.UpperBand
	st := superTrend.SuperTrend

	// rsi := talib.Rsi(hlc3, 14)

	buySize := 0.0
	buyPrice := 0.0
	slRatio := 0.9

	index := risk.ChoppySlice(70, c, h, l)
	choppyEma := risk.ChoppyEma(index, 5)

	isHolding := false

	for i := 1; i < len(choppyEma); i++ {

		if i < atrPeriod {
			// fmt.Printf("Skipping iteration %d due to insufficient data.\n", i)
			continue
		}
		if (c[i-1] < st[i-1] && c[i] >= st[i]) && choppyEma[i] > 50 && !isHolding {

			// fee := 1 - 0.01
			if simple {
				buySize = account.SimpleTradeSize(1)
				buyPrice = c[i]
				accountBalance := account.GetBalance()

				signalEvents.Buy(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, buySize, accountBalance, false)
				isHolding = true

			} else {
				buySize = account.TradeSize(riskSize) / df.Candles[i].Close
				buyPrice = c[i]
				accountBalance := account.GetBalance()
				if account.Buy(df.Candles[i].Close, buySize) {
					signalEvents.Buy(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, buySize, accountBalance, false)
					isHolding = true
				}
			}

		}
		if ((c[i-1] > st[i-1] && c[i] <= st[i]) || (c[i] <= buyPrice*slRatio)) && isHolding {
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

	// fmt.Println(signalEvents)
	return signalEvents
}

func (df *DataFrameCandle) OptimizeSuperTrend2() (performance float64, bestAtrPeriod int, bestFactor float64) {

	// オブジェクティブ関数を定義
	objective := func(trial goptuna.Trial) (float64, error) {
		// ハイパーパラメータの候補をサンプリング
		atrPeriod, _ := trial.SuggestStepInt("atrPeriod", 10, 60, 2)
		factorInt, _ := trial.SuggestInt("factor", 4, 16)
		factor := float64(factorInt) * 0.5

		account := trader.NewAccount(1000) // Move this line inside the objective function

		marketDefault, _ := BuyAndHoldingStrategy(account)
		signalEvents := df.SuperTrendChoppyStrategy(atrPeriod, factor, account, simple)

		if signalEvents == nil {
			return 0.0, nil
		}

		if analytics.TotalTrades(signalEvents) < 30 {
			return 0.0, nil
		}

		if analytics.TotalNetProfit(signalEvents) < marketDefault {
			return 0.0, nil
		}

		p := analytics.SortinoRatio(signalEvents, 0.02)
		// p := analytics.Prr(signalEvents)
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
	err = study.Optimize(objective, 1000)
	if err != nil {
		panic(err)
	}

	// 最適化結果の取得
	v, _ := study.GetBestValue()
	params, _ := study.GetBestParams()
	performance = v
	bestAtrPeriod = params["atrPeriod"].(int)
	bestFactor = float64(params["factor"].(int)) * 0.5

	fmt.Println("最高のパフォーマンス", performance, "最適なATR", bestAtrPeriod, "最適なファクター", bestFactor)

	return performance, bestAtrPeriod, bestFactor

}

func (df *DataFrameCandle) OptimizeSuperTrend() (performance float64, bestAtrPeriod int, bestFactor float64) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	bestAtrPeriod = 21
	bestFactor = 3.0

	limit := 1000
	slots := make(chan struct{}, limit)

	a := trader.NewAccount(1000)
	marketDefault, _ := BuyAndHoldingStrategy(a)

	var mu sync.Mutex
	var wg sync.WaitGroup

	for atrPeriod := 13; atrPeriod < 60; atrPeriod += 2 {
		for factor := 2.0; factor < 8.0; factor += 0.5 {

			wg.Add(1)
			slots <- struct{}{}

			go func(atrPeriod int, factor float64) {
				defer wg.Done()
				account := trader.NewAccount(1000) // Move this line inside the goroutine
				signalEvents := df.SuperTrendChoppyStrategy(atrPeriod, factor, account, simple)

				if signalEvents == nil {
					return
				}

				if analytics.TotalTrades(signalEvents) < 30 {
					<-slots
					return
				}

				if analytics.TotalNetProfit(signalEvents) < marketDefault {
					<-slots
					return
				}

				// if analytics.SQN(signalEvents) < 3.2 {
				// 	<-slots
				// 	return
				// }

				// if analytics.PayOffRatio(signalEvents) < 1 {
				// <-slots
				// 	return
				// }

				// p := analytics.SortinoRatio(signalEvents, 0.02)
				p := analytics.SQN(signalEvents)
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

	fmt.Println("最高のパフォーマンス", performance, "最適なATR", bestAtrPeriod, "最適なファクター", bestFactor)

	return performance, bestAtrPeriod, bestFactor
}

func RunSTOptimize() {

	df, account, _ := RadyBacktest()

	performance, bestAtrPeriod, bestFactor := df.OptimizeSuperTrend()

	if performance > 0 {

		df.Signal = df.SuperTrendChoppyStrategy(bestAtrPeriod, bestFactor, account, simple)
		Result(df.Signal)

	} else {
		fmt.Println("💸マイナスです")
		df.Signal = df.SuperTrendChoppyStrategy(bestAtrPeriod, bestFactor, account, simple)
		Result(df.Signal)

	}

}

func RunSTOptimize2() {

	df, account, _ := RadyBacktest()

	performance, bestAtrPeriod, bestFactor := df.OptimizeSuperTrend2()

	if performance > 0 {

		df.Signal = df.SuperTrendChoppyStrategy(bestAtrPeriod, bestFactor, account, simple)
		Result(df.Signal)

	} else {
		fmt.Println("💸マイナスです")
		df.Signal = df.SuperTrendChoppyStrategy(bestAtrPeriod, bestFactor, account, simple)
		Result(df.Signal)

	}

}

func SuperTrendBacktest() {

	df, account, _ := RadyBacktest()

	df.Signal = df.SuperTrendChoppyStrategy(26, 3.5, account, simple)
	if df.Signal.Signals == nil {
		fmt.Println("トレード結果がありません")
	}
	Result(df.Signal)

}
