package strategey

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"v1/pkg/analytics"
	"v1/pkg/execute"
	"v1/pkg/indicator/indicators"
	"v1/pkg/management/risk"
	"v1/pkg/trader"

	"github.com/c-bata/goptuna"

	"github.com/c-bata/goptuna/tpe"

	"github.com/markcheno/go-talib"
)

func (df *DataFrameCandle) StEmaStrategy(period1, period2 int, account *trader.Account, simple bool) *execute.SignalEvents {

	var StrategyName = "ST_EMA"
	lenCandles := len(df.Candles)
	if lenCandles <= period1 || lenCandles <= period2 {
		return nil
	}
	signalEvents := execute.NewSignalEvents()

	h := df.Highs()
	l := df.Lows()
	c := df.Closes()

	emaValue1 := talib.Ema(df.Hlc3(), period1)
	emaValue2 := talib.Ema(df.Hlc3(), period2)
	// rsi := talib.Rsi(df.Hlc3(), 14)

	buySize := 0.0
	buyPrice := 0.0
	slRatio := 0.9
	index := risk.ChoppySlice(70, c, h, l)
	choppyEma := risk.ChoppyEma(index, 5)

	// hlc3 := df.Hlc3()

	st, _ := indicators.SuperTrend(50, 5.5, h, l, c)

	isBuyHolding := false
	for i := 1; i < lenCandles; i++ {
		if i < period1 || i < period2 || i >= len(choppyEma) {
			continue
		}
		buyCondition := c[i] > st.SuperTrend[i]
		sellCondition := c[i] < st.SuperTrend[i]

		if emaValue1[i-1] < emaValue2[i-1] && emaValue1[i] >= emaValue2[i] && choppyEma[i] > 50 && buyCondition && !isBuyHolding {

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
		if emaValue1[i-1] > emaValue2[i-1] && emaValue1[i] <= emaValue2[i] || (df.Candles[i].Close <= buyPrice*slRatio) && sellCondition && isBuyHolding {
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

func (df *DataFrameCandle) OptimizeStEma2() (performance float64, bestPeriod1 int, bestPeriod2 int) {

	// オブジェクティブ関数を定義
	objective := func(trial goptuna.Trial) (float64, error) {
		// ハイパーパラメータの候補をサンプリング
		period1, _ := trial.SuggestStepInt("period1", 3, 13, 1)
		period2, _ := trial.SuggestStepInt("period2", 13, 89, 1)

		account := trader.NewAccount(1000) // Move this line inside the objective function
		signalEvents := df.StEmaStrategy(period1, period2, account, simple)

		if signalEvents == nil {
			return 0.0, nil
		}

		if analytics.TotalTrades(signalEvents) < 10 {
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
	bestPeriod1 = params["period1"].(int)
	bestPeriod2 = params["period2"].(int)

	fmt.Println("最高パフォーマンス", performance, "最適な短期線", bestPeriod1, "最適な長期線", bestPeriod2)

	return performance, bestPeriod1, bestPeriod2
}

func (df *DataFrameCandle) OptimizeStEma() (performance float64, bestPeriod1 int, bestPeriod2 int) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	bestPeriod1 = 5
	bestPeriod2 = 21

	limit := 3000
	slots := make(chan struct{}, limit)

	// a := trader.NewAccount(1000)
	// marketDefault, _ := BuyAndHoldingStrategy(a)

	var mu sync.Mutex
	var wg sync.WaitGroup
	for period1 := 3; period1 < 13; period1 += 1 {
		for period2 := 5; period2 < 34; period2 += 1 {

			wg.Add(1)
			slots <- struct{}{}

			go func(period1 int, period2 int) {

				defer wg.Done()
				account := trader.NewAccount(1000) // Move this line inside the goroutine
				signalEvents := df.StEmaStrategy(period1, period2, account, simple)

				if signalEvents == nil {
					return
				}

				if analytics.TotalTrades(signalEvents) < 15 {
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

				p := analytics.SortinoRatio(signalEvents, 0.02)
				// p := analytics.Prr(signalEvents)
				mu.Lock()
				if performance == 0 || performance < p {
					performance = p
					bestPeriod1 = period1
					bestPeriod2 = period2

				}
				<-slots
				mu.Unlock()

			}(period1, period2)

		}
	}

	wg.Wait()

	if bestPeriod1 > bestPeriod2 {
		log.Fatalf("数値が逆転しています")
	}

	fmt.Println("最高パフォーマンス", performance, "最適な短期線", bestPeriod1, "最適な長期線", bestPeriod2)

	return performance, bestPeriod1, bestPeriod2
}

func RunStEmaOptimize() {

	df, account, _ := RadyBacktest()

	performance, bestPeriod1, bestPeriod2 := df.OptimizeStEma()

	if performance > 0 {

		df.Signal = df.StEmaStrategy(bestPeriod1, bestPeriod2, account, simple)
		if df.Signal.Signals == nil {
			fmt.Println("トレード結果がありません")
		}
		Result(df.Signal)

	} else {
		fmt.Println("💸マイナスです")
		df.Signal = df.EmaChoppyStrategy(bestPeriod1, bestPeriod2, account, simple)
		Result(df.Signal)

	}

}

func RunStEmaOptimize2() {

	df, account, _ := RadyBacktest()

	performance, bestPeriod1, bestPeriod2 := df.OptimizeStEma2()

	if performance > 0 {

		df.Signal = df.StEmaStrategy(bestPeriod1, bestPeriod2, account, simple)
		if df.Signal.Signals == nil {
			fmt.Println("トレード結果がありません")
		}
		Result(df.Signal)

	} else {
		fmt.Println("💸マイナスです")
		df.Signal = df.EmaChoppyStrategy(bestPeriod1, bestPeriod2, account, simple)
		Result(df.Signal)

	}

}

func StEmaBacktest() {

	df, account, _ := RadyBacktest()

	df.Signal = df.StEmaStrategy(9, 47, account, simple)
	if df.Signal.Signals == nil {
		fmt.Println("トレード結果がありません")
	}
	Result(df.Signal)

}
