package indicators

import (
	"math"
)

// ATR関数の定義
func Atr(prices [][]float64, n int) float64 {

	// TRのスライスを作成
	tr := make([]float64, len(prices))
	// 最初のTRは高値と安値の差とする
	tr[0] = prices[0][0] - prices[0][1]
	// 2日目以降のTRを計算
	for i := 1; i < len(prices); i++ {
		// 高値と安値の差
		hl := prices[i][0] - prices[i][1]
		// 高値と前日終値の差の絶対値
		hc := math.Abs(prices[i][0] - prices[i-1][2])
		// 安値と前日終値の差の絶対値
		lc := math.Abs(prices[i][1] - prices[i-1][2])
		// 3つの値の最大値をTRとする
		tr[i] = math.Max(hl, math.Max(hc, lc))
	}
	// 最初のn日間のTRの平均をATRとする
	atr := 0.0
	for i := 0; i < n; i++ {
		atr += tr[i]
	}
	atr /= float64(n)
	// n+1日目以降のATRを計算
	for i := n; i < len(prices); i++ {
		// ATR = (前日のATR * (n-1) + 今日のTR) / n
		atr = (atr*float64(n-1) + tr[i]) / float64(n)
	}
	// ATRの値を返す
	return atr
}

// package indicators

// import (
// 	"math"
// )

// // ATR構造体の定義
// type ATR struct {
// 	High  float64 // 高値
// 	Low   float64 // 安値
// 	Close float64 // 終値
// 	TR    float64 // 真の変動幅
// 	N     int     // 採用本数
// 	Value float64 // ATRの値
// }

// // ATR関数の定義
// func Atr(prices [][]float64, n int) ATR {

// 	// ATR構造体のインスタンスを作成
// 	atr := ATR{}
// 	// Nをセット
// 	atr.N = n
// 	// 最初の高値、安値、終値をセット
// 	atr.High = prices[0][0]
// 	atr.Low = prices[0][1]
// 	atr.Close = prices[0][2]
// 	// 最初のTRは高値と安値の差、高値と前日終値の差、安値と前日終値の差のうち、最大のものとする
// 	atr.TR = math.Max(atr.High-atr.Low, math.Max(math.Abs(atr.High-atr.Close), math.Abs(atr.Low-atr.Close)))
// 	// 最初のATRはTRとする
// 	atr.Value = atr.TR
// 	// 2日目以降のATRを計算
// 	for i := 1; i < len(prices); i++ {
// 		// 高値、安値、終値を更新
// 		atr.High = prices[i][0]
// 		atr.Low = prices[i][1]
// 		atr.Close = prices[i][2]
// 		// 高値と安値の差
// 		hl := atr.High - atr.Low
// 		// 高値と前日終値の差の絶対値
// 		hc := math.Abs(atr.High - prices[i-1][2])
// 		// 安値と前日終値の差の絶対値
// 		lc := math.Abs(atr.Low - prices[i-1][2])
// 		// 3つの値の最大値をTRとする
// 		atr.TR = math.Max(hl, math.Max(hc, lc))
// 		// ATR = (前日のATR * (n-1) + 今日のTR) / n
// 		atr.Value = (atr.Value*float64(n-1) + atr.TR) / float64(n)
// 	}
// 	// ATR構造体を返す
// 	return atr
// }
