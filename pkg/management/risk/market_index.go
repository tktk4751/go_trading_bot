package risk

import "math"

// ChoppyMarketIndex関数は、市場の方向性とボラティリティを測る指標を返す
func ChoppyMarketIndex(close []float64, high []float64, low []float64) float64 {
	// 配列の長さをチェックする
	if len(close) < 30 || len(high) < 30 || len(low) < 30 {
		return 0.0 // エラーを返す
	}
	// 現在の終値と30日前の終値の差の絶対値を計算する
	direction := math.Abs(close[len(close)-1] - close[len(close)-30])
	// 過去30日間の最高値と最安値の差を計算する
	volatility := Max(high[len(high)-30:]) - Min(low[len(low)-30:])
	// 市場の方向性とボラティリティのバランスを計算する
	index := direction / volatility * 100.0
	// 指標を返す
	return index
}

// Max関数は、配列の中の最大値を返す
func Max(array []float64) float64 {
	// 配列の長さをチェックする
	if len(array) == 0 {
		return 0.0 // エラーを返す
	}
	// 最大値を初期化する
	max := array[0]
	// 配列の要素をループする
	for _, value := range array {
		// 最大値を更新する
		if value > max {
			max = value
		}
	}
	// 最大値を返す
	return max
}

// Min関数は、配列の中の最小値を返す
func Min(array []float64) float64 {
	// 配列の長さをチェックする
	if len(array) == 0 {
		return 0.0 // エラーを返す
	}
	// 最小値を初期化する
	min := array[0]
	// 配列の要素をループする
	for _, value := range array {
		// 最小値を更新する
		if value < min {
			min = value
		}
	}
	// 最小値を返す
	return min
}