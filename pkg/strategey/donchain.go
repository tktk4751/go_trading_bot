package strategey

import (
	"fmt"
	"v1/pkg/execute"
	"v1/pkg/indicator/indicators"

	"github.com/markcheno/go-talib"
)

func (df *DataFrameCandle) DonchainStrategy(period int) *execute.SignalEvents {

	StrategyName := "DBO"

	lenCandles := len(df.Candles)
	if lenCandles <= period {
		return nil
	}

	signalEvents := execute.NewSignalEvents()

	donchain := indicators.Donchain(df.Highs(), df.Low(), period)

	close := df.Closes()

	rsi := talib.Rsi(close, 14)

	for i := 1; i < lenCandles; i++ {
		if i < period {
			continue
		}
		if close[i] > donchain.High[i-1] {
			signalEvents.Buy(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, 1.0, true)

		}
		if close[i] < donchain.Low[i-1] || rsi[i] < 13 {
			signalEvents.Sell(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, 1.0, true)

		}

	}
	return signalEvents

}

func (df *DataFrameCandle) OptimizeDonchain() (performance float64, bestPeriod int) {
	bestPeriod = 40

	for period := 5; period < 300; period++ {

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

	fmt.Println("利益", performance, "最適なピリオド", bestPeriod)
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
