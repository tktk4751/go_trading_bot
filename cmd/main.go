package main

import (

	// "v1/pkg/analytics/metrics"

	// "v1/pkg/config"
	// "v1/pkg/db/models"
	// p "v1/pkg/management/position"

	// "fmt"
	// "log"
	// data "v1/pkg/db"
	"fmt"
	databace "v1/pkg/db"
)

// var path = "/home/lux/dev/go_trading_bot/pkg/data/spot/monthly/klines"

// var close []float64 = utils.GetClosePrice(path)

// var hloc = utils.GetCandleData(path)

// var side = randam_side()

// func randam_side() string {
// 	// Declare a local variable result to store the random side
// 	var result string

// 	for i := 0; i < len(close); i++ {

// 		n := rand.Intn(2)
// 		// Assign "BUY" or "SELL" to result
// 		if n == 0 {
// 			result = "BUY"
// 		} else {
// 			// Otherwise, assign "SELL" to result
// 			result = "SELL"
// 		}

//		}
//		// Return the value of result
//		return result
//	}
func main() {

	databace.GetCloseData("BTCUSDT", "4h")

	// var assets_names []string = []string{"BTCUSDT", "MATICUSDT", "PEPEUSDT", "ARBUSDT", "ETHUSDT", "XRPUSDT", "OPUSDT", "ATOMUSDT", "UNIUSDT", "SEIUSDT", "SUIUSDT", "TIAUSDT", "DOTUSDT", "NEARUSDT", "WLDUSDT", "XRPUSDT"}
	// var durations []string = []string{"1m", "15m", "30m", "4h"}
	// paths := data.GetRelativePaths()

	// groupedPaths := data.GroupAssetNamePaths(paths)

	// asset_data, err := data.LoadOHLCV(groupedPaths, assets_names, durations)
	// if err != nil {
	// 	log.Fatalf("Error loading OHLCV data: %v", err)
	// }

	// // DBに接続する関数を呼び出し
	// db, err := data.ConnectDB("./db/kline.db")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// // DBをクローズするのを遅延実行
	// defer db.Close()
	// // データをDBに保存する関数を呼び出し
	// err = data.SaveAssetDatas(db, asset_data)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// indicators.GetData()
	// 終了メッセージを表示

	// for key, paths := range groupedPaths {

	// 	fmt.Printf("{%s: %v}\n", key, paths)

	// }

	// for _, assetData := range asset_data {
	// 	fmt.Printf("Asset: %s, Duration: %s, OHLCV: %+v\n", assetData.AssetName, assetData.Duration, assetData.Data)
	// }

	// fmt.Println(asset_data)
	defer fmt.Println("メイン関数終了")

}

// data.GetAbsolutePaths()

// db := models.DbConnection

// env := config.GetEnv()

// var wr = metrics.Winrate_arg{
// 	Totall_wintrade: 100,
// 	Totall_trade:    200,
// }
// var winrate float64 = metrics.Calc_winrate(wr.Totall_wintrade, wr.Totall_trade)

// w := 0.4044
// r := 4.699
// d := 0.33

// position := p.PositionSizeCalculator{}

// risk_size := position.Risk_size_calculator(w, r, d) * 100

// sl := position.Stop_loss_price_calc(close, side)

// // management := money_management.PositionSizeCalculator{}
// // sl := management.Stop_loss_price_calc()

// // Call the KellyCriterion function and print the result
// fmt.Println(risk_size, "%")
// fmt.Println(env.TradeDuration, "DURATION")
// fmt.Println(sl, side, "EXITPRICE")
// // fmt.Println(env.ApiKey)
// fmt.Println(winrate)
// fmt.Println(db)
// fmt.Println(hloc)
