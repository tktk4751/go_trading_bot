package main

import (
	"fmt"
	"log"
	"net/http"
	"v1/pkg/data/query"
	"v1/pkg/execute"
	// "v1/pkg/analytics/metrics"
	// "v1/pkg/config"
	// "v1/pkg/db/models"
	// p "v1/pkg/management/position"
	// data "v1/pkg/data/utils"
	// chart "v1/pkg/charts"
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

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
func main() {

	strategyName := "test"
	assetName := "SOLUSDT"
	duration := "4H"

	db, err := execute.DBOpen(strategyName, assetName, duration)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	s := execute.NewSignalEvents()
	p, _ := query.GetKlineData(assetName, duration)

	c1 := p[0].Close
	c2 := p[2000].Close
	err = s.Buy(db, strategyName, assetName, duration, p[8].Date, c1, 1.0)
	if err != nil {
		log.Println("Error buying:", err)
	}
	err = s.Sell(db, strategyName, assetName, duration, p[3].Date, c2, 1.0)
	if err != nil {
		log.Println("Error selling:", err)
	}
	defer fmt.Println("メイン関数終了")

	// チャート呼び出し
	// var c chart.CandleStickChart
	// c.CandleStickChart()

	// query.GetCloseData("BTCUSDT", "4h")

	// var assets_names []string = []string{"RUNEUSDT", "BTCUSDT", "AAVEUSDT", "ORDIUSDT", "SANUSDT", "LTCUSDT", "OKBUSDT", "ASTRUSDT", "MNTUSDT", "FTMUSDT", "SNXUSDT", "DYDXUSDT", "BONKUSDT", "LUNAUSDT", "MAGICUSDT", "XLMUSDT", "DOGEUSDT", "TRSUSDT", "LINKUSDT", "TONUSDT", "ISPUSDT", "BONKUSDT", "GMXUSDT", "INJUSDT", "ETHUSDT", "SOLUSDT", "AVAXUSDT", "MATICUSDT", "ATOMUSDT", "UNIUSDT", "ARBUSDT", "OPUSDT", "PEPEUSDT", "SEIUSDT", "SUIUSDT", "TIAUSDT", "WLDUSDT", "XRPUSDT", "NEARUSDT", "DOTUSDT", "APTUSDT", "XMRUSDT", "LDOUSDT", "FILUSDT", "KASUSDT", "STXUSDT", "RNDRUSDT", "GRTUSDT"}

	// var durations []string = []string{"1m", "3m", "5m", "15m", "30m", "1h", "2h", "4h", "6h", "8h", "12h"}
	// paths := data.GetRelativePaths()

	// groupedPaths := data.GroupAssetNamePaths(paths)

	// asset_data, err := data.LoadOHLCV(groupedPaths, assets_names, durations)
	// if err != nil {
	// 	log.Fatalf("Error loading OHLCV data: %v", err)
	// }

	// // data.SaveAssetDatasCSV(asset_data)

	// // DBに接続する関数を呼び出し
	// db, err := data.ConnectDB("./db/kline.db")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// // DBをクローズするのを遅延実行
	// defer db.Close()
	// // データをDBに保存する関数を呼び出し
	// err = data.SaveAssetDatasDB(db, asset_data)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// for key, paths := range groupedPaths {

	// 	fmt.Printf("{%s: %v}\n", key, paths)

	// }

	// for _, assetData := range asset_data {
	// 	fmt.Printf("Asset: %s, Duration: %s, OHLCV: %+v\n", assetData.AssetName, assetData.Duration, assetData.Data)
	// }

	// fmt.Println(asset_data)

	// fs := http.FileServer(http.Dir("pkg/charts/html"))
	// log.Println("running server at http://localhost:8089")
	// log.Fatal(http.ListenAndServe("localhost:8089", logRequest(fs)))

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
