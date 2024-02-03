package strategey

import (
	"log"
	"v1/pkg/config"
	dbquery "v1/pkg/data/query"
	"v1/pkg/trader"
)

// func getStrageyNameBuyAndHold() string {
// 	return "B&H"
// }

func BuyAndHoldingStrategy(account *trader.Account) (profit float64, multiple float64) {
	// if df == nil {
	// 	return 0.0, 0.0
	// }

	btcfg, err := config.Yaml()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	assetName := btcfg.AssetName
	duration := btcfg.Dration

	close, _ := dbquery.GetCloseData(assetName, duration)

	lenCandles := len(close)
	if lenCandles < 2 {
		return 0.0, 0.0
	}
	// account := trader.NewAccount(1000)

	buySize := account.TradeSize(1) / close[0]
	account.HolderBuy(close[0], buySize)

	account.Sell(close[len(close)-1])

	profit = account.Balance - initialBalance
	multiple = account.Balance / initialBalance
	return profit, multiple
}

// func RunBacktestBuyAndHolding() {

// 	var err error

// 	// account := trader.NewAccount(1000)
// 	btcfg, err := config.Yaml()
// 	if err != nil {
// 		log.Fatalf("error: %v", err)
// 	}

// 	fmt.Println(btcfg.AssetName)

// 	strategyName := getStrageyNameBuyAndHold()
// 	assetName := btcfg.AssetName
// 	duration := btcfg.Dration

// 	account := trader.NewAccount(1000)

// 	df, _ := GetCandleData(assetName, duration)

// 	tableName := strategyName + "_" + assetName + "_" + duration

// 	_, err = execute.CreateDBTable(tableName)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	df.Signal = df.BuyAndHoldingStrategy(account)

// 	Result(df.Signal)

// }
