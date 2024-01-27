package strategey

// import (
// 	"database/sql"
// 	"log"
// 	"time"
// 	"v1/pkg/data/query"
// 	"v1/pkg/execute"
// 	"v1/pkg/indicator/indicators"
// )

// func DonchainStrategeyLong(db *sql.DB, assetName string, duration string, strategyName string, date time.Time, price float64, size float64, save bool) (bool, bool) {

// 	datas := query.GetOHLCData(assetName, duration)

// 	s := execute.NewSignalEvents()

// 	var buy bool = s.Buy(db, strategyName, assetName, duration, date, price, size, save)
// 	var sell bool = s.Sell(db, strategyName, assetName, duration, date, price, size, save)

// 	var shortExitSignal = false
// 	var longExitSignal = false

// 	var ohlc, e = query.GetOHLCData(assetName, duration)
// 	if e != nil {
// 		log.Fatal(e)
// 	}

// 	var h []float64
// 	var l []float64
// 	var c []float64

// 	for _, data := range ohlc {
// 		h = append(h, data.High)
// 		l = append(l, data.Low)
// 		c = append(c, data.Close)
// 	}

// 	d := indicators.Donchain(h, l, 40)

// 	if c[len(c)-1] > d.High[len(d.High)-1] {
// 		buy = true
// 		return buy, shortExitSignal
// 	}

// 	if c[len(c)-1] < d.Low[len(d.Low)-1] {
// 		longExitSignal = true
// 		return sell, longExitSignal
// 	}
// 	return false, false
// }
