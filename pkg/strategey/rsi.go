package strategey

import (
	"v1/pkg/execute"

	"github.com/markcheno/go-talib"
)

func (df *DataFrameCandle) RsiStrategy(period int, buyThread, sellThread float64) *execute.SignalEvents {
	lenCandles := len(df.Candles)
	if lenCandles <= period {
		return nil
	}

	signalEvents := execute.NewSignalEvents()
	values := talib.Rsi(df.Closes(), period)

	buySize := 0.0
	isHolding := false
	for i := 1; i < lenCandles; i++ {
		if values[i-1] == 0 || values[i-1] == 100 {
			continue
		}
		if values[i-1] < buyThread && values[i] >= buyThread {
			buySize = TradeSize(0.2) / df.Candles[i].Close
			signalEvents.Buy(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, buySize, true)
			isHolding = true
		}

		if values[i-1] > sellThread && values[i] <= sellThread {
			signalEvents.Sell(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, buySize, true)
			isHolding = false
		}
	}
	return signalEvents
}

func (df *DataFrameCandle) OptimizeRsi() (performance float64, bestPeriod int, bestBuyThread, bestSellThread float64) {
	bestPeriod = 14
	bestBuyThread, bestSellThread = 30.0, 70.0

	for period := 5; period < 25; period++ {
		signalEvents := df.RsiStrategy(period, bestBuyThread, bestSellThread)
		if signalEvents == nil {
			continue
		}
		profit := Profit(signalEvents)
		if performance < profit {
			performance = profit
			bestPeriod = period
			bestBuyThread = bestBuyThread
			bestSellThread = bestSellThread
		}
	}
	return performance, bestPeriod, bestBuyThread, bestSellThread
}
