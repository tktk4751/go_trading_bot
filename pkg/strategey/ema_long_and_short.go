package strategey

// import (
// 	"fmt"
// 	"log"
// 	"math"
// 	"v1/pkg/config"
// 	"v1/pkg/execute"
// 	"v1/pkg/trader"

// 	"github.com/markcheno/go-talib"
// )

// func (df *DataFrameCandle) EmaStrategy2(period1, period2 int, account *trader.Account) *execute.SignalEvents {

// 	var StrategyName = "EMA"
// 	lenCandles := len(df.Candles)
// 	if lenCandles <= period1 || lenCandles <= period2 {
// 		return nil
// 	}
// 	signalEvents := execute.NewSignalEvents()
// 	emaValue1 := talib.Ema(df.Closes(), period1)
// 	emaValue2 := talib.Ema(df.Closes(), period2)
// 	rsiValue := talib.Rsi(df.Closes(), 14)

// 	buySize := 0.0
// 	buyPrice := 0.0
// 	sellSize := 0.0
// 	sellPrice := 0.0
// 	slRatio := 0.9

// 	isBuyHolding := false
// 	isSellHolding := false
// 	for i := 1; i < lenCandles; i++ {
// 		if i < period1 || i < period2 {
// 			continue
// 		}
// 		// EMAのゴールデンクロスで買いポジションをオープンする
// 		if emaValue1[i-1] < emaValue2[i-1] && emaValue1[i] >= emaValue2[i] && !isBuyHolding {
// 			accountBalance := account.GetBalance()
// 			buySize = account.TradeSize(riskSize) / df.Candles[i].Close
// 			buyPrice = df.Candles[i].Close

// 			// ショートポジションがあればクローズする
// 			if isSellHolding {
// 				signalEvents.Exit(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, sellSize, accountBalance, false)
// 				isSellHolding = false
// 				sellSize = 0.0
// 				account.PositionSize = buySize
// 			}
// 			if account.Buy(df.Candles[i].Close, buySize) {
// 				signalEvents.Buy(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, buySize, accountBalance, false)
// 				isBuyHolding = true

// 			}
// 		}
// 		// EMAのデッドクロスでショートポジションをオープンする
// 		if emaValue1[i-1] > emaValue2[i-1] && emaValue1[i] <= emaValue2[i] && !isSellHolding {
// 			accountBalance := account.GetBalance()
// 			sellSize = -account.TradeSize(riskSize) / df.Candles[i].Close
// 			sellPrice = df.Candles[i].Close

// 			// ロングポジションがあればクローズする
// 			if isBuyHolding {
// 				signalEvents.Exit(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, buySize, accountBalance, false)
// 				isBuyHolding = false
// 				buySize = 0.0
// 				account.PositionSize = sellSize
// 			}

// 			if account.Sell(df.Candles[i].Close) {
// 				signalEvents.Sell(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, math.Abs(sellSize), accountBalance, false)
// 				isSellHolding = true

// 			}
// 		}
// 		// RSIが30以下で買いポジションをクローズする
// 		if rsiValue[i] < 30.0 && isBuyHolding {
// 			accountBalance := account.GetBalance()
// 			signalEvents.Exit(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, buySize, accountBalance, false)
// 			isBuyHolding = false
// 			buySize = 0.0
// 			account.PositionSize = buySize
// 		}
// 		// RSIが70以上でショートポジションをクローズする
// 		if rsiValue[i] > 70.0 && isSellHolding {
// 			accountBalance := account.GetBalance()
// 			signalEvents.Exit(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, sellSize, accountBalance, false)
// 			isSellHolding = false
// 			sellSize = 0.0
// 			account.PositionSize = sellSize
// 		}
// 		// ストップロスでポジションをクローズする
// 		if (df.Candles[i].Close <= buyPrice*slRatio) && isBuyHolding {
// 			accountBalance := account.GetBalance()
// 			signalEvents.Exit(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, buySize, accountBalance, false)
// 			isBuyHolding = false
// 			buySize = 0.0
// 			account.PositionSize = buySize
// 		}
// 		if (df.Candles[i].Close >= sellPrice*slRatio) && isSellHolding {
// 			accountBalance := account.GetBalance()
// 			signalEvents.Exit(StrategyName, df.AssetName, df.Duration, df.Candles[i].Date, df.Candles[i].Close, sellSize, accountBalance, false)
// 			isSellHolding = false
// 			sellSize = 0.0
// 			account.PositionSize = sellSize
// 		}
// 	}
// 	return signalEvents
// }
// func RunBacktestEma2() {

// 	var err error

// 	// account := trader.NewAccount(1000)
// 	btcfg, err := config.Yaml()
// 	if err != nil {
// 		log.Fatalf("error: %v", err)
// 	}

// 	fmt.Println(btcfg.AssetName)

// 	assetName := btcfg.AssetName
// 	duration := btcfg.Dration
// 	// limit := btcfg.Limit

// 	account := trader.NewAccount(1000)

// 	df, _ := GetCandleData(assetName, duration)

// 	df.Signal = df.EmaStrategy2(168, 245, account)
// 	Result(df.Signal)

// }
