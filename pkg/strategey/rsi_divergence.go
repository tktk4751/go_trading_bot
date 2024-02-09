package strategey

// 過去n日間の高値を取得する
// その高値を記録したときのRSIの数値も取得する
// 現在の終値が､過去n日間の高値を突破する
// 現在のRSIの数値が高値のときのRSIよりも小さい場合､ダイバージェンスシグナルをTrueにする
// ダイバージェンスがTrueになったときに､rsi[i-1] > rsi[i]になったらエントリーする

// func Max(array []float64) float64 {
// 	// 配列の長さをチェックする
// 	if len(array) == 0 {
// 		return 0.0 // エラーを返す
// 	}
// 	// 最大値を初期化する
// 	max := array[0]
// 	// 配列の要素をループする
// 	for _, value := range array {
// 		// 最大値を更新する
// 		if value > max {
// 			max = value
// 		}
// 	}
// 	// 最大値を返す
// 	return max
// }

// func (df *DataFrameCandle) RsiDG(period int, duration int, account *trader.Account) *execute.SignalEvents {

// 	var StrategyName = "RSI_DIVERGENCE"
// 	lenCandles := len(df.Candles)
// 	if lenCandles <= period {
// 		return nil
// 	}

// 	signalEvents := execute.NewSignalEvents()

// 	high := df.Highs()
// 	hl3 := df.Hlc3()
// 	rsi := talib.Rsi(hl3, period)

// 	h := Max(high[len(high)-duration:])

// 	buySize := 0.0
// 	buyPrice := 0.0
// 	isBuyHolding := false
// 	dgSignal := false

// }
