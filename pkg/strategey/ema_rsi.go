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
	"github.com/google/uuid"

	"github.com/c-bata/goptuna/tpe"

	"github.com/markcheno/go-talib"
)

func (df *DataFrameCandle) EmaRsiStrategy(period1 int, account *trader.Account, simple bool) *execute.SignalEvents {

	var StrategyName = "EMA_RSI"
	lenCandles := len(df.Candles)
	if lenCandles <= period1 {
		return nil
	}
	signalEvents := execute.NewSignalEvents()

	c := df.Closes()

	ema1 := talib.Ema(df.Hlc3(), period1)
	// ema2 := talib.Ema(df.Hlc3(), period2)
	rsi := talib.Rsi(df.Hlc3(), 5)

	var buySignalId uuid.UUID
	var sellSignalId uuid.UUID

	buySize := 0.0
	sellSize := 0.0
	buyPrice := 0.0
	sellPrice := 0.0
	longSlRatio := 0.9
	shortSlRatio := 1.1
	index := risk.ChoppySlice(70, df.Closes(), df.Highs(), df.Lows())
	choppyEma := risk.ChoppyEma(index, 5)

	isBuyHolding := false
	isSellHolding := false
	for i := 1; i < lenCandles; i++ {
		if i < period1 || i >= len(choppyEma) {
			continue
		}

		if ema1[i] < c[i] && rsi[i-1] < 30 && rsi[i] >= 30 && choppyEma[i] > 50 && !isBuyHolding {

			buyPrice = c[i]
			buySize = account.TradeSize(riskSize) / c[i]
			buySignalId = uuid.New()
			accountBalance := account.GetBalance()
			if account.Entry("BUY", buyPrice, buySize, 0.01) {
				signalEvents.Buy(buySignalId, StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, buySize, accountBalance, false)
				isSellHolding = false
				isBuyHolding = true

			}

		}
		if (ema1[i] > c[i] || c[i] <= sellPrice*shortSlRatio) && isBuyHolding {

			if account.Exit("BUY", c[i]) {
				accountBalance := account.GetBalance()
				signalEvents.Close(buySignalId, StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, buySize, accountBalance, false)
				isBuyHolding = false
				isSellHolding = false
				buySize = 0.0
				account.PositionSize = buySize
			}

		}
		//SELLのエントリーかエグジットに問題がある｡
		if ema1[i] > c[i] && rsi[i-1] > 75 && rsi[i] <= 75 && choppyEma[i] > 50 && !isSellHolding {

			sellPrice = c[i]
			sellSize = account.TradeSize(riskSize) / c[i]
			accountBalance := account.GetBalance()
			sellSignalId = uuid.New()
			if account.Entry("SELL", sellPrice, sellSize, 0.01) {
				signalEvents.Sell(sellSignalId, StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, sellSize, accountBalance, false)
				isBuyHolding = false
				isSellHolding = true

				// account.PositionSize = buySize

			}

			if (ema1[i] < c[i] || (c[i] <= buyPrice*longSlRatio)) && isSellHolding {

				if account.Exit("SELL", c[i]) {
					sellSize = account.PositionSize
					accountBalance := account.GetBalance()
					signalEvents.Close(sellSignalId, StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, sellSize, accountBalance, false)
					isSellHolding = false
					isBuyHolding = false
					sellSize = 0
					account.PositionSize = sellSize
				}

			}
		}
	}
	return signalEvents
}

func (df *DataFrameCandle) OptimizeEmaRsi2() (performance float64, bestPeriod int) {

	// オブジェクティブ関数を定義
	objective := func(trial goptuna.Trial) (float64, error) {
		// ハイパーパラメータの候補をサンプリング
		period1, _ := trial.SuggestStepInt("period1", 75, 250, 1)

		account := trader.NewAccount(1000) // Move this line inside the objective function
		marketDefault, _ := BuyAndHoldingStrategy(account)

		signalEvents := df.EmaRsiStrategy(period1, account, simple)

		if signalEvents == nil {
			return 0.0, nil
		}

		if analytics.TotalTrades(signalEvents) < 10 {
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
	err = study.Optimize(objective, 500)
	if err != nil {
		panic(err)
	}

	// 最適化結果の取得
	v, _ := study.GetBestValue()
	params, _ := study.GetBestParams()
	performance = v
	bestPeriod = params["period1"].(int)

	fmt.Println("最高パフォーマンス", performance, "最適なMA", bestPeriod)

	return performance, bestPeriod
}

func (df *DataFrameCandle) OptimizeEmaRsi() (performance float64, bestPeriod int) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	bestPeriod = 10

	limit := 1000
	slots := make(chan struct{}, limit)

	// a := trader.NewAccount(1000)
	// marketDefault, _ := BuyAndHoldingStrategy(a)

	var mu sync.Mutex
	var wg sync.WaitGroup
	for period := 75; period < 250; period += 5 {

		wg.Add(1)
		slots <- struct{}{}

		go func(period int) {

			defer wg.Done()
			account := trader.NewAccount(1000) // Move this line inside the goroutine
			signalEvents := df.EmaRsiStrategy(period, account, simple)

			if signalEvents == nil {
				return
			}

			if analytics.TotalTrades(signalEvents) < 30 {
				<-slots
				return
			}

			// if analytics.TotalNetProfit(signalEvents) < marketDefault {
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
			p := analytics.SQN(signalEvents)
			mu.Lock()
			if performance == 0 || performance < p {
				performance = p
				fmt.Println(performance)
				bestPeriod = period

			}
			<-slots
			mu.Unlock()

		}(period)

	}

	wg.Wait()

	fmt.Println("最高パフォーマンス", performance, "最適なMA", bestPeriod)

	return performance, bestPeriod
}

func RunEmaRsiOptimize() {

	df, account, _ := RadyBacktest()

	performance, bestPeriod := df.OptimizeEmaRsi()

	if performance > 0 {

		df.Signal = df.EmaRsiStrategy(bestPeriod, account, simple)
		if df.Signal.Signals == nil {
			fmt.Println("トレード結果がありません")
		}
		Result(df.Signal)

	} else {
		fmt.Println("💸マイナスです")
		df.Signal = df.EmaRsiStrategy(bestPeriod, account, simple)
		Result(df.Signal)

	}

}

// func RunEmaRsiOptimize2() {

// 	df, account, _ := RadyBacktest()

// 	performance, bestPeriod1, bestPeriod2 := df.OptimizeEma2()

// 	if performance > 0 {

// 		df.Signal = df.EmaChoppyStrategy(bestPeriod1, bestPeriod2, account, simple)
// 		Result(df.Signal)
// 		if df.Signal.Signals == nil {
// 			fmt.Println("トレード結果がありません")
// 		}

// 	} else {
// 		fmt.Println("💸マイナスです")
// 		df.Signal = df.EmaChoppyStrategy(bestPeriod1, bestPeriod2, account, simple)

// 		Result(df.Signal)

// 	}

// }

func EmaRsiBacktest() {

	df, account, _ := RadyBacktest()

	df.Signal = df.EmaRsiStrategy(100, account, simple)
	if df.Signal.Signals == nil {
		fmt.Println("トレード結果がありません")
	}
	Result(df.Signal)

}
