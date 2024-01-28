package strategey

import (
	"fmt"
	"v1/pkg/analytics"
	"v1/pkg/execute"
	"v1/pkg/indicator/indicators"
)

const StrategyName = "DBO"

func TradeSize(persetege float64) float64 {

	size := AccountBalance * persetege
	return size
}

func (df *DataFrameCandle) DonchainStrategy(period int) *execute.SignalEvents {

	lenCandles := len(df.Candles)
	if lenCandles <= period {
		return nil
	}

	signalEvents := execute.NewSignalEvents()

	donchain := indicators.Donchain(df.Highs(), df.Low(), period)

	close := df.Closes()

	buySize := 0.0
	isHolding := false

	for i := 1; i < lenCandles; i++ {

		if i < period {
			continue
		}
		if close[i] > donchain.High[i-1] && !isHolding {
			buySize = TradeSize(0.2) / df.Candles[i].Close
			signalEvents.Buy(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, buySize, true)
			isHolding = true
		}
		if close[i] < donchain.Low[i-1] && isHolding {
			signalEvents.Sell(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, buySize, true)
			isHolding = false
		}

	}

	return signalEvents

}

func (df *DataFrameCandle) OptimizeProfitDonchain() (performance float64, bestPeriod int) {
	if df == nil {
		return 0.0, 0
	}
	bestPeriod = 40

	for period := 10; period < 333; period++ {

		signalEvents := df.DonchainStrategy(period)
		if signalEvents == nil {
			continue
		}
		profit := Profit(signalEvents)
		if performance < profit {
			performance = profit
			bestPeriod = period

		}

	}

	fmt.Println("最高利益", performance, "最適なピリオド", bestPeriod)

	return performance, bestPeriod
}

func (df *DataFrameCandle) OptimizeWinRateDonchain() (performance float64, bestPeriod int) {
	bestPeriod = 40

	for period := 10; period < 333; period++ {

		signalEvents := df.DonchainStrategy(period)
		if signalEvents == nil {
			continue
		}
		winrate := analytics.WinRate(signalEvents)
		if performance < winrate {
			performance = winrate
			bestPeriod = period

		}

	}

	fmt.Println("最高勝率", performance*100, "% ", "最適なピリオド", bestPeriod)
	return performance, bestPeriod
}

// func DonchainStrategeyBacktest(assetName string, duration string) ([]bool, []bool, []bool, []bool) {

// 	var ohlc, e = query.GetOHLCData(assetName, duration)
// 	if e != nil {
// 		log.Fatal(e)
// 	}

// 	var h []float64
// 	var l []float64
// 	var c []float64

// 	for _, data := range ohlc {
// 		h = append(h, data.High)
// 		l = append(l, data.Low)
// 		c = append(c, data.Close)
// 	}

// 	d := indicators.Donchain(h, l, 40)

// 	var buySignals []bool
// 	var sellSignals []bool
// 	var shortExitSignals []bool
// 	var longExitSignals []bool

// 	for i := range c {
// 		var buySignal bool = false
// 		var sellSignal bool = false
// 		var shortExitSignal bool = false
// 		var longExitSignal bool = false

// 		if c[i] > d.High[i] {
// 			buySignal = true
// 			shortExitSignal = true
// 		}

// 		if c[i] < d.Low[i] {
// 			sellSignal = true
// 			longExitSignal = true
// 		}

// 		buySignals = append(buySignals, buySignal)
// 		sellSignals = append(sellSignals, sellSignal)
// 		shortExitSignals = append(shortExitSignals, shortExitSignal)
// 		longExitSignals = append(longExitSignals, longExitSignal)
// 	}

// 	return buySignals, sellSignals, shortExitSignals, longExitSignals
// }
