package strategey

// import (
// 	"v1/pkg/data/query"
// 	"v1/pkg/execute"

// 	"github.com/markcheno/go-talib"
// )

// var df, err = query.GetCandleData("SOLUSDT", "4h")

// func BackTestBb(n int, k float64) *execute.SignalEvents {
// 	cd

// 	if lenCandles <= n {
// 		return nil
// 	}

// 	signalEvents := &execute.SignalEvents{}
// 	bbUp, _, bbDown := talib.BBands(df.Closes(), n, k, k, 0)
// 	for i := 1; i < lenCandles; i++ {
// 		if i < n {
// 			continue
// 		}
// 		if bbDown[i-1] > df.Candles[i-1].Close && bbDown[i] <= df.Candles[i].Close {
// 			signalEvents.Buy(df.ProductCode, df.Candles[i].Time, df.Candles[i].Close, 1.0, false)
// 		}
// 		if bbUp[i-1] < df.Candles[i-1].Close && bbUp[i] >= df.Candles[i].Close {
// 			signalEvents.Sell(df.ProductCode, df.Candles[i].Time, df.Candles[i].Close, 1.0, false)
// 		}
// 	}
// 	return signalEvents
// }
