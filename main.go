package main

import (
	"fmt"
	"time"
	dbquery "v1/pkg/data/query"
	"v1/pkg/indicator/indicators"
)

func main() {

	start := time.Now()

	// 	close, _ := dbquery.GetCloseData("SOLUSDT", "1h")

	// 	high, _ := dbquery.GetHighData("SOLUSDT", "1h")

	// 	low, _ := dbquery.GetLowData("SOLUSDT", "1h")

	// 	index := risk.ChoppySlice(close, high, low)

	// 	e := risk.ChoppyEma(index)
	// 	fmt.Println(e)
	// strategey.RunBacktestEma()
	// strategey.RunBacktestEmaChoppy()

	// strategey.RunBacktestDonchain()
	// strategey.RunBacktestDonchainChoppy()

	// c, _ := strategey.GetCsvDataFrame("BTCUSDT", "4h", "2022-05", "2023-12")

	// fmt.Println(c)

	a := "SOLUSDT"
	d := "4h"

	h, _ := dbquery.GetHighData(a, d)
	l, _ := dbquery.GetLowData(a, d)
	c, _ := dbquery.GetCloseData(a, d)

	st, _ := indicators.SuperTrend(21, 3.0, h, l, c)

	fmt.Println(st.SuperTrend)

	end := time.Now()

	// 処理時間を計算
	duration := end.Sub(start)

	// 処理時間を表示
	fmt.Printf("処理時間: %v\n", duration)
}
