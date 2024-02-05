package main

import (
	"fmt"
	"time"
	"v1/pkg/strategey"
)

func main() {

	start := time.Now()

	strategey.RunEmaOptimize()
	// strategey.RunBacktestDonchainChoppy()
	// strategey.RunBacktestSuperTrend()

	// assetName := "TIAUSDT"
	// duration := "4h"

	// df, _ := strategey.GetCandleData(assetName, duration)

	// h := df.Highs()
	// l := df.Lows()
	// c := df.Closes()

	// a := risk.ChoppySlice(c, h, l)

	// fmt.Println(a)

	end := time.Now()

	// 処理時間を計算
	duration1 := end.Sub(start)

	// 処理時間を表示
	fmt.Printf("処理時間: %v\n", duration1)
}
